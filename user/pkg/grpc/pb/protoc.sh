#! /bin/sh.
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:.  pkg/grpc/pb/user.proto