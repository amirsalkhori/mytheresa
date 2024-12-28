package handler

import (
	"context"
	"fmt"
	"log"
	"mytheresa/internal/app/dto"
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"
	"strconv"

	"github.com/gin-gonic/gin"

	"net/http"
)

type ProductHandler struct {
	Service ports.ProductService
}

func NewProductHandler(service ports.ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
}

func (h *ProductHandler) CreateProducts(ctx context.Context, productsRoot dto.ProductsRoot) {
	for _, product := range productsRoot.Products {
		productDomain := domain.Product{
			SKU:      product.SKU,
			Name:     product.Name,
			Category: product.Category,
			Price:    uint32(product.Price),
		}
		_, err := h.Service.CreateProduct(ctx, productDomain)
		if err != nil {
			fmt.Println("Could not create product")
		}
	}
	log.Println("Products have been successfully stored DB.")
}

func (h *ProductHandler) GetFilteredProducts(c *gin.Context) {
	queryToFilterMap := map[string]string{
		"category":      "category",
		"priceLessThan": "priceLessThan",
	}

	filters := h.extractFiltersFromQuery(c, queryToFilterMap)

	pageSize, err := parsePaginationParams(c, "pagesize", 5)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lastIDStr := c.Query("lastID")
	var lastIDUint32 uint32
	if lastIDStr == "" {
		lastIDUint32 = 0
	} else {
		lastID, err := strconv.ParseUint(lastIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lastID parameter"})
			return
		}
		lastIDUint32 = uint32(lastID)
	}

	products, pagination, err := h.Service.ListProducts(c.Request.Context(), filters, pageSize, lastIDUint32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products":   products,
		"pagination": pagination,
	})
}

func (h *ProductHandler) extractFiltersFromQuery(c *gin.Context, queryToFilterMap map[string]string) map[string]interface{} {
	filters := make(map[string]interface{})

	for queryKey, filterKey := range queryToFilterMap {
		if value := c.Query(queryKey); value != "" {
			filters[filterKey] = value
		}
	}

	return filters
}

func parsePaginationParams(c *gin.Context, pageSizeKey string, defaultPageSize int) (int, error) {

	pageSize, err := strconv.Atoi(c.DefaultQuery(pageSizeKey, strconv.Itoa(defaultPageSize)))
	if err != nil || pageSize < 1 {
		return 0, fmt.Errorf("invalid %s parameter", pageSizeKey)
	}

	return pageSize, nil
}
