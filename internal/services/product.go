package services

import (
	"context"
	"fmt"
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"

	"github.com/speps/go-hashids"
)

type ProductService struct {
	Repo            ports.ProductRepository
	DiscountService ports.DiscountService
	Hashids         *hashids.HashID
}

func NewProductService(repo ports.ProductRepository, discountService ports.DiscountService, salt string) ports.ProductService {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 6
	hashid, _ := hashids.NewWithData(hd)
	return &ProductService{Repo: repo, DiscountService: discountService, Hashids: hashid}
}

func (s *ProductService) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	return s.Repo.CreateProduct(ctx, product)
}

func (s ProductService) ListProducts(ctx context.Context, filters map[string]interface{}, pageSize int, nextID, prevID uint32) ([]domain.ProductDiscount, domain.HashedPagination, error) {
	products, pagination, err := s.Repo.ListProducts(ctx, filters, pageSize, nextID, prevID)
	if err != nil {
		return nil, domain.HashedPagination{}, err
	}
	var discountedProducts []domain.ProductDiscount

	for _, product := range products {
		var finalPrice = product.Price
		var discountPercentage *string

		discount, err := s.DiscountService.GetBestDiscount(ctx, product.SKU, product.Category)
		if err != nil || discount.ID == 0 {
			discountPercentage = nil
		} else {
			discountPercentageToString := fmt.Sprintf("%d%%", discount.Percentage)
			discountPercentage = &discountPercentageToString
			finalPrice = applyDiscount(finalPrice, discount.Percentage)
		}
		discountedProduct := domain.ProductDiscount{
			ID:       product.ID,
			SKU:      product.SKU,
			Name:     product.Name,
			Category: product.Category,
			Price: domain.Price{
				Original:           product.Price,
				Final:              finalPrice,
				DiscountPercentage: discountPercentage,
				Currency:           product.Currency,
			},
		}
		discountedProducts = append(discountedProducts, discountedProduct)
	}

	var nextHashID, prevHashID string
	if pagination.Next > 0 {
		nextHashID, _ = s.Hashids.EncodeInt64([]int64{int64(pagination.Next)})
	}
	if pagination.Prev > 0 {
		prevHashID, _ = s.Hashids.EncodeInt64([]int64{int64(pagination.Prev)})
	}

	hashedPagination := domain.HashedPagination{
		Next:     nextHashID,
		Prev:     prevHashID,
		PageSize: pagination.PageSize,
	}

	return discountedProducts, hashedPagination, nil
}

func applyDiscount(originalPrice uint32, discountPercentage uint8) uint32 {
	discountAmount := (originalPrice * uint32(discountPercentage)) / 100
	return originalPrice - discountAmount
}
