package handler

import (
	"context"
	"time"

	service "github.com/jaydee029/SeeALie/request"
	"github.com/jaydee029/SeeALie/request/internal/database"
)

type Handler struct {
	DB  *database.Queries
	Svc *service.Service
}

func (h *Handler) Run(ctx context.Context) {
	reqticker := time.NewTicker(1 * time.Second)

	defer func() {
		reqticker.Stop()
	}()

	for {
		select {
		case <-reqticker.C:
			h.Svc.ProcessRequestNotification(ctx)
			h.Svc.ProcessStatusNotification(ctx)
		case <-ctx.Done():
			return
		}
	}

}
