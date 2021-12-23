package listing_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/amelendres/go-catalog"
	"github.com/amelendres/go-catalog/listing"
	"github.com/amelendres/go-catalog/pricing"
	"github.com/stretchr/testify/assert"
)

const (
	givenProductsJSON = `
	[
    	{		
			"sku": "000001",
			"name": "BV Lean leather ankle boots",
			"category": "boots",
			"price": 89000
		}, 
		{		
			"sku": "000002",
			"name": "BV Lean leather ankle boots",
			"category": "boots",
			"price": 99000
		}, 
		{	
			"sku": "000003",
			"name": "Ashlington leather ankle boots",
			"category": "boots",
			"price": 71000
		},
		{
			"sku": "000004",
			"name": "Naima embellished suede sandals",
			"category": "sandals",
			"price": 79500
		}, 
		{
			"sku": "000005",
			"name": "Nathane leather sneakers",
			"category": "sneakers",
			"price": 59000
		}
		{
			"sku": "000006",
			"name": "AA sandals",
			"category": "sandals",
			"price": 80000
		}
	]	
	`
)

var (
	givenProductDiscount  = catalog.DiscountPercentage(15)
	givenCategoryDiscount = catalog.DiscountPercentage(30)
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

type DiscountRepoStub struct {
	discounts []catalog.Discount
	wantErr   error
}

func (r *DiscountRepoStub) Find(search catalog.SearchCriteria) (discounts []*catalog.Discount, err error) {
	if r.wantErr != nil {
		return nil, r.wantErr
	}
	return discounts, nil
}

func newDiscountRepoStub(d []catalog.Discount, wantErr error) *DiscountRepoStub {
	return &DiscountRepoStub{d, wantErr}
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
	discounts = append(discounts, catalog.NewCategoryDiscount("boots", 30), catalog.NewProductDiscount("000003", 15))

	products := newFakePaginatedProducts(givenProducts)
	lister := newFakeProductLister(products, discounts, nil, nil)
	discountedProducts := newFakePaginatedDiscountedProducts(givenProducts)

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

		// "default pagination": {
		// 	in:      lister,
		// 	want:    firstPagePaginatedProducts,
		// 	wantErr: nil,
		// },
		// "second page": {
		// 	in: "title (1)",
		// 	want: "title (1)",

		// },
		"Category discount": {
			in:      lister,
			search:  catalog.SearchCriteria{},
			want:    discountedProducts,
			wantErr: nil,
		},
		// "Product discount": {
		// 	in:      bootsLister,
		// 	search:  catalog.SearchCriteria{},
		// 	want:    bootsDiscountedProducts,
		// 	wantErr: nil,
		// },
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
		// "filtered by category sandals":                   {in: "(subtitle)", want: "(subtitle)"},
		// "filtered by price less than":                    {in: "title ( 1 )", want: "title ( 1 )"},
		// "filtered by category boots and price less than": {in: "title ( 1 )", want: "title ( 1 )"},
		// "empty result":                                   {in: "title (1) (2)", wantErr: devom.ErrTitleWithManyVolumes("title (1) (2)")},

		// "invalid pagination":      {in: "title (1) 2", wantErr: devom.ErrInvalidTitle("title (1) 2")},
		// "repository failure":      {in: "title (1) 2", wantErr: devom.ErrInvalidTitle("title (1) 2")},
	}

	for name, tc := range tests {
		got, err := tc.in.List(tc.search)

		if tc.wantErr != nil {
			assert.Error(t, err, name)
			assert.Nil(t, got, name)
			// assert.Equal(t, got, nil, name)
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
	discountRepo := newDiscountRepoStub(discounts, wantPricingErr)
	calculater := newStubPricingCalculater(discountRepo, givenDiscountedPrices, wantPricingErr)
	return listing.NewProductLister(productRepo, calculater)
}

func newFakePaginatedProducts(products []*catalog.Product) *catalog.PaginatedProducts {
	pagination, _ := catalog.NewPagination(5, 0)
	return catalog.NewPaginatedProducts(
		catalog.PaginationMeta{
			Total:      len(products),
			Pagination: *pagination,
		},
		products,
	)
}

func newFakePaginatedDiscountedProducts(products []*catalog.Product) *catalog.PaginatedDiscountedProducts {
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
