package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/jaydee029/SeeALie/notifications/internal/database"
	"github.com/jaydee029/SeeALie/pubsub"
)

type Service struct {
	Domain      string
	AdminEmail  string
	AdminPasswd string
	DB          *database.Queries
	Pubsub      *pubsub.PubSub
}

func NewService(domain, adminEmail, adminPasswd string, db *database.Queries, pubsub *pubsub.PubSub) *Service {
	return &Service{
		Domain:      domain,
		AdminEmail:  adminEmail,
		AdminPasswd: adminPasswd,
		DB:          db,
		Pubsub:      pubsub,
	}
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
					s.ProcessRequest(ctx)
					s.ProcessNotification(ctx)

				}()
			default:
				fmt.Println("Previous job didnt complete yet")
			}

		case <-ctx.Done():
			return
		}
	}

}
