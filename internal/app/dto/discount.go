package dto

type Discount struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
	Percentage uint8  `json:"percentage"`
}

type DisocuntRoot struct {
	Disocunts []Discount `json:"disocunts"`
}
