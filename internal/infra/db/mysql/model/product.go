package model

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
