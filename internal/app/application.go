package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mytheresa/configs"
	"mytheresa/internal/app/dto"
	"mytheresa/internal/handler"
	mysql "mytheresa/internal/infra/db/mysql/config"
	"mytheresa/internal/infra/db/mysql/repository"
	"mytheresa/internal/infra/db/redis"
	"mytheresa/internal/services"
	"os"

	"github.com/gin-gonic/gin"
)

func StartApplication() {
	cfg := configs.GetConfig()
	ctx := context.Background()
	db, err := mysql.NewMySQLRepository(&cfg)
	if err != nil {
		log.Fatal("MySQL error:", err)
	}

	redis, err := redis.NewRedisAdapter(&cfg.Redis)
	if err != nil {
		log.Fatal("Redis error:", err)
	}

	productRepo := repository.NewProductRepository(db.DB)
	discountRepo := repository.NewDiscountRepository(db.DB)

	discountService := services.NewDiscountService(discountRepo, redis)
	productService := services.NewProductService(productRepo, discountService)

	productHandler := handler.NewProductHandler(productService)
	discountHandler := handler.NewDiscountHandler(discountService)
	initKey := "app_initialized"
	initialized, err := redis.Get(ctx, initKey)
	if err != nil && err.Error() != "redis: nil" {
		log.Fatalf("Error checking initialization status: %v", err)
	}
	if initialized == "" {
		productFilePath := "./resources/product.json"
		discountFilePath := "./resources/discount.json"

		// Load products, discounts from JSON
		productsRoot, err := loadProductsFromJSON(ctx, productFilePath)
		if err != nil {
			log.Fatalf("Error loading products from JSON: %v", err)
		}
		discountsRoot, err := loadDiscountsFromJSON(ctx, discountFilePath)
		if err != nil {
			log.Fatalf("Error loading discounts from JSON: %v", err)
		}
		productHandler.CreateProducts(ctx, productsRoot)
		discountHandler.CreateDiscountFromFile(ctx, discountsRoot)
		err = redis.Set(ctx, initKey, "true", 0)
		if err != nil {
			log.Fatalf("Error setting initialization status: %v", err)
		}
	}
	r := gin.Default()
	r.GET("/products", productHandler.GetFilteredProducts)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

func loadProductsFromJSON(ctx context.Context, filePath string) (dto.ProductsRoot, error) {
	select {
	case <-ctx.Done():
		return dto.ProductsRoot{}, fmt.Errorf("context cancelled while loading JSON: %w", ctx.Err())
	default:
	}

	file, err := os.Open(filePath)
	if err != nil {
		return dto.ProductsRoot{}, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return dto.ProductsRoot{}, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var productsRoot dto.ProductsRoot
	if err := json.Unmarshal(byteValue, &productsRoot); err != nil {
		return dto.ProductsRoot{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return productsRoot, nil
}

func loadDiscountsFromJSON(ctx context.Context, filePath string) (dto.DisocuntRoot, error) {
	select {
	case <-ctx.Done():
		return dto.DisocuntRoot{}, fmt.Errorf("context cancelled while loading JSON from discount: %w", ctx.Err())
	default:
	}

	file, err := os.Open(filePath)
	if err != nil {
		return dto.DisocuntRoot{}, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return dto.DisocuntRoot{}, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var discountRoot dto.DisocuntRoot
	if err := json.Unmarshal(byteValue, &discountRoot); err != nil {
		return dto.DisocuntRoot{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return discountRoot, nil
}
