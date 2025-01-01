package services

import (
	"fmt"
	"log"
	"math/rand"
	mysql "mytheresa/internal/infra/db/mysql/config"
	"time"

	"gorm.io/gorm"
)

type SeederService struct {
	DB *gorm.DB
}

func NewSeederService(repo *mysql.MySQLRepository) *SeederService {
	return &SeederService{
		DB: repo.DB,
	}
}

// Seed predefined categories
func (s *SeederService) SeedCategories(count int) error {
	var categories []map[string]interface{}
	for i := 1; i <= count; i++ {
		categories = append(categories, map[string]interface{}{
			"name": fmt.Sprintf("Category %d", i),
		})
	}

	if err := s.DB.Table("categories").CreateInBatches(categories, 100).Error; err != nil {
		log.Printf("Error seeding categories: %v", err)
		return err
	}

	log.Printf("Successfully seeded 4000 categories")
	return nil
}

// Seed products with formatted SKUs and predefined categories
func (s *SeederService) SeedProducts(count int) error {
	var categoryIDs []int
	if err := s.DB.Table("categories").Pluck("id", &categoryIDs).Error; err != nil {
		log.Printf("Error fetching categories: %v", err)
		return err
	}

	if len(categoryIDs) == 0 {
		log.Println("No categories found. Please add categories first.")
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	var products []map[string]interface{}
	for i := 0; i < count; i++ {
		categoryID := categoryIDs[rand.Intn(len(categoryIDs))]
		products = append(products, map[string]interface{}{
			"sku":         fmt.Sprintf("%06d", i+1), // SKU with leading zeros, e.g., 000001
			"name":        fmt.Sprintf("Product %d", i+1),
			"price":       rand.Intn(100000) + 1000, // Price range: 1000 to 101000
			"category_id": categoryID,
		})
	}

	if err := s.DB.Table("products").CreateInBatches(products, 100).Error; err != nil {
		log.Printf("Error seeding products: %v", err)
		return err
	}

	log.Printf("Successfully seeded %d products", count)
	return nil
}

// Seed discounts linked to SKUs or categories
func (s *SeederService) SeedDiscounts(count int) error {
	var categoryIDs []int
	if err := s.DB.Table("categories").Pluck("id", &categoryIDs).Error; err != nil {
		log.Printf("Error fetching category IDs: %v", err)
		return err
	}

	if len(categoryIDs) == 0 {
		log.Println("No categories found in the database. Cannot seed discounts.")
		return fmt.Errorf("no categories available")
	}

	rand.Seed(time.Now().UnixNano())
	var discounts []map[string]interface{}
	for i := 0; i < count; i++ {
		var sku string
		if rand.Intn(2) == 0 {
			sku = fmt.Sprintf("%06d", rand.Intn(999999)+1) // Random SKU with leading zeros
		}

		categoryID := categoryIDs[rand.Intn(len(categoryIDs))]
		percentage := rand.Intn(41) + 10 // Discount percentage range: 10 to 50

		discounts = append(discounts, map[string]interface{}{
			"sku":         sku,
			"category_id": categoryID,
			"percentage":  percentage,
		})
	}

	if err := s.DB.Table("discounts").CreateInBatches(discounts, 100).Error; err != nil {
		log.Printf("Error seeding discounts: %v", err)
		return err
	}

	log.Printf("Successfully seeded %d discounts", count)
	return nil
}
