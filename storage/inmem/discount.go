package inmem

import (
	. "github.com/amelendres/go-catalog"
)

type DiscountRepo struct {
	products   map[string]Discount
	categories map[string]Discount
}

func (r *DiscountRepo) Find(search SearchCriteria) (discounts []Discount, err error) {
	var resp []Discount
	for _, f := range search.Filters() {
		switch filter := f.(type) {
		case CategoryFilter:
			if discount, ok := r.categories[string(filter.Value())]; ok {
				resp = append(resp, discount)
			}
		case SKUFilter:
			if discount, ok := r.products[string(filter.Value())]; ok {
				resp = append(resp, discount)
			}
		}
	}

	return resp, nil
}

func NewDiscountRepo(discounts []Discount) *DiscountRepo {
	productDiscounts := make(map[string]Discount)
	categoryDiscounts := make(map[string]Discount)

	for _, d := range discounts {
		switch discount := d.(type) {
		case *CategoryDiscount:
			categoryDiscounts[string(discount.Category())] = discount
		case *ProductDiscount:
			productDiscounts[string(discount.SKU())] = discount
		}
	}
	return &DiscountRepo{productDiscounts, categoryDiscounts}
}
