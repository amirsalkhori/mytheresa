package domain

type Price struct {
	ID                 uint32
	Original           uint32
	Final              uint32
	DiscountPercentage uint8
	Currency           string
}
