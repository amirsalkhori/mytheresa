package services_test

import (
	"context"
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
		mockProducts        []domain.Product
		mockCategories      []domain.Category
	)

	ginkgo.BeforeEach(func() {
		mockRepo = new(mocks.MockProductRepository)
		mockDiscountService = new(mocks.MockDiscountService)
		productService = services.NewProductService(mockRepo, mockDiscountService, "mytheresa-salt-value")
		ctx = context.Background()

		mockCategories = []domain.Category{
			{ID: 1, Name: "boots"},
			{ID: 2, Name: "sneakers"},
		}

		mockProducts = []domain.Product{
			{
				ID:         1,
				SKU:        "000001",
				Name:       "BV Lean leather ankle boots",
				Category:   mockCategories[0],
				CategoryID: 1,
				Price:      89000,
				Currency:   "EUR",
			},
			{
				ID:         2,
				SKU:        "000002",
				Name:       "Nathane leather sneakers",
				Category:   mockCategories[1],
				CategoryID: 2,
				Price:      88000,
				Currency:   "EUR",
			},
		}
	})

	ginkgo.Describe("ListProducts", func() {
		ginkgo.It("returns products with discounts applied and correct pagination", func() {
			percentage := "20%"
			expectedProducts := []domain.ProductDiscount{
				{
					ID:       1,
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
					ID:       2,
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

			paginationResult := domain.HashedPagination{
				Next:     "WYxok8",
				Prev:     "Ra2dkD",
				PageSize: 5,
			}

			mockRepo.On("ListProducts", ctx, mock.Anything, 5, uint32(2), uint32(5)).
				Return(mockProducts, domain.Pagination{
					Next:     2,
					Prev:     5,
					PageSize: 5,
				}, nil).Once()

			mockDiscountService.On("GetBestDiscount", ctx, "000001", "boots").
				Return(domain.Discount{ID: 1, Percentage: 20}, nil).Once()

			mockDiscountService.On("GetBestDiscount", ctx, "000002", "sneakers").
				Return(domain.Discount{}, nil).Once()

			filters := map[string]interface{}{"category": "boots"}
			pageSize := 5
			nextID := uint32(2)
			prevID := uint32(5)

			result, pagination, err := productService.ListProducts(ctx, filters, pageSize, nextID, prevID)

			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(result).To(gomega.Equal(expectedProducts))
			gomega.Expect(pagination).To(gomega.Equal(paginationResult))
			gomega.Expect(result[0].Price).To(gomega.Equal(expectedProducts[0].Price))
			gomega.Expect(result[0].Price.Final).To(gomega.Equal(expectedProducts[0].Price.Final))
			gomega.Expect(result[1].Price.Final).To(gomega.Equal(expectedProducts[1].Price.Original))
			gomega.Expect(result[1].Price.DiscountPercentage).To(gomega.Equal(expectedProducts[1].Price.DiscountPercentage))
			gomega.Expect(result[1].Price.DiscountPercentage).To(gomega.BeNil())

			mockRepo.AssertExpectations(ginkgo.GinkgoT())
			mockDiscountService.AssertExpectations(ginkgo.GinkgoT())
		})
	})
})
