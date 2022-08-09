package service

import (
	mygrpc "eventstore/pkg/grpc"
	"eventstore/pkg/grpc/pb"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func Run() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	logger.Log("transport", "net/tcp", "msg", "listening")
	if err != nil {
		logger.Log("transport", "net/tcp", "msg", "failed listening", "err", err)
	}

	grpcServer := grpc.NewServer()
	server := mygrpc.NewEventStoreServer()

	pb.RegisterEventStoreServer(grpcServer, server)

	logger.Log("transport", "net/tcp", "msg", "listening")
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Log("transport", "grpc", "msg", "failed serving", "err", err)
	}
}
