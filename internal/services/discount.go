package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"
	"time"
)

const (
	REDIS_EXPIRE = 3600
)

type DiscountService struct {
	Repo  ports.DisocuntRepository
	Redis ports.Cache
}

func NewDiscountService(repo ports.DisocuntRepository, redis ports.Cache) ports.DiscountService {
	return &DiscountService{Repo: repo, Redis: redis}
}

func (s *DiscountService) GetBestDiscount(ctx context.Context, sku, category string) (domain.Discount, error) {
	attributes := []struct {
		key, value string
	}{
		{"sku", sku},
		{"category", category},
	}

	var discounts []domain.Discount

	for _, attribute := range attributes {
		discount, err := s.getDiscountByAttribute(ctx, attribute.key, attribute.value)
		if err != nil {
			log.Printf("Error fetching discount for %s=%s: %v", attribute.key, attribute.value, err)
			continue
		}

		if discount.ID != 0 {
			discounts = append(discounts, discount)
		}
	}

	if len(discounts) == 0 {
		return domain.Discount{}, nil
	}

	return s.getHighestDiscount(discounts), nil
}

func (s *DiscountService) getDiscountByAttribute(ctx context.Context, key, value string) (domain.Discount, error) {
	redisKey := s.generateRedisKey(key, value)
	discountData, err := s.Redis.Get(ctx, redisKey)
	if err == nil && discountData != "" {
		var discount domain.Discount
		if jsonErr := json.Unmarshal([]byte(discountData), &discount); jsonErr == nil {
			return discount, nil
		}
		log.Printf("Failed to unmarshal cached discount data for key %s: %v", redisKey)
	}

	sku, category := "", ""
	if key == "sku" {
		sku = value
	} else {
		category = value
	}

	discounts, err := s.Repo.GetDiscountsBySKUAndCategory(ctx, sku, category)
	if err != nil {
		log.Printf("Error fetching discounts from DB for %s=%s: %v", key, value, err)
		return domain.Discount{}, err
	}

	if len(discounts) == 0 {
		return domain.Discount{}, nil
	}

	highestDiscount := s.getHighestDiscount(discounts)

	if discountData, err := json.Marshal(highestDiscount); err == nil {
		cacheErr := s.Redis.Set(ctx, redisKey, discountData, time.Duration(REDIS_EXPIRE)*time.Second)
		if cacheErr != nil {
			log.Printf("Error caching discount data for key %s: %v", redisKey, cacheErr)
		}
	} else {
		log.Printf("Error marshalling discount for caching: %v", err)
	}

	return highestDiscount, nil
}

func (s *DiscountService) getHighestDiscount(discounts []domain.Discount) domain.Discount {
	if len(discounts) == 0 {
		return domain.Discount{}
	}

	bestDiscount := discounts[0]
	for _, discount := range discounts[1:] {
		if discount.Percentage > bestDiscount.Percentage {
			bestDiscount = discount
		}
	}
	return bestDiscount
}

func (s *DiscountService) generateRedisKey(discountType, identifier string) string {
	return fmt.Sprintf("discount:%s:%s", discountType, identifier)
}
