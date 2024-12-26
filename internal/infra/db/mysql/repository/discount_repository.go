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
