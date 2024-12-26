package services

import (
	"context"
	"fmt"
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"
)

type ProductService struct {
	Repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) ports.ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	return s.Repo.CreateProduct(ctx, product)
}

func (s ProductService) ListProducts(ctx context.Context, filters map[string]interface{}, pageSize, page int) ([]domain.ProductDiscount, domain.Pagination, error) {
	products, pagination, err := s.Repo.ListProducts(ctx, filters, pageSize, page)
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	var discountedProducts []domain.ProductDiscount

	for _, product := range products {
		var finalPrice = product.Price
		var discountPercentage *string

		if product.SKU == "000003" {
			discountPercentage = applyDiscount(&finalPrice, 15)
		}
		if product.Category == "boots" {
			discountPercentage = applyDiscount(&finalPrice, 30)
		}

		discountedProduct := domain.ProductDiscount{
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

func applyDiscount(price *uint32, percentage int) *string {
	discount := (*price * uint32(percentage)) / 100
	*price -= discount
	discountString := fmt.Sprintf("%d%%", percentage)
	return &discountString
}
