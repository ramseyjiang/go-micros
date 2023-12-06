package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ramseyjiang/go-micros/sales/grpc-gateway/routes"
)

const (
	tradeSalesPort   = "SALES_PORT"
	defaultSalesPort = ":8080"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	salesPort := os.Getenv(tradeSalesPort)
	if salesPort == "" {
		salesPort = defaultSalesPort
	}

	httpHandler := routes.SetupRoutes(ctx)

	log.Printf("Starting gRPC Gateway on port %s\n", salesPort)
	if err := http.ListenAndServe(salesPort, httpHandler); err != nil {
		log.Fatalf("Failed to serve gRPC-Gateway: %v", err)
	}
}
