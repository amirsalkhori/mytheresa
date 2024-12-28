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

func (r *ProductRepository) ListProducts(ctx context.Context, filters map[string]interface{}, pageSize int, lastID uint32) ([]domain.Product, domain.Pagination, error) {
	var modelProducts []model.Product
	query := r.DB.Model(&model.Product{})

	if category, ok := filters["category"]; ok {
		query.Where("category = ?", category)
	}
	if price, ok := filters["priceLessThan"]; ok {
		query.Where("price <= ?", price)
	}
	if lastID > 0 {
		query.Where("id > ?", lastID)
	}

	query = query.Order("id ASC").Limit(pageSize)
	var products []domain.Product
	if err := query.Find(&modelProducts).Error; err != nil {
		log.Printf("Failed to get products: %v", err)
		return nil, domain.Pagination{}, err
	}

	products = model.ToDomainProducts(modelProducts)

	var nextID uint32
	if len(products) > 0 {
		nextID = products[len(products)-1].ID
	}
	pagination := domain.Pagination{
		PageSize: uint32(pageSize),
		Page:     nextID,
	}

	return products, pagination, nil
}
