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

func (m *MockDiscountService) CreateDiscount(ctx context.Context, discount domain.Discount) (domain.Discount, error) {
	args := m.Called(ctx, discount)
	return args.Get(0).(domain.Discount), args.Error(1)
}

func (m *MockDiscountService) StoreDiscountsInRedis(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
