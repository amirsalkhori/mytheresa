package model

import "mytheresa/internal/domain"

type Discount struct {
	ID         uint32 `gorm:"primaryKey;autoIncrement"`
	SKU        string `gorm:"index;type:varchar(50);not null"`
	Percentage uint8  `gorm:"not null"`
	CategoryID uint32 `gorm:"not null"`
	Category   Category
}

func (Discount) TableName() string {
	return "discounts"
}

func ToModelDiscount(domainDiscount domain.Discount) Discount {
	return Discount{
		ID:         domainDiscount.ID,
		SKU:        domainDiscount.SKU,
		CategoryID: domainDiscount.CategoryID,
		Percentage: domainDiscount.Percentage,
		Category:   ToModelCategory(domainDiscount.Category),
	}
}

func ToDomainDiscount(modelDiscount Discount) domain.Discount {
	return domain.Discount{
		ID:         modelDiscount.ID,
		SKU:        modelDiscount.SKU,
		CategoryID: modelDiscount.CategoryID,
		Percentage: modelDiscount.Percentage,
		Category:   ToDomainCategory(modelDiscount.Category),
	}
}

func ToDomainDiscounts(modelDiscounts []Discount) []domain.Discount {
	if len(modelDiscounts) == 0 {
		return nil
	}
	discounts := make([]domain.Discount, 0, len(modelDiscounts))
	for _, discount := range modelDiscounts {
		discounts = append(discounts, ToDomainDiscount(discount))
	}

	return discounts
}
