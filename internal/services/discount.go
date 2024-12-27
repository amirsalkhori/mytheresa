package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	discountKey := s.generateRedisKey(createdDiscount.Identifier)

	discountValue, _ := json.Marshal(createdDiscount)

	err = s.redis.Set(ctx, discountKey, discountValue, 0)
	if err != nil {
		log.Fatal("error during the discount persist into the Redis!")
	}

	return createdDiscount, nil
}

func (s *DiscountService) GetDiscount(ctx context.Context, identifier string) (domain.Discount, error) {
	discountKey := s.generateRedisKey(identifier)
	discountData, err := s.redis.Get(ctx, discountKey)
	if err == nil {
		var discount domain.Discount
		if err := json.Unmarshal([]byte(discountData), &discount); err == nil {
			return discount, nil
		}
		log.Fatal("Log unmarshalling error but proceed to fallback")
	}

	discount, dbErr := s.repo.GetDiscountsBySKUAndCategory(ctx, identifier)
	if dbErr != nil {
		return domain.Discount{}, errors.New("unable to fetch discount from database")
	}

	if discount.ID == 0 {
		return domain.Discount{}, nil
	}
	discountValue, _ := json.Marshal(discount)
	_ = s.redis.Set(ctx, discountKey, discountValue, 0)

	return discount, nil
}

func (s *DiscountService) GetBestDiscount(ctx context.Context, product domain.Product) (domain.Discount, error) {
	productAttributes := []string{product.SKU, product.Category}
	discounts := make([]domain.Discount, 0, len(productAttributes))

	for _, attribute := range productAttributes {
		discount, err := s.getDiscountByAttribute(ctx, attribute)
		if err != nil {
			log.Printf("Error fetching discount for attribute %s: %v", attribute, err)
			continue
		}

		if discount.ID != 0 {
			discounts = append(discounts, discount)
		}
	}

	return s.getHighestDiscount(discounts), nil
}

func (s *DiscountService) getDiscountByAttribute(ctx context.Context, attribute string) (domain.Discount, error) {
	redisKey := s.generateRedisKey(attribute)
	discountData, err := s.redis.Get(ctx, redisKey)
	if err == nil {
		var discount domain.Discount
		if jsonErr := json.Unmarshal([]byte(discountData), &discount); jsonErr == nil {
			return discount, nil
		}
		log.Printf("Failed to unmarshal discount data for key %s: %v", redisKey, err)
	}
	discount, dbErr := s.repo.GetDiscountsBySKUAndCategory(ctx, attribute)
	if dbErr != nil {
		return domain.Discount{}, fmt.Errorf("failed to fetch discount from database for attribute %s: %w", attribute, dbErr)
	}
	if discount.ID != 0 {
		discountValue, _ := json.Marshal(discount)
		if setErr := s.redis.Set(ctx, redisKey, discountValue, 0); setErr != nil {
			log.Printf("Failed to cache discount data for key %s: %v", redisKey, setErr)
		}
	}

	return discount, nil
}

func (s *DiscountService) getHighestDiscount(discounts []domain.Discount) domain.Discount {
	var bestDiscount domain.Discount
	for _, discount := range discounts {
		if discount.Percentage > bestDiscount.Percentage {
			bestDiscount = discount
		}
	}
	return bestDiscount
}

func (s *DiscountService) generateRedisKey(identifier string) string {
	return "discount:type:" + identifier
}
