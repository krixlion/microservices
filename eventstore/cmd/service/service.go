package service

import (
	mygrpc "eventstore/pkg/grpc"
	"eventstore/pkg/grpc/pb"
	"eventstore/pkg/log"
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func Run() {
	logger := log.MakeLogger()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		logger.Log("transport", "net/tcp", "msg", "failed listening", "err", err)
	}

	grpcServer := grpc.NewServer()
	server := mygrpc.MakeEventStoreServer()

	pb.RegisterEventStoreServer(grpcServer, server)

	logger.Log("transport", "net/tcp", "msg", "listening")
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Log("transport", "grpc", "msg", "failed to serve", "err", err)
	}
}
