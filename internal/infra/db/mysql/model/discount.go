package model

import "mytheresa/internal/domain"

type Discount struct {
	ID         uint32 `gorm:"primaryKey;autoIncrement"`
	Type       string `gorm:"type:ENUM('CATEGORY', 'SKU');"`
	Identifier string `gorm:"index;type:varchar(50);not null"`
	Percentage uint8  `gorm:"not null"`
}

func (Discount) TableName() string {
	return "discounts"
}

func ToModelDiscount(domainDiscount domain.Discount) Discount {
	return Discount{
		ID:         domainDiscount.ID,
		Type:       domainDiscount.Type,
		Identifier: domainDiscount.Identifier,
		Percentage: domainDiscount.Percentage,
	}
}

func ToDomainDiscount(modelDiscount Discount) domain.Discount {
	return domain.Discount{
		ID:         modelDiscount.ID,
		Type:       modelDiscount.Type,
		Identifier: modelDiscount.Identifier,
		Percentage: modelDiscount.Percentage,
	}
}
