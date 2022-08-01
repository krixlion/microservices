package grpc

import (
	"context"
	"eventstore/pkg/grpc/pb"
	"eventstore/pkg/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventStoreServer struct {
	pb.UnimplementedEventStoreServer
	db repository.Repository[pb.Event, string]
}

func NewEventStoreServer() *EventStoreServer {
	return &EventStoreServer{}
}

func (s *EventStoreServer) Create(context.Context, *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	panic("not implemented")
}

func (s *EventStoreServer) Get(context.Context, *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	panic("not implemented")
}

func (s *EventStoreServer) GetStream(*pb.GetEventsRequest, pb.EventStore_GetStreamServer) error {
	panic("not implemented")
}

func (s *EventStoreServer) Publish(event *pb.Event) error {
	rabbitmq, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer rabbitmq.Close()
	ch, err := rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
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
