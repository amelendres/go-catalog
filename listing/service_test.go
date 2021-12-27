package listing_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/amelendres/go-catalog/catalog"
	"github.com/amelendres/go-catalog/listing"
	"github.com/amelendres/go-catalog/pricing"
	"github.com/amelendres/go-catalog/testing/stub"
	"github.com/stretchr/testify/assert"
)

var (
	givenCategoryDiscount = catalog.DiscountPercentage(30)
	givenProductDiscount  = catalog.DiscountPercentage(15)
	givenProducts         = []*catalog.Product{
		catalog.NewProduct("000001", "BV Lean leather ankle boots", "boots", 89000),
		catalog.NewProduct("000002", "BV Lean leather ankle boots", "boots", 99000),
		catalog.NewProduct("000003", "Ashlington leather ankle boots", "boots", 71000),
		catalog.NewProduct("000004", "Naima embellished suede sandals", "sandals", 79500),
		catalog.NewProduct("000005", "Nathane leather sneakers", "sneakers", 59000),
	}
	givenDiscountedPrices = map[string]*catalog.DiscountedPrice{
		"000001": catalog.NewDiscountedPrice(catalog.Price(89000), &givenCategoryDiscount),
		"000002": catalog.NewDiscountedPrice(catalog.Price(99000), &givenCategoryDiscount),
		"000003": catalog.NewDiscountedPrice(catalog.Price(71000), &givenCategoryDiscount),
		"000004": catalog.NewDiscountedPrice(catalog.Price(79500), nil),
		"000005": catalog.NewDiscountedPrice(catalog.Price(59000), nil),
	}
)

type ProductRepoStub struct {
	products *catalog.PaginatedProducts
	wantErr  error
}

func (r *ProductRepoStub) List(search catalog.SearchCriteria) (products *catalog.PaginatedProducts, err error) {
	if r.wantErr != nil {
		return nil, r.wantErr
	}
	return r.products, nil
}

func newProductRepoStub(p *catalog.PaginatedProducts, wantErr error) *ProductRepoStub {
	return &ProductRepoStub{p, wantErr}
}

type StubPricingCalculater struct {
	repository       catalog.DiscountRepository
	discountedPrices map[string]*catalog.DiscountedPrice
	wantErr          error
}

func newStubPricingCalculater(
	repo catalog.DiscountRepository,
	prices map[string]*catalog.DiscountedPrice,
	wantErr error,
) pricing.Calculater {
	return &StubPricingCalculater{repository: repo, discountedPrices: prices, wantErr: wantErr}
}

func (s StubPricingCalculater) Calculate(p catalog.Product) (*catalog.DiscountedPrice, error) {
	if s.wantErr != nil {
		return nil, s.wantErr
	}
	return s.discountedPrices[string(p.SKU)], nil
}

func TestProductLister_List(t *testing.T) {
	var discounts []catalog.Discount
	discounts = append(
		discounts,
		catalog.NewCategoryDiscount("boots", givenCategoryDiscount),
		catalog.NewProductDiscount("000003", givenProductDiscount),
	)

	products := newPaginatedProducts(givenProducts)
	lister := newFakeProductLister(products, discounts, nil, nil)
	discountedProducts := newPaginatedDiscountedProducts(givenProducts)

	productRepoError := errors.New("fails product repository")
	listerWithProductErr := newFakeProductLister(products, discounts, productRepoError, nil)

	discountRepoError := errors.New("fails discount repository")
	listerWithPricingErr := newFakeProductLister(products, discounts, nil, discountRepoError)

	tests := map[string]struct {
		in      listing.ProductLister
		search  catalog.SearchCriteria
		want    *catalog.PaginatedDiscountedProducts
		wantErr error
	}{
		"Category discount": {
			in:      lister,
			search:  catalog.SearchCriteria{},
			want:    discountedProducts,
			wantErr: nil,
		},
		"Product repository error": {
			in:      listerWithProductErr,
			search:  catalog.SearchCriteria{},
			want:    nil,
			wantErr: productRepoError,
		},
		"Pricing calculater error": {
			in:      listerWithPricingErr,
			search:  catalog.SearchCriteria{},
			want:    nil,
			wantErr: discountRepoError,
		},
	}

	for name, tc := range tests {
		got, err := tc.in.List(tc.search)

		if tc.wantErr != nil {
			assert.Error(t, err, name)
			assert.Nil(t, got, name)
			continue
		}
		assert.NoError(t, err, name)
		assert.Equal(t, tc.want, got, name)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.want, got)
		}
	}
}

func newFakeProductLister(
	products *catalog.PaginatedProducts,
	discounts []catalog.Discount,
	wantProductErr error,
	wantPricingErr error,
) listing.ProductLister {
	productRepo := newProductRepoStub(products, wantProductErr)
	discountRepo := stub.NewStubDiscountRepo(discounts, wantPricingErr)
	calculater := newStubPricingCalculater(discountRepo, givenDiscountedPrices, wantPricingErr)
	return listing.NewProductLister(productRepo, calculater)
}

func newPaginatedProducts(products []*catalog.Product) *catalog.PaginatedProducts {
	pagination, _ := catalog.NewPagination(5, 0)
	return catalog.NewPaginatedProducts(
		catalog.PaginationMeta{
			Total:      len(products),
			Pagination: *pagination,
		},
		products,
	)
}

func newPaginatedDiscountedProducts(products []*catalog.Product) *catalog.PaginatedDiscountedProducts {
	var discountedProducts []*catalog.DiscountedProduct
	for _, p := range products {
		discountedProducts = append(discountedProducts, catalog.NewDiscountedProduct(
			p.SKU,
			p.Name,
			p.Category,
			*givenDiscountedPrices[string(p.SKU)],
		))
	}

	pagination, _ := catalog.NewPagination(5, 0)
	return catalog.NewPaginatedDiscountedProducts(
		catalog.PaginationMeta{
			Total:      len(products),
			Pagination: *pagination,
		},
		discountedProducts,
	)
}
