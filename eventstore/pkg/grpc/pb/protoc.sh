#! /bin/sh.
DIR="$(dirname "${BASH_SOURCE[0]}")"
DIR="$(realpath "${DIR}")"

protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:.  pkg/grpc/pb/event_store.proto 
protoc-go-inject-tag -input="$DIR/*.pb.go"