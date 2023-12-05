package services

import (
	"context"
	"testing"

	tradepb "github.com/ramseyjiang/go-micros/sales/trade/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockTradeRepository is a mock implementation of the TradeRepository
type MockTradeRepository struct {
	// Add fields to simulate different responses
	productExists bool
	price         float32
	err           error
}

func (m *MockTradeRepository) CheckProductExists(ctx context.Context, productID string) (bool, float32, error) {
	return m.productExists, m.price, m.err
}

func TestSalesService_CreateSale(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		lineItems      []*tradepb.LineItem
		discountAmount float32
		mockRepo       *MockTradeRepository
		wantTotalPrice float32
		wantErrCode    codes.Code
	}{
		{
			name: "Successful Sale",
			lineItems: []*tradepb.LineItem{
				{ProductId: "1", Quantity: 2},
				{ProductId: "2", Quantity: 1},
			},
			discountAmount: 5,
			mockRepo: &MockTradeRepository{
				productExists: true,
				price:         10,
			},
			wantTotalPrice: 25, // (2*10 + 1*10) - 5
			wantErrCode:    codes.OK,
		},
		{
			name: "Product Not Exists",
			lineItems: []*tradepb.LineItem{
				{ProductId: "3", Quantity: 1},
			},
			discountAmount: 0,
			mockRepo: &MockTradeRepository{
				productExists: false,
				price:         0,
			},
			wantTotalPrice: 0,
			wantErrCode:    codes.NotFound,
		},
		// Add more test cases as needed...
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSalesService(tt.mockRepo)
			req := &tradepb.CreateSaleRequest{
				LineItems:      tt.lineItems,
				DiscountAmount: tt.discountAmount,
			}
			resp, err := service.CreateSale(context.Background(), req)

			if err != nil && status.Code(err) != tt.wantErrCode {
				t.Errorf("CreateSale() error = %v, wantErr %v", err, tt.wantErrCode)
				return
			}
			if err == nil && resp.TotalPrice.Value != tt.wantTotalPrice {
				t.Errorf("CreateSale() got total price = %v, want %v", resp.TotalPrice.Value, tt.wantTotalPrice)
			}
		})
	}
}
