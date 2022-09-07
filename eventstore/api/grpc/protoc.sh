#! /bin/sh.
DIR="$(dirname "${BASH_SOURCE[0]}")"
DIR="$(realpath "${DIR}")"
GO_PKG_PATH="$(cd $DIR/../../internal/pb/ && pwd)"

protoc --go_out=paths=source_relative:$GO_PKG_PATH --go-grpc_out=paths=source_relative:$GO_PKG_PATH -I $DIR event_store.proto
protoc-go-inject-tag -input="$GO_PKG_PATH/*.pb.go"