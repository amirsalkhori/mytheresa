package domain

type Price struct {
	Original           uint32
	Final              uint32
	DiscountPercentage *string
	Currency           string
}
