package mocks

import (
	"context"
	"mytheresa/internal/domain"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDiscountRepository struct {
	mock.Mock
}

func (m *MockDiscountRepository) GetDiscountsBySKUAndCategory(ctx context.Context, SKU, categoryName string) ([]domain.Discount, error) {
	args := m.Called(ctx, SKU, categoryName)
	return args.Get(0).([]domain.Discount), args.Error(1)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}
