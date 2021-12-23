package catalog

const EURCurrency = Currency("EUR")

type DiscountPercentage int

type Currency string

type Discount interface {
	Percentage() DiscountPercentage
}

type DiscountRepository interface {
	Find(search SearchCriteria) (discounts []Discount, err error)
}

type ProductDiscount struct {
	sku        SKU
	percentage DiscountPercentage
}

func (d *ProductDiscount) Percentage() DiscountPercentage {
	return d.percentage
}

func (d *ProductDiscount) SKU() SKU {
	return d.sku
}

func NewProductDiscount(sku SKU, dp DiscountPercentage) Discount {
	return &ProductDiscount{sku, dp}
}

type CategoryDiscount struct {
	category   Category
	percentage DiscountPercentage
}

func (d *CategoryDiscount) Percentage() DiscountPercentage {
	return d.percentage
}

func (d *CategoryDiscount) Category() Category {
	return d.category
}

func NewCategoryDiscount(cat Category, dp DiscountPercentage) Discount {
	return &CategoryDiscount{cat, dp}
}

type DiscountedPrice struct {
	Original           Price               `json:"original"`
	Final              Price               `json:"final"`
	DiscountPercentage *DiscountPercentage `json:"discount_percentage"`
	Currenty           Currency            `json:"currency"`
}

func NewDiscountedPrice(original Price, dp *DiscountPercentage) *DiscountedPrice {
	final := original
	if dp != nil {
		final = Price(uint(original) - (uint(original) * uint(*dp) / 100))
	}
	return &DiscountedPrice{original, final, dp, EURCurrency}
}

type DiscountedProduct struct {
	SKU      SKU             `json:"sku"`
	Name     string          `json:"name"`
	Category Category        `json:"category"`
	Price    DiscountedPrice `json:"price"`
}

func NewDiscountedProduct(SKU SKU, name string, category Category, price DiscountedPrice) *DiscountedProduct {
	return &DiscountedProduct{SKU, name, category, price}
}

type PaginatedDiscountedProducts struct {
	Meta     PaginationMeta       `json:"meta"`
	Products []*DiscountedProduct `json:"items"`
}

func NewPaginatedDiscountedProducts(meta PaginationMeta, products []*DiscountedProduct) *PaginatedDiscountedProducts {
	return &PaginatedDiscountedProducts{meta, products}
}

func (p *PaginatedDiscountedProducts) MetaData() PaginationMeta {
	return p.Meta
}

func (p *PaginatedDiscountedProducts) Items() []*DiscountedProduct {
	return p.Products
}
