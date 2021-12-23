package inmem

import (
	. "github.com/amelendres/go-catalog"
)

type ProductRepo struct {
	products []*Product
}

func (r *ProductRepo) List(search SearchCriteria) (products *PaginatedProducts, err error) {

	filteredProducts := r.filter(search.Filters())
	paginated := paginate(filteredProducts, *search.Pagination())

	return paginated, nil
}

func (r *ProductRepo) filter(filters []Filter) []*Product {

	var filteredProducts []*Product
	if len(filters) == 0 {
		filteredProducts = r.products
	}

	for _, f := range filters {
		switch filter := f.(type) {
		case CategoryFilter:
			if filteredProducts != nil {
				filteredProducts = append(filteredProducts, filterByCategory(filteredProducts, filter.Value())...)
			}
			filteredProducts = append(filteredProducts, filterByCategory(r.products, filter.Value())...)
		case PriceLessThanFilter:
			if filteredProducts != nil {
				filteredProducts = append(filteredProducts, filterByPriceLessThan(filteredProducts, filter.Value())...)
			}
			filteredProducts = append(filteredProducts, filterByPriceLessThan(r.products, filter.Value())...)

		}
	}
	return filteredProducts
}

func paginate(products []*Product, pag Pagination) *PaginatedProducts {

	if len(products) == 0 || pag.Offset > len(products) {
		return NewPaginatedProducts(
			PaginationMeta{Total: len(products), Pagination: pag},
			nil,
		)
	}
	to := pag.Limit
	if len(products) < pag.Offset+to {
		to = len(products)
	}

	return NewPaginatedProducts(
		PaginationMeta{Total: len(products), Pagination: pag},
		products[pag.Offset:to],
	)
}

func NewProductRepo(p []*Product) *ProductRepo {
	return &ProductRepo{p}
}

func filterByCategory(products []*Product, cat Category) []*Product {
	var filtered []*Product
	for _, p := range products {
		if p.Category == cat {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func filterByPriceLessThan(products []*Product, price Price) []*Product {
	var filtered []*Product
	for _, p := range products {
		if p.Price <= price {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
