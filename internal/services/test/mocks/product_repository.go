package mocks

import (
	"context"
	"mytheresa/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *MockProductRepository) ListProducts(ctx context.Context, filters map[string]interface{}, pageSize int, lastID uint32) ([]domain.Product, domain.Pagination, error) {
	args := m.Called(ctx, filters, pageSize, lastID)
	return args.Get(0).([]domain.Product), args.Get(1).(domain.Pagination), args.Error(2)
}
