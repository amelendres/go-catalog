package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/amelendres/go-catalog"
	"github.com/amelendres/go-catalog/http/rest"
	"github.com/amelendres/go-catalog/listing"
	"github.com/amelendres/go-catalog/pricing"
	"github.com/amelendres/go-catalog/storage/inmem"
)

const (
	givenProductsJSON = `
	[
    	{		
			"sku": "000001",
			"name": "BV Lean leather ankle boots",
			"category": "boots",
			"price": 89000
		}, 
		{		
			"sku": "000002",
			"name": "BV Lean leather ankle boots",
			"category": "boots",
			"price": 99000
		}, 
		{	
			"sku": "000003",
			"name": "Ashlington leather ankle boots",
			"category": "boots",
			"price": 71000
		},
		{
			"sku": "000004",
			"name": "Naima embellished suede sandals",
			"category": "sandals",
			"price": 79500
		}, 
		{
			"sku": "000005",
			"name": "Nathane leather sneakers",
			"category": "sneakers",
			"price": 59000
		},
		{
			"sku": "000006",
			"name": "AA hat",
			"category": "hats",
			"price": 72000
		}
	]	
	`
)

var (
	givenCategoryDiscount = catalog.DiscountPercentage(30)
	givenProductDiscount  = catalog.DiscountPercentage(15)
	discounts             = []catalog.Discount{
		catalog.NewCategoryDiscount("boots", givenCategoryDiscount),
		catalog.NewProductDiscount("000003", givenProductDiscount),
	}
)

func main() {
	productRepo := inmem.NewProductRepo(newProductsFromJSON(givenProductsJSON))
	discountRepo := inmem.NewDiscountRepo(discounts)
	pricingCalculater := pricing.NewCalculater(discountRepo)
	productLister := listing.NewProductLister(productRepo, pricingCalculater)

	cs := rest.NewCatalogServer(productLister)

	if err := http.ListenAndServe(":5000", cs); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

func newProductsFromJSON(jsonStr string) []*catalog.Product {
	var products []*catalog.Product
	err := json.Unmarshal([]byte(jsonStr), &products)
	if err != nil {
		log.Fatalf("could not load products %v", err)
	}

	return products
}
