package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"
)

type DiscountService struct {
	Repo  ports.DisocuntRepository
	Redis ports.Cache
}

func NewDiscountService(repo ports.DisocuntRepository, redis ports.Cache) ports.DiscountService {
	return &DiscountService{Repo: repo, Redis: redis}
}

func (s *DiscountService) CreateDiscount(ctx context.Context, discount domain.Discount) (domain.Discount, error) {
	createdDiscount, err := s.Repo.CreateDiscount(ctx, discount)
	if err != nil {
		return domain.Discount{}, err
	}

	discountKey := s.generateRedisKey(createdDiscount.Type, createdDiscount.Identifier)

	discountValue, _ := json.Marshal(createdDiscount)

	err = s.Redis.Set(ctx, discountKey, discountValue, 0)
	if err != nil {
		log.Fatal("error during the discount persist into the Redis!")
	}

	return createdDiscount, nil
}

func (s *DiscountService) StoreDiscountsInRedis(ctx context.Context) error {
	discounts, err := s.Repo.GetAllDiscounts()
	if err != nil {
		return err
	}

	for _, discount := range discounts {
		discountKey := s.generateRedisKey(discount.Type, discount.Identifier)

		discountData, err := json.Marshal(discount)
		if err != nil {
			log.Printf("Error marshalling discount: %v", err)
			continue
		}
		err = s.Redis.Set(ctx, discountKey, discountData, 0)
		if err != nil {
			log.Fatal("error during the discount persist into the Redis!")
		}
	}

	log.Println("Discounts have been successfully stored in Redis.")
	return nil
}

func (s *DiscountService) GetBestDiscount(ctx context.Context, sku, category string) (domain.Discount, error) {
	attributes := []struct {
		key, value string
	}{
		{"sku", sku},
		{"category", category},
	}

	discounts := make([]domain.Discount, 0, len(attributes))

	for _, attribute := range attributes {
		discount, err := s.getDiscountByAttribute(ctx, attribute.key, attribute.value)
		if err != nil {
			log.Printf("Error fetching discount for attribute %s: %v", attribute, err)
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

func (s *DiscountService) getDiscountByAttribute(ctx context.Context, key, attribute string) (domain.Discount, error) {
	redisKey := s.generateRedisKey(key, attribute)
	discountData, err := s.Redis.Get(ctx, redisKey)
	if err == nil {
		var discount domain.Discount
		if jsonErr := json.Unmarshal([]byte(discountData), &discount); jsonErr == nil {
			return discount, nil
		}
		log.Printf("Failed to unmarshal discount data for key %s: %v", redisKey, err)
	}

	return domain.Discount{}, nil
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

func (s *DiscountService) generateRedisKey(discountType, identifier string) string {
	return fmt.Sprintf("discount_%s_%s", discountType, identifier)
}
