package services

import (
	"context"
	"fmt"
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"
)

type ProductService struct {
	Repo            ports.ProductRepository
	DiscountService ports.DiscountService
}

func NewProductService(repo ports.ProductRepository, discountService ports.DiscountService) ports.ProductService {
	return &ProductService{Repo: repo, DiscountService: discountService}
}

func (s *ProductService) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	return s.Repo.CreateProduct(ctx, product)
}

func (s ProductService) ListProducts(ctx context.Context, filters map[string]interface{}, pageSize int, lastID uint32) ([]domain.ProductDiscount, domain.Pagination, error) {
	products, pagination, err := s.Repo.ListProducts(ctx, filters, pageSize, lastID)
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	var discountedProducts []domain.ProductDiscount

	for _, product := range products {
		var finalPrice = product.Price
		var discountPercentage *string

		discount, err := s.DiscountService.GetBestDiscount(ctx, product.SKU, product.Category)
		if err != nil || discount.ID == 0 {
			discountPercentage = nil
		} else {
			discountPercentageToString := fmt.Sprintf("%d%%", discount.Percentage)
			discountPercentage = &discountPercentageToString
			finalPrice = applyDiscount(finalPrice, discount.Percentage)
		}
		discountedProduct := domain.ProductDiscount{
			ID:       product.ID,
			SKU:      product.SKU,
			Name:     product.Name,
			Category: product.Category,
			Price: domain.Price{
				Original:           product.Price,
				Final:              finalPrice,
				DiscountPercentage: discountPercentage,
				Currency:           product.Currency,
			},
		}
		discountedProducts = append(discountedProducts, discountedProduct)
	}

	return discountedProducts, pagination, nil
}

func applyDiscount(originalPrice uint32, discountPercentage uint8) uint32 {
	discountAmount := (originalPrice * uint32(discountPercentage)) / 100
	return originalPrice - discountAmount
}
