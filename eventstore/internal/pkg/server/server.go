package server

import (
	"context"
	"encoding/json"
	"eventstore/internal/pb"
	"eventstore/internal/pkg/event_repository"
	"eventstore/pkg/log"
	"eventstore/pkg/repository"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	exchangeName = "eventstore"
	rabbitUri    = "amqp://guest:guest@rabbitmq-service:5672/"
)

type EventStoreServer struct {
	pb.UnimplementedEventStoreServer
	repo     repository.Repository[*pb.Event]
	rabbitmq *amqp.Connection
	// logger kitlog.Logger
}

func MakeEventStoreServer() EventStoreServer {
	rabbitmq, err := amqp.Dial(rabbitUri)
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

func (s EventStoreServer) Create(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	// Save document to DB
	if err := s.repo.Create(ctx, req.Event); err != nil {
		log.PrintLn("transport", "grpc", "procedure", "create", "msg", "failure", "err", err)
		return &pb.CreateEventResponse{
			IsSuccess: false,
		}, status.Convert(err).Err()
	}
	log.PrintLn("transport", "grpc", "procedure", "create", "msg", "success")
	return &pb.CreateEventResponse{
		IsSuccess: true,
	}, nil
}

func (s EventStoreServer) Get(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	id := req.GetEventId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid ID")
	}
	// Get document from DB
	event, err := s.repo.Get(ctx, id)

	if err != nil {
		log.PrintLn("transport", "grpc", "procedure", "get", "msg", "failure", "err", err)
		return nil, status.Error(codes.NotFound, "Event not found")
	}

	log.PrintLn("transport", "grpc", "procedure", "get", "msg", "success")
	return &pb.GetEventsResponse{Event: event}, nil
}

func (s EventStoreServer) GetStream(req *pb.GetEventsRequest, stream pb.EventStore_GetStreamServer) error {
	ctx := stream.Context()
	events, err := s.repo.Index(ctx)
	if err != nil {
		log.PrintLn("transport", "grpc", "procedure", "getStream", "msg", "failure", "err", err)
		return status.Convert(err).Err()
	}

	// Begin streaming events
	for _, event := range events {
		// If client is unavailable Send() will return an error and abort streaming
		if err := stream.Send(event); err != nil {
			log.PrintLn("transport", "grpc", "procedure", "getStream", "msg", "failure", "err", err)
			return status.Convert(err).Err()
		}
		log.PrintLn("transport", "grpc", "procedure", "getStream", "msg", "success")
	}
	return nil
}

func (s EventStoreServer) PublishEvent(ctx context.Context, event *pb.Event) error {

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ch, err := s.rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,       // name
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		exchangeName,    // exchange
		event.EventType, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonEvent),
		})
	if err != nil {
		return err
	}

	return nil
}

func (s EventStoreServer) ListenEvents(abort <-chan struct{}) (<-chan string, error) {

	ch, err := s.rabbitmq.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,       // name
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		queue.Name,   // queue name
		"",           // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}

	events, err := ch.Consume(
		queue.Name,   // queue
		"eventstore", // consumer
		true,         // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}

	channel := make(chan string)
	go func() {
		defer close(channel)
		for msg := range events {
			select {
			case channel <- string(msg.Body):
				log.PrintLn("transport", "amqp", "procedure", "ListenEvents", "msg", "received")
			case <-abort: // receive on closed channel can proceed immediately
				log.PrintLn("transport", "amqp", "procedure", "ListenEvents", "msg", "aborting")
				return
			}
		}
	}()

	return channel, nil
}
