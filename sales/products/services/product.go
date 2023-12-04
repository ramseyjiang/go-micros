package service

import (
	"context"
	pb "github.com/ramseyjiang/go-micros/sales/products/proto"
	"github.com/ramseyjiang/go-micros/sales/products/repos"
)

type ProductService struct {
	repo repos.ProductRepositoryInterface
}

func NewProductService(repo repos.ProductRepositoryInterface) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	products, err := s.repo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductsResponse{Products: products}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	product := &pb.Product{
		// Generate a unique ID for the product and assign it here
		Name:  req.Name,
		Price: req.Price,
	}
	err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}
