package routes

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ramseyjiang/go-micros/sales/grpc-gateway/middleware"
	"github.com/ramseyjiang/go-micros/sales/grpc-gateway/protos/products"
	"github.com/ramseyjiang/go-micros/sales/grpc-gateway/protos/trade"
	"google.golang.org/grpc"
)

const (
	defaultProductService = ":9011"
	defaultTradeService   = ":9012"
	productServiceEnvVar  = "PRODUCT_SERVICE_ADDR"
	tradeServiceEnvVar    = "TRADE_SERVICE_ADDR"
)

func SetupRoutes(ctx context.Context) http.Handler {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register gRPC handlers
	productServiceAddress := os.Getenv(productServiceEnvVar)
	if productServiceAddress == "" {
		productServiceAddress = defaultProductService
	}
	if err := products.RegisterProductServiceHandlerFromEndpoint(ctx, mux, productServiceAddress, opts); err != nil {
		log.Fatalf("Failed to register product gRPC gateway: %v", err)
	}

	tradeServiceAddress := os.Getenv(tradeServiceEnvVar)
	if tradeServiceAddress == "" {
		tradeServiceAddress = defaultTradeService
	}
	if err := trade.RegisterSalesServiceHandlerFromEndpoint(ctx, mux, tradeServiceAddress, opts); err != nil {
		log.Fatalf("Failed to register trade gRPC gateway: %v", err)
	}

	// Apply the rate limiting middleware
	return middleware.RateLimit(mux)
}
