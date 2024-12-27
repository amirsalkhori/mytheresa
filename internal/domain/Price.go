package domain

type Price struct {
	Original           uint32
	Final              uint32
	DiscountPercentage *uint8
	Currency           string
}
