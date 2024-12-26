package ports

import (
	"context"
	"mytheresa/internal/domain"

	"github.com/gin-gonic/gin"
)

type DisocuntRepository interface {
	CreateDiscount(ctx context.Context, disocunt domain.Discount) (domain.Discount, error)
	// ListProducts(ctx context.Context, filters map[string]interface{}, pageSize, page int) ([]domain.Product, domain.Pagination, error)
}

type DiscountService interface {
	CreateDiscount(ctx context.Context, discount domain.Discount) (domain.Discount, error)
	// ListProducts(ctx context.Context, filters map[string]interface{}, pageSize, page int) ([]domain.Product, domain.Pagination, error)
}

type DiscountHandler interface {
	CreateDiscount(c *gin.Context)
	// ListProducts(ctx context.Context, filters map[string]interface{}, pageSize, page int) ([]domain.Product, domain.Pagination, error)
}
