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
	"github.com/speps/go-hashids"

	"net/http"
)

type ProductHandler struct {
	Service ports.ProductService
	hashids *hashids.HashID
}

func NewProductHandler(service ports.ProductService, salt string) *ProductHandler {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 6
	hashid, _ := hashids.NewWithData(hd)
	return &ProductHandler{Service: service, hashids: hashid}
}

func (h *ProductHandler) CreateProducts(ctx context.Context, productsRoot dto.ProductsRoot) {
	for _, product := range productsRoot.Products {
		productDomain := domain.Product{
			SKU:  product.SKU,
			Name: product.Name,
			// Category: product.Category,
			Price: uint32(product.Price),
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
	if pageSize > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pagesize can not be more than 5"})
		return
	}
	nextHash := c.Query("next")
	prevHash := c.Query("prev")
	var nextID, prevID uint32
	if nextHash != "" {
		decoded, err := h.hashids.DecodeInt64WithError(nextHash)
		if err != nil || len(decoded) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"einvalid next hash": err.Error()})
			return
		}
		nextID = uint32(decoded[0])
	}

	if prevHash != "" {
		decoded, err := h.hashids.DecodeInt64WithError(prevHash)
		if err != nil || len(decoded) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"invalid prev hash": err.Error()})
			return
		}
		prevID = uint32(decoded[0])
	}
	products, pagination, err := h.Service.ListProducts(c.Request.Context(), filters, pageSize, nextID, prevID)
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
