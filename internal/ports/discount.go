package ports

import (
	"context"
	"mytheresa/internal/domain"
)

type DisocuntRepository interface {
	GetDiscountsBySKUAndCategory(ctx context.Context, SKU, categoryName string) ([]domain.Discount, error)
}

type DiscountService interface {
	GetBestDiscount(ctx context.Context, sku, category string) (domain.Discount, error)
}
