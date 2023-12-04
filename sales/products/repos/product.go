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
	var products []*pb.Product

	iter := r.redisClient.Scan(ctx, 0, "product:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

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
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product keys: %v", err)
	}

	return products, nil
}

// CreateProduct stores a new product in Redis.
func (r *ProductRepository) CreateProduct(ctx context.Context, product *pb.Product) error {
	nextID, err := r.redisClient.Incr(ctx, "product:next_id").Result()
	if err != nil {
		return fmt.Errorf("error generating new ID for product: %v", err)
	}

	productID := strconv.FormatInt(nextID, 10)
	productKey := fmt.Sprintf("product:%s", productID)

	// Store the product data in a hash
	_, err = r.redisClient.HMSet(ctx, productKey, map[string]interface{}{
		"id":    productID,
		"name":  product.Name,
		"price": product.Price,
	}).Result()

	if err != nil {
		return fmt.Errorf("error storing product in Redis: %v", err)
	}

	// Update the product's ID with the new unique ID
	product.Id = productID

	return nil
}
