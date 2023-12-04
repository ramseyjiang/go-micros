package repos

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	pb "github.com/ramseyjiang/go-micros/sales/products/proto"
)

type ProductRepositoryInterface interface {
	GetProducts(ctx context.Context) ([]*pb.Product, error)
	CreateProduct(ctx context.Context, product *pb.Product) error
}

// ProductRepository handles the interaction with Redis for product data.
type ProductRepository struct {
	redisClient *redis.Client
}

// NewProductRepository creates a new instance of ProductRepository.
func NewProductRepository(redisClient *redis.Client) *ProductRepository {
	return &ProductRepository{
		redisClient: redisClient,
	}
}

// GetProducts retrieves all products from Redis.
func (r *ProductRepository) GetProducts(ctx context.Context) ([]*pb.Product, error) {
	// Get all keys for product from redis cli
	productKeys, err := r.redisClient.Keys(ctx, "product:*").Result()
	if err != nil {
		return nil, err
	}

	var products []*pb.Product
	for _, key := range productKeys {
		if key == "product:next_id" {
			// Skip the product:next_id key
			continue
		}

		// Because the data in the key is hash, that's why here uses HGetAll method
		productData, err := r.redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("error retrieving product from Redis: %v", err)
		}

		product := &pb.Product{
			Id:    productData["id"],
			Name:  productData["name"],
			Price: productData["price"],
		}
		products = append(products, product)
	}

	return products, nil
}

// CreateProduct stores a new product in Redis.
func (r *ProductRepository) CreateProduct(ctx context.Context, product *pb.Product) error {
	// Increment an integer value stored at product:next_id.
	// If it doesn’t exist, Redis creates it and sets it to 1.
	// This is a way to keep track of the number of products and ensure each one has a unique ID.
	nextID, err := r.redisClient.Incr(ctx, "product:next_id").Result()
	if err != nil {
		return fmt.Errorf("error generating new ID for product: %v", err)
	}

	productID := strconv.FormatInt(nextID, 10)
	productKey := fmt.Sprintf("product:%s", productID)

	// Store the product data in a hash
	if _, err = r.redisClient.HMSet(ctx, productKey, map[string]interface{}{
		"id":    productID,
		"name":  product.Name,
		"price": product.Price,
	}).Result(); err != nil {
		return fmt.Errorf("error storing product in Redis: %v", err)
	}

	// Update the product's ID with the new unique ID
	product.Id = productID

	return nil
}
