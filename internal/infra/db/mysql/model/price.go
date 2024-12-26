package model

import "mytheresa/internal/domain"

type Price struct {
	ID                 uint32  `gorm:"primaryKey;autoIncrement"`
	Original           uint32  `gorm:"not null"`
	Final              uint32  `gorm:"not null"`
	DiscountPercentage *string `gorm:"not null"`
	Currency           string  `gorm:"type:varchar(10);not null"`
	ProductID          uint32  `gorm:"not null"`
}

func (Price) TableName() string {
	return "prices"
}

func ToModelPrice(domainPrice domain.Price) Price {
	return Price{
		ID:                 domainPrice.ID,
		Original:           domainPrice.Original,
		Final:              domainPrice.Final,
		DiscountPercentage: domainPrice.DiscountPercentage,
		Currency:           domainPrice.Currency,
	}
}

func ToDomainPrice(modelPrice Price) domain.Price {
	return domain.Price{
		ID:                 modelPrice.ID,
		Original:           modelPrice.Original,
		Final:              modelPrice.Final,
		DiscountPercentage: modelPrice.DiscountPercentage,
		Currency:           modelPrice.Currency,
	}
}
