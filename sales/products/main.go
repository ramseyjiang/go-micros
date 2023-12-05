package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/ramseyjiang/go-micros/sales/products/internal/repos"
	"github.com/ramseyjiang/go-micros/sales/products/internal/services"
	pb "github.com/ramseyjiang/go-micros/sales/products/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	productServicePort = "localhost:9011"
	redisAddress       = "localhost:26379" // localhost:6379, docker use 26379
)

func main() {
	flag.Parse()

	// Initialize a Redis client
	redisAddr := os.Getenv("REDIS_ADDR") // For example: "localhost:6379"
	if redisAddr == "" {
		redisAddr = redisAddress
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize the product repository
	productRepo := repos.NewProductRepository(redisClient)

	// Initialize the product service with the repository
	productSvc := services.NewProductService(productRepo)

	// Set up gRPC server
	lis, err := net.Listen("tcp", productServicePort) // Port should be configurable
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register ProductServiceServer
	pb.RegisterProductServiceServer(grpcServer, productSvc)

	// Register reflection service on gRPC server.
	// It is also used for grpcurl sending in the terminal.
	reflection.Register(grpcServer)

	// Graceful shutdown handling
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		redisClient.Close()
		log.Println("Server has been stopped.")
	}()

	log.Println("Starting server on port " + productServicePort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
