package domain

type Product struct {
	ID         uint32
	SKU        string
	Name       string
	Price      uint32
	Currency   string
	CategoryID uint32
	Category   Category
}

type ProductDiscount struct {
	ID       uint32
	SKU      string
	Name     string
	Category string
	Price    Price
}
