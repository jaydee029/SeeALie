package internal

import (
	"context"
	"log"

	"github.com/jaydee029/SeeALie/pubsub"
)

func eventhandler(val any) pubsub.Acktype {

	
}

func (s *Service) Eventlistener(ctx context.Context) {

	errch := make(chan error)
	go func() {
		err := s.Pubsub.SubscribeGob("SeeALie_User_direct", "user_email_queue", "user.email", pubsub.DurableQueue, eventhandler)
		if err != nil {
			errch <- err
		}
	}()

	go func() {
		err := s.Pubsub.SubscribeGob("SeeALie_User_direct", "user_mapping_queue", "user.mapping", pubsub.DurableQueue, eventhandler)
		if err != nil {
			errch <- err
		}
	}()

	select {
	case <-ctx.Done():
		close(errch)
		log.Println("closing event listener ...")
		return
	case err := <-errch:
		log.Panicln(err)
	}
}

/*
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
*/
