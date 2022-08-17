package grpc

import (
	"context"
	"eventstore/pkg/grpc/pb"
	"eventstore/pkg/log"
	"eventstore/pkg/repository"

	kitlog "github.com/go-kit/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventStoreServer struct {
	pb.UnimplementedEventStoreServer
	repo   repository.Repository[*pb.Event]
	logger kitlog.Logger
}

func MakeEventStoreServer() EventStoreServer {
	return EventStoreServer{
		repo:   repository.MakeEventRepository(),
		logger: log.MakeLogger(),
	}
}

func (s EventStoreServer) Create(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	err := s.repo.Create(ctx, req.Event)

	if err != nil {
		return &pb.CreateEventResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, err
	}

	return &pb.CreateEventResponse{
		IsSuccess: true,
		Error:     "",
	}, nil
}

func (s EventStoreServer) Get(context.Context, *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	panic("not implemented")
}

func (s EventStoreServer) GetStream(*pb.GetEventsRequest, pb.EventStore_GetStreamServer) error {
	panic("not implemented")
}

func (s EventStoreServer) Publish(event *pb.Event) error {
	const uri = "amqp://guest:guest@rabbitmq-service:5672/"
	rabbitmq, err := amqp.Dial(uri)
	if err != nil {
		return err
	}
	defer rabbitmq.Close()
	ch, err := rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		event.ChannelName, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(
		context.Background(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event.EventData),
		})
	if err != nil {
		return err
	}
	return nil
}
