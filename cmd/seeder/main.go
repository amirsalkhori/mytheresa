package main

import (
	"flag"
	"log"
	"mytheresa/configs"
	mysql "mytheresa/internal/infra/db/mysql/config"
	"mytheresa/internal/services"
)

func main() {
	categoriesCount := flag.Int("categories", 10, "Number of categories to create")
	productsCount := flag.Int("products", 100, "Number of products to create")
	discountsCount := flag.Int("discounts", 50, "Number of discounts to create")
	flag.Parse()

	cfg := configs.GetConfig()
	cfg.Mysql.Host = "127.0.0.1"
	cfg.Mysql.Port = 13306
	cfg.Mysql.User = "mytheresa"
	cfg.Mysql.Pass = "mytheresa"
	cfg.Mysql.Name = "mytheresa"

	repo, err := mysql.NewMySQLRepository(&cfg)
	if err != nil {
		log.Fatalf("Error initializing repository: %v", err)
	}

	seederService := services.NewSeederService(repo)

	if err := seederService.SeedCategories(*categoriesCount); err != nil {
		log.Fatalf("Error seeding categories: %v", err)
	}
	if err := seederService.SeedProducts(*productsCount); err != nil {
		log.Fatalf("Error seeding products: %v", err)
	}
	if err := seederService.SeedDiscounts(*discountsCount); err != nil {
		log.Fatalf("Error seeding discounts: %v", err)
	}

	log.Println("Seeding completed successfully!")
}
