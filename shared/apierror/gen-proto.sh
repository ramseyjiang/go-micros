#!/bin/bash

protoc ./apierror.proto \
  -I . \
  -I "$GOPATH"/src/ \
  -I "$GOPATH"/src/github.com/googleapis/ \
  --go_opt=paths=source_relative \
  --go_out=./ \
  --go-grpc_out=./
