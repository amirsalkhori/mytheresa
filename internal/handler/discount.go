package handler

import (
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"

	"github.com/gin-gonic/gin"

	"net/http"
)

type DiscountHandler struct {
	service ports.DiscountService
}

func NewDiscountHandler(service ports.DiscountService) ports.DiscountHandler {
	return &DiscountHandler{service: service}
}

func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var discount domain.Discount
	if err := c.ShouldBindJSON(&discount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount request"})
		return
	}

	createdDiscount, err := h.service.CreateDiscount(c.Request.Context(), discount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product"})
		return
	}

	c.JSON(http.StatusCreated, createdDiscount)
}

func (h *DiscountHandler) GetDiscount(c *gin.Context) {
	sku := c.Query("sku")
	category := c.Query("category")

	if sku == "" && category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one of SKU or category must be provided"})
		return
	}

	discount, err := h.service.GetDiscount(c.Request.Context(), sku, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch discount"})
		return
	}

	if (discount == domain.Discount{}) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Discount not found"})
		return
	}

	c.JSON(http.StatusOK, discount)
}
