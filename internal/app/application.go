package app

import (
	"log"
	"mytheresa/configs"
	"mytheresa/internal/handler"
	mysql "mytheresa/internal/infra/db/mysql/config"
	"mytheresa/internal/infra/db/mysql/repository"
	"mytheresa/internal/services"

	"github.com/gin-gonic/gin"
)

func StartApplication() {
	cfg := configs.GetConfig()

	db, err := mysql.NewMySQLRepository(&cfg)
	if err != nil {
		log.Fatal("MySQL error:", err)
	}

	productRepo := repository.NewProductRepository(db.DB)
	productService := services.NewProductService(productRepo)

	r := gin.Default()
	productHandler := handler.NewProductHandler(productService)

	r.POST("/products", productHandler.CreateProduct)
	r.GET("/products", productHandler.GetFilteredProducts)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
