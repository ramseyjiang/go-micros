package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	pb "github.com/ramseyjiang/go-micros/sales/products/proto"
)

type mockProductRepository struct {
	products []*pb.Product
	err      error
}

func (m *mockProductRepository) GetProducts(ctx context.Context) ([]*pb.Product, error) {
	return m.products, m.err
}

func (m *mockProductRepository) CreateProduct(ctx context.Context, product *pb.Product) error {
	if m.err != nil {
		return m.err
	}
	m.products = append(m.products, product)
	return nil
}

func TestGetProducts(t *testing.T) {
	ctx := context.Background()
	mockProducts := []*pb.Product{
		{Id: "1", Name: "Product 1", Price: "10.99"},
		{Id: "2", Name: "Product 2", Price: "15.99"},
	}

	tests := []struct {
		name         string
		mockRepo     *mockProductRepository
		wantProducts []*pb.Product
		wantErr      bool
	}{
		{
			name:         "Success",
			mockRepo:     &mockProductRepository{products: mockProducts},
			wantProducts: mockProducts,
			wantErr:      false,
		},
		{
			name:         "RepoError",
			mockRepo:     &mockProductRepository{err: errors.New("error")},
			wantProducts: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductService(tt.mockRepo)
			got, err := s.GetProducts(ctx, &pb.GetProductsRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductService.GetProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.GetProducts(), tt.wantProducts) {
				t.Errorf("ProductService.GetProducts() got = %v, want %v", got.GetProducts(), tt.wantProducts)
			}
		})
	}
}

func TestCreateProduct(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		mockRepo *mockProductRepository
		req      *pb.CreateProductRequest
		wantErr  bool
	}{
		{
			name:     "Success",
			mockRepo: &mockProductRepository{},
			req:      &pb.CreateProductRequest{Name: "New Product", Price: "20.99"},
			wantErr:  false,
		},
		{
			name:     "RepoError",
			mockRepo: &mockProductRepository{err: errors.New("error")},
			req:      &pb.CreateProductRequest{Name: "New Product", Price: "20.99"},
			wantErr:  true,
		},
		{
			name:     "EmptyProductName",
			mockRepo: &mockProductRepository{},
			req:      &pb.CreateProductRequest{Name: "", Price: "20.99"},
			wantErr:  true,
		},
		{
			name:     "InvalidProductPrice",
			mockRepo: &mockProductRepository{},
			req:      &pb.CreateProductRequest{Name: "Valid Name", Price: "0"},
			wantErr:  true,
		},
		{
			name:     "NonNumericPrice",
			mockRepo: &mockProductRepository{},
			req:      &pb.CreateProductRequest{Name: "Valid Name", Price: "abc"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductService(tt.mockRepo)
			_, err := s.CreateProduct(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductService.CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
