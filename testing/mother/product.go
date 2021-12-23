package mother

import "github.com/amelendres/go-catalog"

func NewPaginatedDiscountedProducts(
	products []*catalog.Product,
	discountedPrices map[string]*catalog.DiscountedPrice,
	pag catalog.Pagination,
	total int,
) *catalog.PaginatedDiscountedProducts {
	var discountedProducts []*catalog.DiscountedProduct
	for _, p := range products {
		discountedProducts = append(discountedProducts, catalog.NewDiscountedProduct(
			p.SKU,
			p.Name,
			p.Category,
			*discountedPrices[string(p.SKU)],
		))
	}

	return catalog.NewPaginatedDiscountedProducts(
		catalog.PaginationMeta{
			Total:      total,
			Pagination: pag,
		},
		discountedProducts,
	)
}
