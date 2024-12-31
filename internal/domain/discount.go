package domain

type Discount struct {
	ID         uint32
	SKU        string
	Percentage uint8
	CategoryID uint32
	Category   Category
}
