package handler

import (
	"context"
	"fmt"
	"log"
	"mytheresa/internal/app/dto"
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

func (h *DiscountHandler) CreateDiscountFromFile(ctx context.Context, discountsRoot dto.DisocuntRoot) {
	for _, discount := range discountsRoot.Disocunts {
		discountDomain := domain.Discount{
			Type:       discount.Type,
			Identifier: discount.Identifier,
			Percentage: discount.Percentage,
		}
		_, err := h.service.CreateDiscount(ctx, discountDomain)
		if err != nil {
			fmt.Println("Could not create discount")
		}
	}
	log.Println("Discounts have been successfully stored in DB.")
}
