package model

import "mytheresa/internal/domain"

type Discount struct {
	ID         uint32 `gorm:"primaryKey;autoIncrement"`
	SKU        string `gorm:"index;type:varchar(255);not null"`
	Category   string `gorm:"index;type:varchar(255);not null"`
	Percentage string `gorm:"index;type:varchar(10);not null"`
}

func (Discount) TableName() string {
	return "discounts"
}

func ToModelDiscount(domainDiscount domain.Discount) Discount {
	return Discount{
		ID:         domainDiscount.ID,
		SKU:        domainDiscount.SKU,
		Category:   domainDiscount.Category,
		Percentage: domainDiscount.Percentage,
	}
}

func ToDomainDiscount(modelDiscount Discount) domain.Discount {
	return domain.Discount{
		ID:         modelDiscount.ID,
		SKU:        modelDiscount.SKU,
		Category:   modelDiscount.Category,
		Percentage: modelDiscount.Percentage,
	}
}
