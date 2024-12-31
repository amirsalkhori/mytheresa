package ports

import (
	"context"
	"mytheresa/internal/domain"
)

type ProductRepository interface {
	ListProducts(ctx context.Context, filters map[string]interface{}, pageSize int, nextID, prevID uint32) ([]domain.Product, domain.Pagination, error)
}

type ProductService interface {
	ListProducts(ctx context.Context, filters map[string]interface{}, pageSize int, nextID, prevID uint32) ([]domain.ProductDiscount, domain.HashedPagination, error)
}
