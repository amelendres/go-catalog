package listing

import (
	. "github.com/amelendres/go-catalog/catalog"
	"github.com/amelendres/go-catalog/pricing"
)

type ProductLister interface {
	List(search SearchCriteria) (*PaginatedDiscountedProducts, error)
}

type service struct {
	repository        ProductRepository
	pricingCalculater pricing.Calculater
}

func NewProductLister(r ProductRepository, pc pricing.Calculater) ProductLister {
	return &service{r, pc}
}

func (s service) List(search SearchCriteria) (*PaginatedDiscountedProducts, error) {

	paginatedProducts, err := s.repository.List(search)
	if err != nil {
		return nil, err
	}

	var discountedProducts []*DiscountedProduct
	for _, p := range paginatedProducts.Items() {

		price, err := s.pricingCalculater.Calculate(*p)
		if err != nil {
			return nil, err
		}
		dp := NewDiscountedProduct(p.SKU, p.Name, p.Category, *price)
		discountedProducts = append(discountedProducts, dp)
	}

	return NewPaginatedDiscountedProducts(paginatedProducts.Meta, discountedProducts), nil
}
