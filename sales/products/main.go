package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/ramseyjiang/go-micros/sales/products/proto"
	"github.com/ramseyjiang/go-micros/sales/products/repos"
	"github.com/ramseyjiang/go-micros/sales/products/services"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()

	// Initialize a Redis client
	redisAddr := os.Getenv("REDIS_ADDR") // For example: "localhost:6379"
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default address
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
	productSvc := service.NewProductService(productRepo)

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":9011") // Port should be configurable
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register ProductServiceServer
	pb.RegisterProductServiceServer(grpcServer, productSvc)

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

	log.Println("Starting server on port :9011")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
