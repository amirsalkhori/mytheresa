package services_test

import (
	"context"
	"errors"
	"fmt"
	"mytheresa/internal/domain"
	"mytheresa/internal/ports"
	"mytheresa/internal/services"
	"mytheresa/internal/services/test/mocks"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = ginkgo.Describe("ProductService", func() {
	var (
		mockRepo            *mocks.MockProductRepository
		mockDiscountService *mocks.MockDiscountService
		productService      ports.ProductService
		ctx                 context.Context
		mockProduct         domain.Product
		mockProduct2        domain.Product
	)

	ginkgo.BeforeEach(func() {
		mockRepo = new(mocks.MockProductRepository)
		mockDiscountService = new(mocks.MockDiscountService)
		productService = services.NewProductService(mockRepo, mockDiscountService)
		ctx = context.Background()
		mockProduct = domain.Product{
			ID:       1,
			SKU:      "000001",
			Name:     "BV Lean leather ankle boots",
			Category: "boots",
			Price:    89000,
			Currency: "EUR",
		}
		mockProduct2 = domain.Product{
			ID:       2,
			SKU:      "000002",
			Name:     "Nathane leather sneakers",
			Category: "sneakers",
			Price:    88000,
			Currency: "EUR",
		}
	})

	ginkgo.Describe("CreateProduct", func() {
		ginkgo.It("creates a product successfully", func() {
			mockRepo.On("CreateProduct", mock.Anything, mockProduct).Return(mockProduct, nil).Once()

			result, err := productService.CreateProduct(ctx, mockProduct)

			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(result).To(gomega.Equal(mockProduct))

			mockRepo.AssertExpectations(ginkgo.GinkgoT())
		})

		ginkgo.It("returns an error when repository call fails", func() {
			mockRepo.On("CreateProduct", mock.Anything, mockProduct).Return(domain.Product{}, errors.New("repo error")).Once()

			result, err := productService.CreateProduct(ctx, mockProduct)

			gomega.Expect(err).To(gomega.MatchError("repo error"))
			gomega.Expect(result).To(gomega.Equal(domain.Product{}))

			mockRepo.AssertExpectations(ginkgo.GinkgoT())
		})
	})

	ginkgo.Describe("ListProducts", func() {
		ginkgo.It("returns products with discounts applied", func() {
			percentage := "20%"
			fmt.Println("mockProduct2", mockProduct2)
			expectedProducts := []domain.ProductDiscount{
				{
					SKU:      "000001",
					Name:     "BV Lean leather ankle boots",
					Category: "boots",
					Price: domain.Price{
						Original:           89000,
						Final:              71200,
						DiscountPercentage: &percentage,
						Currency:           "EUR",
					},
				},
				{
					SKU:      "000002",
					Name:     "Nathane leather sneakers",
					Category: "sneakers",
					Price: domain.Price{
						Original:           88000,
						Final:              88000,
						DiscountPercentage: nil,
						Currency:           "EUR",
					},
				},
			}
			mockRepo.On("ListProducts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return([]domain.Product{mockProduct, mockProduct2}, domain.Pagination{Page: 1, PageSize: 2}, nil).
				Once()

			mockDiscountService.On("GetBestDiscount", mock.Anything, "000001", "boots").
				Return(domain.Discount{ID: 1, Percentage: 20}, nil).
				Once()

			mockDiscountService.On("GetBestDiscount", mock.Anything, "000002", "sneakers").
				Return(domain.Discount{}, nil). // No discount for the second product
				Once()

			filters := map[string]interface{}{"category": "boots"}
			pageSize := 5
			lastID := uint32(0)

			result, paginationResult, err := productService.ListProducts(ctx, filters, pageSize, lastID)

			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(result[0].Price).To(gomega.Equal(expectedProducts[0].Price))
			gomega.Expect(result[0].Price.Final).To(gomega.Equal(expectedProducts[0].Price.Final))
			gomega.Expect(result[1].Price.Final).To(gomega.Equal(expectedProducts[1].Price.Original))
			gomega.Expect(result[1].Price.DiscountPercentage).To(gomega.Equal(expectedProducts[1].Price.DiscountPercentage))
			gomega.Expect(result[1].Price.DiscountPercentage).To(gomega.BeNil())
			gomega.Expect(paginationResult.Page).To(gomega.Equal(uint32(1)))

			mockRepo.AssertExpectations(ginkgo.GinkgoT())
			mockDiscountService.AssertExpectations(ginkgo.GinkgoT())
		})
	})
})
