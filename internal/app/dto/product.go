package dto

type Product struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

type ProductsRoot struct {
	Products []Product `json:"products"`
}
