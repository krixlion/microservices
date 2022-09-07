package server

import (
	"context"
	"eventstore/internal/pb"
	"eventstore/pkg/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
