package repository

import (
	"context"
	"log"
	"mytheresa/internal/domain"
	"mytheresa/internal/infra/db/mysql/model"
	"mytheresa/internal/ports"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &ProductRepository{
		DB: db,
	}
}

func (r ProductRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	modelProduct := model.ToModelProduct(product)
	if err := r.DB.Create(&modelProduct).Error; err != nil {
		return domain.Product{}, err
	}
	return model.ToDomainProduct(modelProduct), nil
}

func (r *ProductRepository) ListProducts(ctx context.Context, filters map[string]interface{}, pageSize int, next, prev uint32) ([]domain.Product, domain.Pagination, error) {
	var modelProducts []model.Product
	query := r.DB.Model(&model.Product{})

	// Apply filters
	if category, ok := filters["category"]; ok {
		query.Where("category = ?", category)
	}
	if price, ok := filters["priceLessThan"]; ok {
		query.Where("price <= ?", price)
	}

	if next > 0 {
		query.Where("id > ?", next).Order("id ASC").Limit(pageSize)
	} else if prev > 0 {
		query.Where("id < ?", prev).Order("id DESC").Limit(pageSize)
	} else {
		query.Order("id ASC").Limit(pageSize)
	}

	query.Order("id DESC").Limit(pageSize)
	if err := query.Find(&modelProducts).Error; err != nil {
		log.Printf("Failed to get products: %v", err)
		return nil, domain.Pagination{}, err
	}

	products := model.ToDomainProducts(modelProducts)

	var nextID, prevID uint32
	if len(products) > 0 {
		nextID = products[len(products)-1].ID
		prevID = products[0].ID
	}

	if prev > 0 {
		for i, j := 0, len(products)-1; i < j; i, j = i+1, j-1 {
			products[i], products[j] = products[j], products[i]
		}
	}

	pagination := domain.Pagination{
		Next:     nextID,
		Prev:     prevID,
		PageSize: uint32(pageSize),
	}

	return products, pagination, nil
}
