package repos

import (
	"context"
	"fmt"
	"strconv"

	// Import the product proto package if you're using gRPC to get products.
	productpb "github.com/ramseyjiang/go-micros/sales/products/proto"
	"google.golang.org/grpc"
)

type TradeRepository interface {
	CheckProductExists(ctx context.Context, productID string) (bool, float32, error)
}

type tradeRepositoryImpl struct {
	productServiceClient productpb.ProductServiceClient
}

// NewTradeRepository creates a new instance of a TradeRepository.
func NewTradeRepository(productServiceAddress string) (TradeRepository, error) {
	conn, err := grpc.Dial(productServiceAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %v", err)
	}
	client := productpb.NewProductServiceClient(conn)
	return &tradeRepositoryImpl{productServiceClient: client}, nil
}

// CheckProductExists checks if the product with the given ID exists and returns its price.
func (r *tradeRepositoryImpl) CheckProductExists(ctx context.Context, productID string) (bool, float32, error) {
	// Get all products from the product service.
	products, err := r.productServiceClient.GetProducts(ctx, &productpb.GetProductsRequest{})
	if err != nil {
		return false, 0, fmt.Errorf("failed to retrieve products: %v", err)
	}

	// Search for the product by ID and retrieve its price.
	for _, product := range products.Products {
		if product.Id == productID {
			// Assuming the price is a string, convert it to float32.
			price, err := strconv.ParseFloat(product.Price, 32)
			if err != nil {
				return true, 0, fmt.Errorf("failed to parse product price: %v", err)
			}
			return true, float32(price), nil
		}
	}

	// Product not found.
	return false, 0, nil
}
