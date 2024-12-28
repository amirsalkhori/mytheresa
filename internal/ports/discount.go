package ports

import (
	"context"
	"mytheresa/internal/domain"

	"github.com/gin-gonic/gin"
)

type DisocuntRepository interface {
	CreateDiscount(ctx context.Context, disocunt domain.Discount) (domain.Discount, error)
	GetDiscountsBySKUAndCategory(ctx context.Context, identifier string) (domain.Discount, error)
}

type DiscountService interface {
	CreateDiscount(ctx context.Context, discount domain.Discount) (domain.Discount, error)
	GetBestDiscount(ctx context.Context, product domain.Product) (domain.Discount, error)
}

type DiscountHandler interface {
	CreateDiscount(c *gin.Context)
}
