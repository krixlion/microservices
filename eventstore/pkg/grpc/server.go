package grpc

import (
	"context"
	"eventstore/pkg/grpc/pb"
	"eventstore/pkg/log"
	"eventstore/pkg/repository"
	"fmt"

	kitlog "github.com/go-kit/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func FmtDupKeyErr(err error) error {
	logger := log.MakeLogger()
	badReq := &errdetails.BadRequest{}
	violation := &errdetails.BadRequest_FieldViolation{
		Field:       "event_id",
		Description: err.Error(),
	}
	badReq.FieldViolations = append(badReq.FieldViolations, violation)

	st, statusErr := status.New(codes.AlreadyExists, "Event with given ID already exists").WithDetails(badReq)
	if statusErr != nil {
		logger.Log("transport", "gRPC", "msg", "Unexpected error attaching metadata", "err", err)
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}

	err = st.Err()
	return err
}

func (s EventStoreServer) Create(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	// Save document to DB
	if err := s.repo.Create(ctx, req.Event); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			err = FmtDupKeyErr(err)
		}
		return &pb.CreateEventResponse{
			IsSuccess: false,
		}, err
	}

	return &pb.CreateEventResponse{
		IsSuccess: true,
	}, nil
}

func (s EventStoreServer) Get(ctx context.Context, rq *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	id := rq.GetEventId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid ID")
	}
	// Get document from DB
	event, err := s.repo.Get(ctx, id)
	if event == nil {
		return nil, status.Error(codes.NotFound, "Event not found")
	}
	if err != nil {
		return nil, err
	}

	s.logger.Log("transport", "grpc", "msg", "succesfully sent get response")
	return &pb.GetEventsResponse{Event: event}, nil
}

func (s EventStoreServer) GetStream(req *pb.GetEventsRequest, stream pb.EventStore_GetStreamServer) error {
	ctx := stream.Context()
	events, err := s.repo.Index(ctx)
	if err != nil {
		s.logger.Log("transport", "mongodb", "msg", "failed to retrieve events", "err", err)
		return status.Convert(err).Err()
	}

	// Begin streaming events
	for _, event := range events {
		// If client is unavailable Send() will return an error and abort streaming
		if err := stream.Send(event); err != nil {
			s.logger.Log("transport", "grpc", "msg", "stopped streaming events", "err", err)
			return status.Convert(err).Err()
		}
		s.logger.Log("transport", "grpc", "msg", "succesfully sent stream response")
	}
	return nil
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
