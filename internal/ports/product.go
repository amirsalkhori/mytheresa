package ports

import (
	"context"
	"mytheresa/internal/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	ListProducts(ctx context.Context, filters map[string]interface{}, pageSize int, lastID uint32) ([]domain.Product, domain.Pagination, error)
}

type ProductService interface {
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	ListProducts(ctx context.Context, filters map[string]interface{}, pageSize int, lastID uint32) ([]domain.ProductDiscount, domain.Pagination, error)
}
