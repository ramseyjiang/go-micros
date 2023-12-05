package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ramseyjiang/go-micros/sales/trade/internal/repos"
	"github.com/ramseyjiang/go-micros/sales/trade/internal/services"
	tradepb "github.com/ramseyjiang/go-micros/sales/trade/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	tradeServicePort      = "localhost:9012"
	productServiceAddress = "localhost:9011"
)

func main() {
	// Set up a connection to the ProductService
	tradeRepo, err := repos.NewTradeRepository(productServiceAddress)
	if err != nil {
		log.Fatalf("Failed to initialize trade repository: %v", err)
	}

	// Create a new SalesService instance
	salesService := services.NewSalesService(tradeRepo)

	// Listen on a port
	listener, err := net.Listen("tcp", tradeServicePort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()

	// Register the service with the gRPC server
	tradepb.RegisterSalesServiceServer(grpcServer, salesService)

	// Register reflection service on gRPC server.
	// It is also used for grpcurl sending in the terminal.
	reflection.Register(grpcServer)

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()
	log.Println("Starting server on port ", tradeServicePort)

	// Start the server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
