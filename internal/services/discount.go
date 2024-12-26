package services

import (
	"context"
	"encoding/json"
	"log"
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"
)

type DiscountService struct {
	repo  ports.DisocuntRepository
	redis ports.Cache
}

func NewDiscountService(repo ports.DisocuntRepository, redis ports.Cache) ports.DiscountService {
	return &DiscountService{repo: repo, redis: redis}
}

func (s *DiscountService) CreateDiscount(ctx context.Context, discount domain.Discount) (domain.Discount, error) {
	createdDiscount, err := s.repo.CreateDiscount(ctx, discount)
	if err != nil {
		return domain.Discount{}, err
	}

	discountKey := s.generateRedisKey(createdDiscount)
	discountValue, _ := json.Marshal(createdDiscount)

	err = s.redis.Set(ctx, discountKey, discountValue, 0)
	if err != nil {
		log.Fatal("error during the discount persist into the Redis!")
	}

	return createdDiscount, nil
}

func (s *DiscountService) generateRedisKey(discount domain.Discount) string {
	if discount.SKU != "" {
		return "discount:sku:" + discount.SKU
	}
	return "discount:category:" + discount.Category
}
