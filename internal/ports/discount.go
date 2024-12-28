package ports

import (
	"context"
	"mytheresa/internal/app/dto"
	"mytheresa/internal/domain"

	"github.com/gin-gonic/gin"
)

type DisocuntRepository interface {
	CreateDiscount(ctx context.Context, disocunt domain.Discount) (domain.Discount, error)
	GetAllDiscounts() ([]domain.Discount, error)
}

type DiscountService interface {
	CreateDiscount(ctx context.Context, discount domain.Discount) (domain.Discount, error)
	GetBestDiscount(ctx context.Context, sku, category string) (domain.Discount, error)
	StoreDiscountsInRedis(ctx context.Context) error
}

type DiscountHandler interface {
	CreateDiscount(c *gin.Context)
	CreateDiscountFromFile(ctx context.Context, discountsRoot dto.DisocuntRoot)
}
