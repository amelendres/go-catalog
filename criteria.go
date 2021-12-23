package catalog

type Filter interface {
	// Value() string
	// Value() interface
}

// type Filter struct {
// 	name  string
// 	value string
// }

type CategoryFilter struct {
	category Category
}

func (f *CategoryFilter) Value() Category {
	return f.category
}

func NewCategoryFilter(cat Category) Filter {
	return &CategoryFilter{cat}
}

type PriceLessThanFilter struct {
	price Price
}

func (f *PriceLessThanFilter) Value() Price {
	return f.price
}

func NewPriceLessThanFilter(p Price) Filter {
	return &PriceLessThanFilter{p}
}

type SearchCriteria struct {
	pagination *Pagination
	filters    []Filter
}

func (s *SearchCriteria) Pagination() *Pagination {
	return s.pagination
}

func (s *SearchCriteria) Filters() []Filter {
	return s.filters
}

func NewSearchCriteria(pag *Pagination, filters []Filter) (*SearchCriteria, error) {
	return &SearchCriteria{pag, filters}, nil
}
