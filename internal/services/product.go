package services

import (
	"context"
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

func (s ProductService) ListProducts(ctx context.Context, filters map[string]interface{}, pageSize, page int) ([]domain.Product, domain.Pagination, error) {
	products, pagination, err := s.Repo.ListProducts(ctx, filters, pageSize, page)
	// discounted := ApplyDiscounts(products)

	return products, pagination, err
}
