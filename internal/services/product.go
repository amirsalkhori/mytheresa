package services

import (
	"context"
	"fmt"
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"
	"strconv"
	"strings"
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

func (s ProductService) ListProducts(ctx context.Context, filters map[string]interface{}, pageSize, page int) ([]domain.ProductDiscount, domain.Pagination, error) {
	products, pagination, err := s.Repo.ListProducts(ctx, filters, pageSize, page)
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	var discountedProducts []domain.ProductDiscount

	for _, product := range products {
		var finalPrice = product.Price
		var discountPercentage *string

		discount, err := s.DiscountService.GetDiscount(ctx, product.SKU, product.Category)
		if err != nil || discount.ID == 0 {
			discountPercentage = nil
		} else {
			discountPercentage = &discount.Percentage
			// Convert the discount percentage string to an integer and apply the discount
			percentage, err := parseDiscountPercentage(*discountPercentage)
			if err != nil {
				return nil, domain.Pagination{}, fmt.Errorf("invalid discount percentage: %v", err)
			}
			// Apply the discount to the price
			finalPrice = applyDiscount(finalPrice, percentage)
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

func applyDiscount(originalPrice uint32, discountPercentage int) uint32 {
	discountAmount := (originalPrice * uint32(discountPercentage)) / 100
	return originalPrice - discountAmount
}

func parseDiscountPercentage(discount string) (int, error) {
	// Trim any whitespace and remove the "%" symbol
	discount = strings.TrimSpace(discount)
	if len(discount) == 0 || discount[len(discount)-1] != '%' {
		return 0, fmt.Errorf("invalid discount format: missing '%' symbol")
	}

	// Remove the "%" and convert the remaining part to an integer
	percentageStr := discount[:len(discount)-1]
	percentage, err := strconv.Atoi(percentageStr)
	if err != nil {
		return 0, fmt.Errorf("invalid discount value: %v", err)
	}

	return percentage, nil
}
