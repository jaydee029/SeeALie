package main

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jaydee029/SeeALie/request/internal/database"
	"github.com/jaydee029/SeeALie/request/protos"
)

func (s *server) Run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.ProcessRequest(ctx)
		case <-ctx.Done():
			return
		}
	}

}

func (s *server) EnqueRequest(ctx context.Context, r *protos.EmailRequest) (*protos.EmailResponse, error) {
	Sentto := r.Sentto
	Sentby := r.Sentby

	if_exists, err := s.DB.User_exists(ctx, Sentto)

	if !if_exists {
		return &protos.EmailResponse{
			Response: "",
			Error:    true,
		}, errors.New("user doesn't exist")

	} else if err != nil {
		return &protos.EmailResponse{
			Response: "",
			Error:    true,
		}, err
	}
	_, err = s.DB.Add_connection_Info(ctx, database.Add_connection_InfoParams{
		RequestBy:    Sentby,
		RequestTo:    Sentto,
		ConnectionID: uuid.New(),
		CreatedAt:    time.Now().UTC(),
	})

	if err != nil {
		return &protos.EmailResponse{
			Response: "",
			Error:    true,
		}, err
	}

	return &protos.EmailResponse{
		Response: "connection request is being processed",
		Error:    false,
	}, nil
}