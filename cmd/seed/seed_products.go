package main

import (
	"fmt"
	"log"
)

type Product struct {
	SKU      string
	Name     string
	Category string
	Price    int
}

func main() {
	// Load configuration
	// cfg := configs.GetConfig()

	// Initialize MySQL repository
	// repo, err := mysql.NewMySQLRepository(&cfg)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize MySQL repository: %v", err)
	// }
	// create category
	// create discount
	// Prepare bulk insert
	const batchSize = 1000 // Number of products per batch
	products := make([]Product, 100)
	for i := 1; i <= 100; i++ {
		products[i-1] = Product{
			SKU:      fmt.Sprintf("000%03d", i),
			Name:     fmt.Sprintf("Product %03d", i),
			Category: "boots",
			Price:    89000,
		}
	}

	// Insert products in batches
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}

		// Create a batch of products
		batch := products[i:end]

		// Build bulk insert query
		values := []string{}
		args := []interface{}{}
		for _, product := range batch {
			values = append(values, "(?, ?, ?, ?)")
			args = append(args, product.SKU, product.Name, product.Category, product.Price)
		}

		// query := fmt.Sprintf(
		// 	`INSERT INTO products (sku, name, category, price) VALUES %s`,
		// 	strings.Join(values, ","),
		// )

		// Execute bulk insert
		// _, err := repo.DB.Exec(query, args...)
		// if err != nil {
		// 	log.Printf("[ERROR] Bulk insert failed for batch starting at %d: %v", i, err)
		// } else {
		// 	log.Printf("[INFO] Successfully inserted batch starting at %d", i)
		// }
	}

	log.Println("[INFO] Finished seeding 100 products.")
}
