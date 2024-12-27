package repository

import (
	"context"
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

func (r DiscountRepository) GetDiscountsBySKUAndCategory(ctx context.Context, identifier string) (domain.Discount, error) {
	var modelDiscount model.Discount

	// Correct SQL query with explicit field selection
	err := r.DB.Table("discounts").
		Where("discounts.identifier = ? ", identifier).
		Order("discounts.percentage DESC").
		Limit(1).
		Scan(&modelDiscount).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Discount{}, nil
		}
		return domain.Discount{}, err
	}
	return model.ToDomainDiscount(modelDiscount), nil
}
