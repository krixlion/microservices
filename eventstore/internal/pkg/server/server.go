package server

import (
	"context"
	"eventstore/internal/pb"
	"eventstore/internal/pkg/event_repository"
	"eventstore/pkg/repository"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "eventstore"
)

type EventStoreServer struct {
	pb.UnimplementedEventStoreServer
	repo     repository.Repository[*pb.Event]
	rabbitmq *amqp.Connection
	// logger kitlog.Logger
}

func MakeEventStoreServer() EventStoreServer {
	pass := os.Getenv("RABBIT_PASS")
	port := os.Getenv("RABBIT_PORT")
	host := os.Getenv("RABBIT_HOST")
	user := os.Getenv("RABBIT_USER")

	uri := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)

	rabbitmq, err := amqp.Dial(uri)
	if err != nil {
		panic(err)
	}

	return EventStoreServer{
		repo:     event_repository.MakeEventRepository(),
		rabbitmq: rabbitmq,
		// logger: log.MakeLogger(),
	}
}

// Closes all connections with th DB and RabbitMQ, returns multiple errors wrapped into one error
func (s EventStoreServer) Close(ctx context.Context) (err error) {
	errMq := s.rabbitmq.Close()
	errRepo := s.repo.Close(ctx)

	switch {
	case errMq != nil:
		err = fmt.Errorf("failed to close conn: %q", errMq)
	case errRepo != nil:
		err = fmt.Errorf("failed to close conn: %q", errRepo)
	case errMq != nil && errRepo != nil:
		err = fmt.Errorf("failed to close conns: %q; %q", errMq, errRepo)
	default:
		return nil
	}
	return err
}
