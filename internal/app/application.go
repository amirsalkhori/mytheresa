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
	hashIDSalt := cfg.HashIDSalt
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
	productService := services.NewProductService(productRepo, discountService, hashIDSalt)

	productHandler := handler.NewProductHandler(productService, hashIDSalt)

	r := gin.Default()
	r.GET("/v1/products", productHandler.GetFilteredProducts)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
