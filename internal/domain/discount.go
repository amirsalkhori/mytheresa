package domain

type Discount struct {
	ID         uint32
	Type       string
	Identifier string
	Percentage uint8
}
