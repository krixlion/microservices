package service

import (
	"context"
	"log"
	"net"
	"user/pkg/grpc/pb"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type UserHandler struct {
	pb.UnimplementedUserServer
}

// Get all Users with filter - A server-to-client streaming RPC.
func (s *UserHandler) Index(_ *pb.UserFilter, stream pb.User_IndexServer) error {
	err := stream.Send(&pb.UserRequest{
		Id:    1,
		Name:  "name",
		Email: "email",
		Phone: "phone lol",
	})
	if err != nil {
		return err
	}
	return nil
}

// Create a new User - A simple RPC
func (s *UserHandler) Create(_ context.Context, _ *pb.UserRequest) (*pb.UserResponse, error) {
	panic("not implemented") // TODO: Implement
}

func Run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterUserServer(s, &UserHandler{})
	s.Serve(lis)
}
