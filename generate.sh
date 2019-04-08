#!/bin/bash

echo running...

set -x

if ! which protoc-gen-go 1> /dev/null; then go get -u github.com/golang/protobuf/protoc-gen-go; fi

set -e

# Setup greet pb file
# protoc messagepb/message.proto --go_out=plugins=grpc:.

# Setup grpc pb
protoc -I. \
  -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  messagepb/message.proto

# Setup rest pb
protoc -I. \
  -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. \
  messagepb/message.proto

# Setup swagger pb
protoc -I. \
  -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true:. \
  messagepb/message.proto
