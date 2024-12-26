package services

import (
	"context"
	"encoding/json"
	"errors"
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

	discountKey := s.generateRedisKey(createdDiscount.SKU, createdDiscount.Category)

	discountValue, _ := json.Marshal(createdDiscount)

	err = s.redis.Set(ctx, discountKey, discountValue, 0)
	if err != nil {
		log.Fatal("error during the discount persist into the Redis!")
	}

	return createdDiscount, nil
}

func (s *DiscountService) GetDiscount(ctx context.Context, sku, category string) (domain.Discount, error) {
	discountKey := s.generateRedisKey(sku, category)
	discountData, err := s.redis.Get(ctx, discountKey)
	if err == nil {
		var discount domain.Discount
		if err := json.Unmarshal([]byte(discountData), &discount); err == nil {
			return discount, nil
		}
		log.Fatal("Log unmarshalling error but proceed to fallback")
	}

	discount, dbErr := s.repo.GetDiscountsBySKUAndCategory(ctx, sku, category)
	if dbErr != nil {
		return domain.Discount{}, errors.New("unable to fetch discount from Redis and database")
	}

	discountValue, _ := json.Marshal(discount)
	_ = s.redis.Set(ctx, discountKey, discountValue, 0)

	return discount, nil
}

func (s *DiscountService) generateRedisKey(sku, category string) string {
	if sku != "" {
		return "discount:sku:" + sku
	}
	return "discount:category:" + category
}
