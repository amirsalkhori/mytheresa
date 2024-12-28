package services_test

import (
	"context"
	"encoding/json"
	"mytheresa/internal/domain"
	"mytheresa/internal/services"
	"mytheresa/internal/services/test/mocks"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestServices(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Services Suite")
}

var _ = ginkgo.Describe("DiscountService", func() {
	var (
		mockDisocuntRepo   *mocks.MockDiscountRepository
		mockCache          *mocks.MockCache
		discountService    services.DiscountService
		ctx                context.Context
		mockDiscount       domain.Discount
		discountID         uint32
		discountType       string
		discountIdentifier string
		discountPercentage uint8
	)

	ginkgo.BeforeEach(func() {
		mockDisocuntRepo = new(mocks.MockDiscountRepository)
		mockCache = new(mocks.MockCache)
		discountService = services.DiscountService{
			Repo:  mockDisocuntRepo,
			Redis: mockCache,
		}
		ctx = context.Background()
		discountID = 1
		discountType = "category"
		discountIdentifier = "boots"
		discountPercentage = 30
		mockDiscount = domain.Discount{
			ID:         discountID,
			Type:       discountType,
			Identifier: discountIdentifier,
			Percentage: discountPercentage,
		}
	})

	ginkgo.Describe("CreateDiscount", func() {
		ginkgo.It("should create a discount and set it in Redis", func() {
			discount := mockDiscount
			mockDisocuntRepo.On("CreateDiscount", ctx, discount).Return(discount, nil)

			discountData, err := json.Marshal(discount)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			mockCache.On("Set", ctx, "discount_category_boots", discountData, time.Duration(0)).Return(nil)

			createdDiscount, err := discountService.CreateDiscount(ctx, discount)

			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(createdDiscount.ID).To(gomega.Equal(discount.ID))
			gomega.Expect(createdDiscount.Type).To(gomega.Equal(discount.Type))
		})
	})

	ginkgo.Describe("GetBestDiscount", func() {
		ginkgo.It("should return the best discount when both sku and category discounts are found", func() {
			discount := mockDiscount
			mockCache.On("Get", ctx, "discount_category_boots").Return(`{"ID": 6, "Percentage": 30}`, nil)
			mockCache.On("Get", ctx, "discount_sku_").Return(`{"ID": 5, "Percentage": 20}`, nil)

			bestDiscount, err := discountService.GetBestDiscount(ctx, "", discount.Identifier)

			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(bestDiscount.Percentage).To(gomega.Equal(uint8(30)))
		})

		ginkgo.It("should return an empty discount when no discount is found for both sku and category", func() {
			mockCache.On("Get", ctx, "discount_sku_000001").Return("", nil)
			mockCache.On("Get", ctx, "discount_category_hats").Return("", nil)

			bestDiscount, err := discountService.GetBestDiscount(ctx, "000001", "hats")

			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(bestDiscount).To(gomega.Equal(domain.Discount{}))
		})
	})
})
