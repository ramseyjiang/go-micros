#!/bin/bash

# Navigate to the directory where your proto file is located, if necessary

# Run the protoc command with the Go plugins to generate the Go files

protoc --go_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:. \
    --grpc-gateway_out=paths=source_relative:. \
    ./protos/products/product.proto ./protos/trade/trade.proto