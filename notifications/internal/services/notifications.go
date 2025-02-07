package internal

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"sync"

	"github.com/jaydee029/SeeALie/notifications/internal/database"
	"golang.org/x/sync/semaphore"
	"gopkg.in/gomail.v2"
)

type NotificationInfo struct {
	request_init_by string
	request_to      string
	Domain          string
}

func (s *Service) ProcessNotification(ctx context.Context) {
	Notifications, err := s.DB.DequeNotifications(ctx)
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(10)

	for _, notification := range Notifications {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Println(err)
		}

		wg.Add(1)
		go func(req *database.DequeNotificationsRow) {
			defer wg.Done()
			defer sem.Release(1)

			err := s.SendNotification(ctx, req)
			err = s.UpdateNotificationSentStatus(ctx, req, err)
			if err != nil {
				log.Println(err)
			}

		}(&notification)

		go func() {
			wg.Wait()
		}()
	}
}

func (s *Service) SendNotification(ctx context.Context, req *database.DequeNotificationsRow) error {

	i := &NotificationInfo{
		request_init_by: req.RequestInitBy,
		request_to:      req.RequestTo,
		Domain:          s.Domain,
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
	m.SetHeader("Subject", "Request Notification")
	m.SetBody("text/html", content)

	d := gomail.NewDialer("smtp.gmail.com", 587, s.AdminEmail, s.AdminPasswd)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateNotificationSentStatus(ctx context.Context, req *database.DequeNotificationsRow, recieved_err error) error {

	switch recieved_err {
	case nil:
		err := s.DB.NotificationSent(ctx, database.NotificationSentParams{
			RequestInitBy: req.RequestInitBy,
			RequestTo:     req.RequestTo,
		})
		return err
	default:
		row, err := s.DB.NotificationnotSent(ctx, database.NotificationnotSentParams{
			RequestInitBy: req.RequestInitBy,
			RequestTo:     req.RequestTo,
		})
		if row.SentAttempts.Valid && row.SentAttempts.Int32 == 3 && !row.StatusSent {
			log.Println("ALL attempts to send email request exhausted")
		}
		return err
	}
}
