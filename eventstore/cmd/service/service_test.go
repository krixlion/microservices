package service_test

import (
	"context"
	mygrpc "eventstore/pkg/grpc"
	"eventstore/pkg/grpc/pb"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	eventstore := mygrpc.MakeEventStoreServer()
	pb.RegisterEventStoreServer(s, eventstore)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreateAndGet(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventStoreClient(conn)
	resp, err := client.Get(ctx, &pb.GetEventsRequest{
		EventId:     "",
		AggregateId: "",
	})
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}
	want := &pb.GetEventsResponse{
		Events: []*pb.Event{},
	}
	if resp != want {
		t.Fatalf("Received response is not equal to expected response")
	}
}
