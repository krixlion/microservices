#! /bin/sh.
DIR="$(dirname "${BASH_SOURCE[0]}")"
DIR="$(realpath "${DIR}")"
GO_PB_PATH="$(cd $DIR/../../internal/pb/ && pwd)"

protoc --go_out=paths=source_relative:$GO_PB_PATH --go-grpc_out=paths=source_relative:$GO_PB_PATH -I $DIR event_store.proto
protoc-go-inject-tag -input="$GO_PB_PATH/*.pb.go"