package service

import (
	"context"
	"net"
	"os"
	"user/pkg/grpc/pb"

	"github.com/go-kit/log"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var logger log.Logger

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
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	if err != nil {
		return err
	}
	return nil
}

// Create a new User - A simple RPC
func (s *UserHandler) Create(_ context.Context, _ *pb.UserRequest) (*pb.UserResponse, error) {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	logger.Log("msg", "Create method", "transport", "gRPC", "port", port)
	panic("not implemented") // TODO: Implement
}

func Run() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Log("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterUserServer(s, &UserHandler{})
	logger.Log("msg", "Listening", "transport", "gRPC", "port", port)
	s.Serve(lis)
}
