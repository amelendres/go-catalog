package rest_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/amelendres/go-catalog"
	"github.com/amelendres/go-catalog/http/rest"
	"github.com/amelendres/go-catalog/listing"
	"github.com/amelendres/go-catalog/pricing"
	"github.com/amelendres/go-catalog/storage/inmem"
	"github.com/amelendres/go-catalog/testing/mother"
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
		catalog.NewProduct("000006", "AA hat", "hats", 72000),
	}
	givenDiscountedPrices = map[string]*catalog.DiscountedPrice{
		"000001": catalog.NewDiscountedPrice(catalog.Price(89000), &givenCategoryDiscount),
		"000002": catalog.NewDiscountedPrice(catalog.Price(99000), &givenCategoryDiscount),
		"000003": catalog.NewDiscountedPrice(catalog.Price(71000), &givenCategoryDiscount),
		"000004": catalog.NewDiscountedPrice(catalog.Price(79500), nil),
		"000005": catalog.NewDiscountedPrice(catalog.Price(59000), nil),
		"000006": catalog.NewDiscountedPrice(catalog.Price(72000), nil),
	}
)

func TestCatalogServer_listProducts(t *testing.T) {
	var discounts []catalog.Discount
	discounts = append(
		discounts,
		catalog.NewCategoryDiscount("boots", givenCategoryDiscount),
		catalog.NewProductDiscount("000003", givenProductDiscount),
	)

	productRepo := inmem.NewProductRepo(givenProducts)
	discountRepo := inmem.NewDiscountRepo(discounts)
	pricingCalculater := pricing.NewCalculater(discountRepo)
	productLister := listing.NewProductLister(productRepo, pricingCalculater)
	catalogService := rest.NewCatalogServer(productLister)

	firstPage, _ := catalog.NewPagination(5, 0)
	secondPage, _ := catalog.NewPagination(5, 5)

	searchWithoutFilter := catalog.NewSearchCriteria(firstPage, nil)
	productsWithoutFilter, _ := productRepo.List(searchWithoutFilter)
	listWithoutFilter := mother.NewPaginatedDiscountedProducts(productsWithoutFilter.Items(), givenDiscountedPrices, *firstPage, 6)

	searchSecondPage := catalog.NewSearchCriteria(secondPage, nil)
	productsSecondPage, _ := productRepo.List(searchSecondPage)
	listSecondPage := mother.NewPaginatedDiscountedProducts(productsSecondPage.Items(), givenDiscountedPrices, *secondPage, 6)

	searchBoots := catalog.NewSearchCriteria(firstPage, []catalog.Filter{catalog.NewCategoryFilter("boots")})
	bootsProducts, _ := productRepo.List(searchBoots)
	listBoots := mother.NewPaginatedDiscountedProducts(bootsProducts.Items(), givenDiscountedPrices, *firstPage, 3)

	searchByPrice := catalog.NewSearchCriteria(firstPage, []catalog.Filter{catalog.NewPriceLessThanFilter(71000)})
	priceProducts, _ := productRepo.List(searchByPrice)
	listByPrice := mother.NewPaginatedDiscountedProducts(priceProducts.Items(), givenDiscountedPrices, *firstPage, 2)

	listEmpty := mother.NewPaginatedDiscountedProducts(nil, nil, *firstPage, 0)

	response := httptest.NewRecorder()

	tests := map[string]struct {
		search catalog.SearchCriteria
		want   *catalog.PaginatedDiscountedProducts
		length int
		status int
	}{
		"Without filters": {
			search: searchWithoutFilter,
			want:   listWithoutFilter,
			length: 5,
			status: 200,
		},
		"Second page": {
			search: searchSecondPage,
			want:   listSecondPage,
			length: 1,
			status: 200,
		},
		"Filtered by category boots": {
			search: searchBoots,
			want:   listBoots,
			length: 3,
			status: 200,
		},
		"Filtered by priceLessThan": {
			search: searchByPrice,
			want:   listByPrice,
			length: 2,
			status: 200,
		},
		"Doesn't match filter": {
			search: catalog.NewSearchCriteria(firstPage, []catalog.Filter{catalog.NewCategoryFilter("category-not-found")}),
			want:   listEmpty,
			length: 0,
			status: 200,
		},
	}

	for name, tc := range tests {
		catalogService.ServeHTTP(response, newListProductsRequestFromSearchCriteria(t, tc.search))

		assert.Equal(t, tc.status, response.Code, name)
		got := newPaginatedDiscountedProductsFromJSON(t, response.Body)

		assert.Equal(t, tc.want, got, name)
		assert.Equal(t, tc.length, len(got.Items()))
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.want, got)
		}
	}
}

func TestCatalogServer_listProducts_WithInvalidRequest(t *testing.T) {
	var discounts []catalog.Discount
	discounts = append(discounts, catalog.NewCategoryDiscount("boots", 30), catalog.NewProductDiscount("000003", 15))

	productRepo := inmem.NewProductRepo(givenProducts)
	discountRepo := inmem.NewDiscountRepo(discounts)
	pricingCalculater := pricing.NewCalculater(discountRepo)
	productLister := listing.NewProductLister(productRepo, pricingCalculater)
	catalogService := rest.NewCatalogServer(productLister)

	invalidPrice := map[string]string{"priceLessThan": "Hello"}
	invalidLimit := map[string]string{
		"limit":  "Hi",
		"offset": "0",
	}
	invalidOffset := map[string]string{
		"limit":  "5",
		"offset": "null",
	}

	response := httptest.NewRecorder()

	tests := map[string]struct {
		search map[string]string
		status int
	}{
		"Invalid price filter": {
			search: invalidPrice,
			status: 400,
		},
		"Invalid pagination": {
			search: invalidLimit,
			status: 400,
		},
		"Invalid pagination offset": {
			search: invalidOffset,
			status: 400,
		},
	}

	for name, tc := range tests {
		catalogService.ServeHTTP(response, newListProductsRequest(t, tc.search))

		assert.Equal(t, tc.status, response.Code, name)
	}
}

func newPaginatedDiscountedProductsFromJSON(t *testing.T, rdr io.Reader) *catalog.PaginatedDiscountedProducts {
	t.Helper()
	var products *catalog.PaginatedDiscountedProducts
	err := json.NewDecoder(rdr).Decode(&products)

	if err != nil {
		t.Fatalf("fails json decoding to PaginatedDiscountedProducts %v", err)
	}

	return products
}

func newListProductsRequestFromSearchCriteria(t *testing.T, search catalog.SearchCriteria) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	if err != nil {
		t.Fatalf("fails creating newListProductsRequest %v", err)
	}
	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(search.Pagination().Limit))
	q.Add("offset", strconv.Itoa(search.Pagination().Offset))
	for _, f := range search.Filters() {
		switch filter := f.(type) {
		case catalog.CategoryFilter:
			q.Add("category", string(filter.Value()))
		case catalog.PriceLessThanFilter:
			q.Add("priceLessThan", strconv.Itoa(int(filter.Value())))
		}
	}
	req.URL.RawQuery = q.Encode()

	return req
}

func newListProductsRequest(t *testing.T, params map[string]string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	if err != nil {
		t.Fatalf("fails creating newListProductsRequest %v", err)
	}
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	return req
}
