package services

// import (
// 	"mytheresa/internal/domain"
// )

// func ApplyDiscounts(products []domain.Product) []models.Product {
// 	for i, product := range products {
// 		originalPrice := product.Price.Original
// 		var finalPrice int
// 		var discountPercentage string

// 		// Apply discounts
// 		if product.Category == "boots" {
// 			finalPrice = originalPrice * 70 / 100
// 			discountPercentage = "30%"
// 		}
// 		if product.SKU == "000003" {
// 			skuDiscountPrice := originalPrice * 85 / 100
// 			if finalPrice == 0 || skuDiscountPrice < finalPrice {
// 				finalPrice = skuDiscountPrice
// 				discountPercentage = "15%"
// 			}
// 		}

// 		// If no discount applied
// 		if finalPrice == 0 {
// 			finalPrice = originalPrice
// 			product.Price.DiscountPercentage = nil
// 		} else {
// 			product.Price.DiscountPercentage = &discountPercentage
// 		}

// 		product.Price.Final = finalPrice
// 		products[i] = product
// 	}
// 	return products
// }
