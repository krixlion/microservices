package service

import (
	"context"
	mygrpc "eventstore/pkg/grpc"
	"eventstore/pkg/grpc/pb"
	"eventstore/pkg/log"
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

var (
	port int
)

func init() {
	portFlag := flag.Int("port", 50051, "The server port")
	flag.Parse()
	port = *portFlag
}

func Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.PrintLn("transport", "net/tcp", "msg", "failed listening", "err", err)
	}

	grpcServer := grpc.NewServer()
	server := mygrpc.MakeEventStoreServer()
	defer log.PrintLn(server.Close(context.Background()))

	pb.RegisterEventStoreServer(grpcServer, server)

	log.PrintLn("transport", "net/tcp", "msg", "listening")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.PrintLn("transport", "grpc", "msg", "failed to serve", "err", err)
	}
}
