package app

import (
	"log"
	"mytheresa/configs"
	"mytheresa/internal/handler"
	mysql "mytheresa/internal/infra/db/mysql/config"
	"mytheresa/internal/infra/db/mysql/repository"
	"mytheresa/internal/infra/db/redis"
	"mytheresa/internal/services"

	"github.com/gin-gonic/gin"
)

func StartApplication() {
	cfg := configs.GetConfig()

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

	productService := services.NewProductService(productRepo)
	disocuntService := services.NewDiscountService(discountRepo, redis)

	r := gin.Default()
	productHandler := handler.NewProductHandler(productService)
	discountHandler := handler.NewDiscountHandler(disocuntService)

	r.POST("/products", productHandler.CreateProduct)
	r.GET("/products", productHandler.GetFilteredProducts)
	r.POST("/discounts", discountHandler.CreateDiscount)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
