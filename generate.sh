#!/bin/bash

echo running...

set -ex
# Setup greet pb file
protoc messagepb/message.proto --go_out=plugins=grpc:.
