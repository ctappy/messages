#!/bin/bash

echo running...

if ! which protoc-gen-go 1> /dev/null; then go get -u github.com/golang/protobuf/protoc-gen-go; fi

set -ex
# Setup greet pb file
protoc messagepb/message.proto --go_out=plugins=grpc:.
