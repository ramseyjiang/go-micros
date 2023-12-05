package services

import (
	"context"
	"fmt"

	"github.com/ramseyjiang/go-micros/sales/trade/internal/repos"
	tradepb "github.com/ramseyjiang/go-micros/sales/trade/proto"
)

type SalesService struct {
	//  the UnimplementedSalesServiceServer satisfies the gRPC interface, including the forward compatibility method.
	tradepb.UnimplementedSalesServiceServer
	repo repos.TradeRepository
}

func NewSalesService(repo repos.TradeRepository) *SalesService {
	return &SalesService{repo: repo}
}

func (s *SalesService) CreateSale(ctx context.Context, req *tradepb.CreateSaleRequest) (*tradepb.CreateSaleResponse, error) {
	var totalSalePrice float32

	for _, item := range req.LineItems {
		// Check if product exists and get its price
		exists, price, err := s.repo.CheckProductExists(ctx, item.ProductId)
		if err != nil {
			return nil, fmt.Errorf("error checking product existence: %v", err)
		}
		if !exists {
			return nil, fmt.Errorf("product with ID %s does not exist", item.ProductId)
		}

		// Calculate total price for the line item
		lineTotal := price * float32(item.Quantity)
		totalSalePrice += lineTotal
	}

	// Create the sale response with the total sale price and line items
	return &tradepb.CreateSaleResponse{
		TotalPrice: totalSalePrice,
		LineItems:  req.LineItems,
	}, nil
}
