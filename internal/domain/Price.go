package domain

type Price struct {
	ID                 uint32
	Original           uint32
	Final              uint32
	DiscountPercentage *string
	Currency           string
}
