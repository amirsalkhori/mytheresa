package domain

type Product struct {
	ID       uint32
	SKU      string
	Name     string
	Category string
	Price    Price
}
