package pricing

import (
	. "github.com/amelendres/go-catalog/catalog"
)

type Calculater interface {
	Calculate(p Product) (*DiscountedPrice, error)
}

type service struct {
	repository DiscountRepository
}

func NewCalculater(r DiscountRepository) Calculater {
	return service{repository: r}
}

func (s service) Calculate(p Product) (*DiscountedPrice, error) {

	criteria := NewSearchCriteria(nil, []Filter{
		NewCategoryFilter(p.Category),
		NewSKUFilter(p.SKU),
	})
	discounts, err := s.repository.Find(criteria)
	if err != nil {
		return nil, err
	}

	if discounts == nil {
		return NewDiscountedPrice(p.Price, nil), nil
	}

	var discount = discounts[0]
	for _, d := range discounts {
		if d.Percentage() > discount.Percentage() {
			discount = d
		}
	}
	discountPercentage := discount.Percentage()
	return NewDiscountedPrice(p.Price, &discountPercentage), nil
}
