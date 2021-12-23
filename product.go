package catalog

type SKU string
type Price int
type Category string

type ProductRepository interface {
	List(search SearchCriteria) (products *PaginatedProducts, err error)
}

type Product struct {
	SKU      SKU
	Name     string
	Category Category
	Price    Price
}

func NewProduct(SKU SKU, name string, category Category, price Price) *Product {
	return &Product{SKU, name, category, price}
}

type PaginatedProducts struct {
	Meta     PaginationMeta `json:"meta"`
	Products []*Product     `json:"items"`
}

func (p *PaginatedProducts) MetaData() PaginationMeta {
	return p.Meta
}

func (p *PaginatedProducts) Items() []*Product {
	return p.Products
}

func NewPaginatedProducts(meta PaginationMeta, p []*Product) *PaginatedProducts {
	return &PaginatedProducts{meta, p}
}
