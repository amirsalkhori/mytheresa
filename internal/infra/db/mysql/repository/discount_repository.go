package repository

import (
	"context"
	"log"
	"mytheresa/internal/domain"
	"mytheresa/internal/infra/db/mysql/model"
	"mytheresa/internal/ports"

	"gorm.io/gorm"
)

type DiscountRepository struct {
	DB *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) ports.DisocuntRepository {
	return &DiscountRepository{
		DB: db,
	}
}

func (r DiscountRepository) CreateDiscount(ctx context.Context, disocunt domain.Discount) (domain.Discount, error) {
	modelDiscount := model.ToModelDiscount(disocunt)
	if err := r.DB.Create(&modelDiscount).Error; err != nil {
		return domain.Discount{}, err
	}
	return model.ToDomainDiscount(modelDiscount), nil
}

func (r DiscountRepository) GetDiscountsBySKUAndCategory(ctx context.Context, SKU, categoryName string) ([]domain.Discount, error) {
	var modelDiscount []model.Discount
	query := r.DB.Model(&model.Discount{}).Preload("Category")

	if categoryName != "" {
		query.Joins("JOIN categories c ON c.id = discounts.category_id")
		query.Where("c.name = ?", categoryName)
	}

	if SKU != "" {
		query.Where("discounts.sku <= ?", SKU)
	}

	if err := query.Find(&modelDiscount).Error; err != nil {
		log.Printf("Failed to get discounts: %v", err)
		return nil, err
	}
	return model.ToDomainDiscounts(modelDiscount), nil
}

func (r DiscountRepository) GetAllDiscounts() ([]domain.Discount, error) {
	var modelDiscounts []model.Discount
	query := r.DB.Model(&model.Discount{})
	if err := query.Find(&modelDiscounts).Error; err != nil {
		log.Printf("Failed to get discounts: %v", err)
		return nil, err
	}
	return model.ToDomainDiscounts(modelDiscounts), nil
}
