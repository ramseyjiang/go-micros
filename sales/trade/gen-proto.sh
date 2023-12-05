#!/bin/bash

# Navigate to the directory where your proto file is located, if necessary

# Run the protoc command with the Go plugins to generate the Go files
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       ./proto/trade.proto

echo "Protobuf files have been generated."