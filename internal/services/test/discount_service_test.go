package services_test

import (
	"context"
	"mytheresa/internal/domain"
	"mytheresa/internal/services"
	"mytheresa/internal/services/test/mocks"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestServices(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Services Suite")
}

var _ = ginkgo.Describe("DiscountService", func() {
	var (
		mockDisocuntRepo *mocks.MockDiscountRepository
		mockCache        *mocks.MockCache
		discountService  services.DiscountService
		ctx              context.Context
		mockDiscount1    domain.Discount
		mockDiscount2    domain.Discount
		mockCategory     domain.Category
	)

	ginkgo.BeforeEach(func() {
		mockDisocuntRepo = new(mocks.MockDiscountRepository)
		mockCache = new(mocks.MockCache)
		discountService = services.DiscountService{
			Repo:  mockDisocuntRepo,
			Redis: mockCache,
		}
		ctx = context.Background()
		mockCategory = domain.Category{
			ID:   1,
			Name: "boots",
		}
		mockDiscount1 = domain.Discount{
			ID:         1,
			Category:   mockCategory,
			CategoryID: 1,
			Percentage: 30,
		}
		mockDiscount2 = domain.Discount{
			ID:         2,
			SKU:        "000001",
			Percentage: 40,
		}
	})

	ginkgo.Describe("GetBestDiscount", func() {
		ginkgo.It("should return the best discount when both sku and category discounts are found", func() {
			categoryDiscount := `{"ID": 1, "Percentage": 30, "CategoryID": 1, "Category": {"Name": "boots"}}`
			skuDiscount := `{"ID": 2, "Percentage": 40, "SKU": "000001"}`

			mockCache.On("Get", ctx, "discount:category:boots").Return(categoryDiscount, nil).Once()
			mockCache.On("Get", ctx, "discount:sku:000001").Return(skuDiscount, nil).Once()

			bestDiscount, err := discountService.GetBestDiscount(ctx, mockDiscount2.SKU, mockDiscount1.Category.Name)

			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Expected no error from GetBestDiscount")
			gomega.Expect(bestDiscount.Percentage).To(gomega.Equal(uint8(40)), "Expected the best discount percentage to be 40")
			gomega.Expect(bestDiscount.ID).To(gomega.Equal(uint32(2)), "Expected the best discount ID to be 2")
			gomega.Expect(bestDiscount.SKU).To(gomega.Equal(mockDiscount2.SKU), "Expected the SKU to match mockDiscount2.SKU")
		})

		ginkgo.It("should return an empty discount when no discount is found for both sku and category", func() {

			categoryDiscount := `{}`
			skuDiscount := `{}`

			mockCache.On("Get", ctx, "discount:category:hats").Return(categoryDiscount, nil).Once()
			mockCache.On("Get", ctx, "discount:sku:000002").Return(skuDiscount, nil).Once()

			bestDiscount, err := discountService.GetBestDiscount(ctx, "000002", "hats")

			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(bestDiscount).To(gomega.Equal(domain.Discount{}))
			gomega.Expect(bestDiscount.ID).To(gomega.Equal(uint32(0)))
			gomega.Expect(bestDiscount.ID).To(gomega.Equal(uint32(0)))
		})
	})
})
