package services

import (
	"fmt"
	"log"
	mysql "mytheresa/internal/infra/db/mysql/config"

	"golang.org/x/exp/rand"
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

func (s *SeederService) SeedCategories(count int) error {
	var categories []map[string]interface{}
	for i := 0; i < count; i++ {
		categories = append(categories, map[string]interface{}{
			"name": fmt.Sprintf("Category %d", i+1),
		})
	}

	if err := s.DB.Table("categories").CreateInBatches(categories, 100).Error; err != nil {
		log.Printf("Error seeding categories: %v", err)
		return err
	}

	log.Printf("Successfully seeded %d categories", count)
	return nil
}

func (s *SeederService) SeedProducts(count int) error {
	var products []map[string]interface{}

	var categoryIDs []int
	if err := s.DB.Table("categories").Pluck("id", &categoryIDs).Error; err != nil {
		log.Printf("Error fetching categories: %v", err)
		return err
	}

	if len(categoryIDs) == 0 {
		log.Println("No categories found. Please add categories first.")
		return nil
	}

	for i := 0; i < count; i++ {
		categoryID := categoryIDs[i%len(categoryIDs)]
		products = append(products, map[string]interface{}{
			"sku":         fmt.Sprintf("SKU-%d", i+1),
			"name":        fmt.Sprintf("Product %d", i+1),
			"price":       1000 + (i % 10000),
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

func (s *SeederService) SeedDiscounts(count int) error {
	var categories []string
	if err := s.DB.Table("categories").Select("name").Find(&categories).Error; err != nil {
		log.Printf("Error fetching categories: %v", err)
		return err
	}

	if len(categories) == 0 {
		log.Println("No categories found in the database. Cannot seed discounts.")
		return fmt.Errorf("no categories available")
	}

	var discounts []map[string]interface{}
	for i := 0; i < count; i++ {
		discountType := "category"
		var identifier string

		if rand.Intn(2) == 0 { 
			discountType = "SKU"
			identifier = fmt.Sprintf("%04d", rand.Intn(9999)+1) // SKU with leading zeros (e.g., 0001)
		} else {
			category := categories[rand.Intn(len(categories))]               // Pick a random category
			identifier = fmt.Sprintf("%s-%04d", category, rand.Intn(9999)+1) // Category-SKU format (e.g., boots-0001)
		}

		percentage := uint8(rand.Intn(41) + 10) // Range: 10 to 50

		discounts = append(discounts, map[string]interface{}{
			"type":       discountType,
			"identifier": identifier,
			"percentage": percentage,
		})
	}

	if err := s.DB.Table("discounts").CreateInBatches(discounts, 100).Error; err != nil {
		log.Printf("Error seeding discounts: %v", err)
		return err
	}

	log.Printf("Successfully seeded %d discounts", count)
	return nil
}
