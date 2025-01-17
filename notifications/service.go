package service

import (
	"context"
	"log"

	"github.com/jaydee029/SeeALie/request/internal/database"
)

type Service struct {
	Domain      string
	AdminEmail  string
	AdminPasswd string
	DB          *database.Queries
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
