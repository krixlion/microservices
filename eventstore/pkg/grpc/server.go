package grpc

import (
	"context"
	"eventstore/pkg/grpc/pb"
	"eventstore/pkg/repository"
	"os"

	"github.com/go-kit/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventStoreServer struct {
	pb.UnimplementedEventStoreServer
	db repository.Repository[*pb.Event]
}

func NewEventStoreServer() *EventStoreServer {
	return &EventStoreServer{
		db: repository.MakeEventRepository(),
	}
}

func (s *EventStoreServer) Create(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	panic("not implemented")
	s.db.Create(req.Event)
	return &pb.CreateEventResponse{}, nil
}

func (s *EventStoreServer) Get(ctx context.Context, rq *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	panic("not implemented")
	id := rq.GetEventId()
	s.db.Get(id)
	return &pb.GetEventsResponse{}, nil
}

func (s *EventStoreServer) GetStream(req *pb.GetEventsRequest, srv pb.EventStore_GetStreamServer) error {

	panic("not implemented")

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	events, err := s.db.Index(req.GetAggregateId())
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := srv.Send(event); err != nil {
			logger.Log("transport", "grpc", "msg", "failed to send in stream", "err", err)
			return err
		}
	}
	return nil
}

func (s *EventStoreServer) Publish(event *pb.Event) error {
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
