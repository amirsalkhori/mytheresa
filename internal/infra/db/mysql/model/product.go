package model

import "mytheresa/internal/domain"

type Product struct {
	ID       uint32 `gorm:"primaryKey;autoIncrement"`
	SKU      string `gorm:"index;type:varchar(255);not null"`
	Name     string `gorm:"index;type:varchar(255);not null"`
	Category string `gorm:"index;type:varchar(255);not null"`
	Price    Price
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
		Price:    ToModelPrice(domainProduct.Price),
	}
}

func ToDomainProduct(modelProduct Product) domain.Product {
	return domain.Product{
		ID:       modelProduct.ID,
		SKU:      modelProduct.SKU,
		Name:     modelProduct.Name,
		Category: modelProduct.Category,
		Price:    ToDomainPrice(modelProduct.Price),
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
