package grpc

import (
	"context"
	"eventstore/pkg/grpc/pb"
	"eventstore/pkg/log"
	"eventstore/pkg/repository"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventStoreServer struct {
	pb.UnimplementedEventStoreServer
	repo repository.Repository[*pb.Event]
	// logger kitlog.Logger
}

func MakeEventStoreServer() EventStoreServer {
	return EventStoreServer{
		repo: repository.MakeEventRepository(),
		// logger: log.MakeLogger(),
	}
}

func (s EventStoreServer) Create(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	// Save document to DB
	if err := s.repo.Create(ctx, req.Event); err != nil {
		return &pb.CreateEventResponse{
			IsSuccess: false,
		}, status.Convert(err).Err()
	}
	log.Println("transport", "grpc", "procedure", "create", "msg", "success")
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
		log.Println("transport", "grpc", "procedure", "get", "msg", "failure", "err", err)
		return nil, status.Error(codes.NotFound, "Event not found")
	}

	log.Println("transport", "grpc", "procedure", "get", "msg", "success")
	return &pb.GetEventsResponse{Event: event}, nil
}

func (s EventStoreServer) GetStream(req *pb.GetEventsRequest, stream pb.EventStore_GetStreamServer) error {
	ctx := stream.Context()
	events, err := s.repo.Index(ctx)
	if err != nil {
		log.Println("transport", "grpc", "procedure", "getStream", "msg", "failure", "err", err)
		return status.Convert(err).Err()
	}

	// Begin streaming events
	for _, event := range events {
		// If client is unavailable Send() will return an error and abort streaming
		if err := stream.Send(event); err != nil {
			log.Println("transport", "grpc", "procedure", "getStream", "msg", "failure", "err", err)
			return status.Convert(err).Err()
		}
		log.Println("transport", "grpc", "procedure", "getStream", "msg", "success")
	}
	return nil
}

func (s EventStoreServer) Publish(ctx context.Context, event *pb.Event) error {
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
		"events", // name
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
		ctx,
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
