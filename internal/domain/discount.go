package domain

type Discount struct {
	ID         uint32
	Category   string
	SKU        string
	Percentage uint8
}
