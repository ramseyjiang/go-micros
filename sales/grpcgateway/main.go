package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ramseyjiang/go-micros/sales/grpc-gateway/products"
	"github.com/ramseyjiang/go-micros/sales/grpc-gateway/trade"
	"google.golang.org/grpc"
)

const (
	salesPort             = ":8080"
	productServiceAddress = ":9011"
	tradeServiceAddress   = ":9012"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()} // Use WithInsecure for non-TLS connections
	err := products.RegisterProductServiceHandlerFromEndpoint(ctx, mux, productServiceAddress, opts)
	if err != nil {
		log.Fatalf("Failed to register product gRPC gateway: %v", err)
	}

	err = trade.RegisterSalesServiceHandlerFromEndpoint(ctx, mux, tradeServiceAddress, opts)
	if err != nil {
		log.Fatalf("Failed to register trade gRPC gateway: %v", err)
	}

	log.Println("Starting serve on port ", salesPort)
	if err := http.ListenAndServe(salesPort, mux); err != nil {
		log.Fatalf("Failed to serve gRPC-Gateway: %v", err)
	}
}
