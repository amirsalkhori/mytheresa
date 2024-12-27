package model

import "mytheresa/internal/domain"

var Currency = "EUR"

type Product struct {
	ID       uint32 `gorm:"primaryKey;autoIncrement"`
	SKU      string `gorm:"index;type:varchar(255);not null"`
	Name     string `gorm:"index;type:varchar(255);not null"`
	Category string `gorm:"index;type:varchar(255);not null"`
	Price    uint32 `gorm:"index;not null"`
}

func (Product) TableName() string {
	return "products"
}

func ToModelProduct(domainProduct domain.Product) Product {
	return Product{
		ID:       domainProduct.ID,
		SKU:      domainProduct.SKU,
		Name:     domainProduct.Name,
		Category: domainProduct.Category,
		Price:    domainProduct.Price,
	}
}

func ToDomainProduct(modelProduct Product) domain.Product {
	return domain.Product{
		ID:       modelProduct.ID,
		SKU:      modelProduct.SKU,
		Name:     modelProduct.Name,
		Category: modelProduct.Category,
		Currency: Currency,
		Price:    modelProduct.Price,
	}
}

func ToDomainProducts(modelProducts []Product) []domain.Product {
	if len(modelProducts) == 0 {
		return []domain.Product{}
	}

	products := make([]domain.Product, 0, len(modelProducts))

	for _, product := range modelProducts {
		products = append(products, ToDomainProduct(product))
	}

	return products
}
