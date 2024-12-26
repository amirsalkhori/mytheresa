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
	// Step 1: Generate Redis key
	discountKey := s.generateRedisKey(sku, category)

	// Step 2: Try to fetch from Redis
	discountData, err := s.redis.Get(ctx, discountKey)
	if err == nil {
		var discount domain.Discount
		if err := json.Unmarshal([]byte(discountData), &discount); err == nil {
			return discount, nil
		}
		// Log unmarshalling error but proceed to fallback
	}

	// Step 3: Fallback to database query
	discount, dbErr := s.repo.GetDiscount(ctx, sku, category)
	if dbErr != nil {
		return domain.Discount{}, errors.New("unable to fetch discount from Redis and database")
	}

	// Step 4: Optionally, cache the result in Redis for future requests
	discountValue, _ := json.Marshal(discount)
	_ = s.redis.Set(ctx, discountKey, discountValue, 0) // Log errors but don't block

	return discount, nil
}

func (s *DiscountService) generateRedisKey(sku, category string) string {
	if sku != "" {
		return "discount:sku:" + sku
	}
	return "discount:category:" + category
}

// func (s *DiscountService) GetLargestDiscount(ctx context.Context, sku, category string) (domain.Discount, error) {
// 	var discounts []domain.Discount

// 	// Step 1: Try fetching discounts from Redis
// 	redisKey := s.generateRedisKey(sku, category)
// 	discountData, err := s.redis.Get(ctx, redisKey)
// 	if err == nil && discountData != nil {
// 		if err := json.Unmarshal([]byte(discountData), &discounts); err != nil {
// 			log.Printf("error unmarshalling discounts from Redis: %v", err)
// 		}
// 	}

// 	// Step 2: If Redis fails or has no data, fallback to the database
// 	if len(discounts) == 0 {
// 		dbDiscounts, err := s.repo.GetDiscountsBySKUAndCategory(ctx, sku, category)
// 		if err != nil {
// 			return domain.Discount{}, fmt.Errorf("failed to fetch discounts from database: %v", err)
// 		}
// 		discounts = dbDiscounts
// 	}

// 	// Step 3: Find the largest discount
// 	var largestDiscount domain.Discount
// 	maxPercentage := 0.0
// 	for _, discount := range discounts {
// 		percentage, err := strconv.ParseFloat(discount.Percentage, 64)
// 		if err != nil {
// 			log.Printf("invalid discount percentage for SKU %s: %v", discount.SKU, err)
// 			continue
// 		}
// 		if percentage > maxPercentage {
// 			maxPercentage = percentage
// 			largestDiscount = discount
// 		}
// 	}

// 	// Step 4: Return the largest discount
// 	return largestDiscount, nil
// }
