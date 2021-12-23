package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/amelendres/go-catalog"
	"github.com/amelendres/go-catalog/listing"
	"github.com/gorilla/mux"
)

type CatalogServer struct {
	productLister listing.ProductLister
	http.Handler
}

const (
	jsonContentType = "application/json"

	defaultLimit  = 5
	defaultOffset = 0
)

func NewCatalogServer(pl listing.ProductLister) *CatalogServer {
	cs := new(CatalogServer)
	cs.productLister = pl

	router := mux.NewRouter()
	router.HandleFunc("/products", cs.listProducts).Methods(http.MethodGet)

	cs.Handler = router

	return cs
}

func (cs *CatalogServer) listProducts(w http.ResponseWriter, r *http.Request) {
	searchCriteria, err := buildSearchCriteria(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lp, err := cs.productLister.List(*searchCriteria)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	_ = json.NewEncoder(w).Encode(lp)
}

func buildSearchCriteria(r *http.Request) (search *catalog.SearchCriteria, err error) {
	//pagination
	iLimit := defaultLimit
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		if iLimit, err = strconv.Atoi(limit); err != nil {
			return nil, err
		}
	}
	iOffset := defaultOffset
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		if iOffset, err = strconv.Atoi(offset); err != nil {
			return nil, err
		}
	}

	pag, err := catalog.NewPagination(iLimit, iOffset)
	if err != nil {
		return nil, err
	}

	//Filters
	var filters []catalog.Filter
	category := r.URL.Query().Get("category")
	if category != "" {
		filters = append(filters, catalog.NewCategoryFilter(catalog.Category(category)))
	}
	priceLessThan := r.URL.Query().Get("priceLessThan")
	if priceLessThan != "" {
		i, err := strconv.Atoi(priceLessThan)
		if err != nil {
			return nil, err
		}
		filters = append(filters, catalog.NewPriceLessThanFilter(catalog.Price(i)))
	}
	criteria := catalog.NewSearchCriteria(pag, filters)
	return &criteria, nil
}
