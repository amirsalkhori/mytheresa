package mocks

import (
	"context"
	"mytheresa/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockDiscountService struct {
	mock.Mock
}

func (m *MockDiscountService) GetBestDiscount(ctx context.Context, sku string, category string) (domain.Discount, error) {
	args := m.Called(ctx, sku, category)
	return args.Get(0).(domain.Discount), args.Error(1)
}
