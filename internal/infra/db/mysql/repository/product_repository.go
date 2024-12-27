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

func (r ProductRepository) ListProducts(ctx context.Context, filters map[string]interface{}, pageSize, page int) ([]domain.Product, domain.Pagination, error) {
	var modelProducts []model.Product
	var totalCount int64
	query := r.DB.Model(&model.Product{})

	if category, ok := filters["category"]; ok {
		query = query.Where("products.category = ?", category)
	}

	if price, ok := filters["priceLessThan"]; ok {
		query = query.Where("products.price <= ?", price)
	}
	if err := query.Count(&totalCount).Error; err != nil {
		log.Fatal("Error while counting the rows", err)
		return nil, domain.Pagination{}, err
	}

	offset := (page - 1) * pageSize
	query = query.Limit(pageSize).Offset(offset)

	err := query.Find(&modelProducts).Error
	if err != nil {
		log.Fatal("Failed to get data from products by filter", err)
		return nil, domain.Pagination{}, err
	}

	pagination := domain.Pagination{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Total:    int32(totalCount),
	}

	return model.ToDomainProducts(modelProducts), pagination, nil
}
