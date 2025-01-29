package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jaydee029/SeeALie/request/internal/database"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Service struct {
	Domain      string
	AdminEmail  string
	AdminPasswd string
	DB          *database.Queries
	Pubsub      *amqp.Connection
}

func (s *Service) Run(ctx context.Context) {
	ntfticker := time.NewTicker(1 * time.Second)
	done := make(chan struct{}, 1)

	defer func() {
		ntfticker.Stop()
		close(done)
	}()

	for {
		select {
		case <-ntfticker.C:
			select {
			case done <- struct{}{}:
				go func() {
					defer func() { <-done }()
					s.ProcessRequestNotification(ctx)
					s.ProcessStatusNotification(ctx)

				}()
			default:
				fmt.Println("Previous job didnt complete yet")
			}

		case <-ctx.Done():
			return
		}
	}

}

func (s *Service) UpdateSentStatus(ctx context.Context, req *database.DequeRequestsRow, recieved_err error) error {

	switch recieved_err {
	case nil:
		err := s.DB.MailSent(ctx, req.ConnectionID)
		return err
	default:
		row, err := s.DB.MailnotSent(ctx, req.ConnectionID)
		if row.SentAttempts.Valid && row.SentAttempts.Int32 == 3 && row.StatusSent == "PENDING" {
			log.Println("ALL attempts to send email request exhausted")
		}
		return err
	}
}
