package internal

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"sync"

	"gopkg.in/gomail.v2"

	"github.com/jaydee029/SeeALie/notifications/internal/database"
	"golang.org/x/sync/semaphore"
)

type RequestInfo struct {
	UserName     string
	Domain       string
	Connectionid string
}

func (s *Service) ProcessRequest(ctx context.Context) {
	Requests, err := s.DB.DequeRequests(ctx)
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(10)

	for _, request := range Requests {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Println(err)
		}

		wg.Add(1)
		go func(req *database.DequeRequestsRow) {
			defer wg.Done()
			defer sem.Release(1)

			err := s.SendRequest(ctx, req)
			err = s.UpdateRequestSentStatus(ctx, req, err)
			if err != nil {
				log.Println(err)
			}

		}(&request)

		go func() {
			wg.Wait()
		}()
	}
}

func (s *Service) SendRequest(ctx context.Context, req *database.DequeRequestsRow) error {

	i := &RequestInfo{
		UserName:     req.RequestBy,
		Domain:       s.Domain,
		Connectionid: req.ConnectionID.String(),
	}
	targetEmail, err := s.DB.GetRequestInfo(ctx, req.RequestTo)
	if err != nil {
		log.Println(err)
	}
	t := template.New("../../template/request.html")

	t, err = t.ParseFiles("../../template/request.html")

	if err != nil {
		log.Println("error parsing mail template", err)
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		log.Println(err)
	}

	content := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", s.AdminEmail)
	m.SetHeader("To", targetEmail)
	m.SetHeader("Subject", "Connection Request")
	m.SetBody("text/html", content)

	d := gomail.NewDialer("smtp.gmail.com", 587, s.AdminEmail, s.AdminPasswd)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateRequestSentStatus(ctx context.Context, req *database.DequeRequestsRow, recieved_err error) error {

	switch recieved_err {
	case nil:
		err := s.DB.MailSent(ctx, req.ConnectionID)
		return err
	default:
		row, err := s.DB.MailnotSent(ctx, req.ConnectionID)
		if row.SentAttempts.Valid && row.SentAttempts.Int32 == 3 && !row.StatusSent {
			log.Println("ALL attempts to send email request exhausted")
		}
		return err
	}
}
