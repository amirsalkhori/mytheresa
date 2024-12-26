package domain

type Product struct {
	ID       uint32
	SKU      string
	Name     string
	Category string
	Price    uint32
	Currency string
}

type ProductDiscount struct {
	SKU      string
	Name     string
	Category string
	Price    Price
}
