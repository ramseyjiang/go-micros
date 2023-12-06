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
	productServiceEnvVar      = "PRODUCT_SERVICE_ADDR"
	redisEnvVar               = "REDIS_ADDR"
	defaultProductServicePort = ":9011"
	defaultRedisPort          = "localhost:6379" // listen on localhost:6379 port, not listen all :6379
)

func main() {
	flag.Parse()

	// Initialize a Redis client
	redisAddr := os.Getenv(redisEnvVar) // For example: "localhost:6379"
	if redisAddr == "" {
		redisAddr = defaultRedisPort
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

	productServicePort := os.Getenv(productServiceEnvVar)
	if productServicePort == "" {
		productServicePort = defaultProductServicePort
	}
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
