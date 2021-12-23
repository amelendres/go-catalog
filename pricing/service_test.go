package pricing_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/amelendres/go-catalog"
	"github.com/amelendres/go-catalog/pricing"
	"github.com/amelendres/go-catalog/testing/stub"
	"github.com/stretchr/testify/assert"
)

var (
	givenCategoryDiscount = catalog.DiscountPercentage(30)
	givenProductDiscount  = catalog.DiscountPercentage(15)
)

func TestPricingCalculater_Calculate(t *testing.T) {
	var discounts []catalog.Discount
	discounts = append(discounts, catalog.NewProductDiscount("000003", givenProductDiscount), catalog.NewCategoryDiscount("boots", givenCategoryDiscount))
	bootsProduct := *catalog.NewProduct("000003", "Ashlington leather ankle boots", "boots", 71000)
	sandalsProduct := *catalog.NewProduct("000004", "Naima embellished suede sandals", "sandals", 79500)

	discountRepo := stub.NewStubDiscountRepo(discounts, nil)
	pricingCalculater := pricing.NewCalculater(discountRepo)

	pricingCalculaterWithoutDiscounts := pricing.NewCalculater(stub.NewStubDiscountRepo(nil, nil))

	discountRepoErr := errors.New("fails discount repository")
	discountRepoWithErr := stub.NewStubDiscountRepo(discounts, discountRepoErr)
	pricingCalculaterWithErr := pricing.NewCalculater(discountRepoWithErr)

	tests := map[string]struct {
		in      pricing.Calculater
		to      catalog.Product
		want    *catalog.DiscountedPrice
		wantErr error
	}{
		"With discount": {
			in:      pricingCalculater,
			to:      bootsProduct,
			want:    catalog.NewDiscountedPrice(bootsProduct.Price, &givenCategoryDiscount),
			wantErr: nil,
		},
		"Without Discount": {
			in:      pricingCalculaterWithoutDiscounts,
			to:      sandalsProduct,
			want:    catalog.NewDiscountedPrice(sandalsProduct.Price, nil),
			wantErr: nil,
		},
		"Discount repository error": {
			in:      pricingCalculaterWithErr,
			to:      bootsProduct,
			want:    nil,
			wantErr: discountRepoErr,
		},
	}

	for name, tc := range tests {
		got, err := tc.in.Calculate(tc.to)

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
