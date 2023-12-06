package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ramseyjiang/go-micros/sales/grpc-gateway/products"
	"github.com/ramseyjiang/go-micros/sales/grpc-gateway/trade"
	"google.golang.org/grpc"
)

const (
	tradeSalesPort        = "SALES_PORT"
	defaultSalesPort      = ":8080"
	defaultProductService = ":9011"
	defaultTradeService   = ":9012"
	productServiceEnvVar  = "PRODUCT_SERVICE_ADDR"
	tradeServiceEnvVar    = "TRADE_SERVICE_ADDR"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()} // Use WithInsecure for non-TLS connections

	// Get the product service address from the environment variable or default to localhost
	productServiceAddress := os.Getenv(productServiceEnvVar)
	if productServiceAddress == "" {
		productServiceAddress = defaultProductService
	}
	err := products.RegisterProductServiceHandlerFromEndpoint(ctx, mux, productServiceAddress, opts)
	if err != nil {
		log.Fatalf("Failed to register product gRPC gateway: %v", err)
	}

	// Get the trade service address from the environment variable or default to localhost
	tradeServiceAddress := os.Getenv(tradeServiceEnvVar)
	if tradeServiceAddress == "" {
		tradeServiceAddress = defaultTradeService
	}
	err = trade.RegisterSalesServiceHandlerFromEndpoint(ctx, mux, tradeServiceAddress, opts)
	if err != nil {
		log.Fatalf("Failed to register trade gRPC gateway: %v", err)
	}

	salesPort := os.Getenv(tradeSalesPort)
	if salesPort == "" {
		salesPort = defaultSalesPort
	}

	log.Printf("Starting gRPC Gateway on port %s\n", salesPort)
	if err = http.ListenAndServe(salesPort, mux); err != nil {
		log.Fatalf("Failed to serve gRPC-Gateway: %v", err)
	}
}
