package catalog

// type PaginatedCollection interface {
// 	MetaData() PaginationMeta
// 	// Items() []*interface{}
// 	Items() []*interface{}
// }

// type PaginatedCollection struct {
// 	Meta  PaginationMeta `json:"meta"`
// 	Items []*interface{} `json:"items"`
// }

type PaginationMeta struct {
	Total      int `json:"meta"`
	Pagination `json:"pagination"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func NewPagination(limit, offset int) (*Pagination, error) {
	if err := assertLimit(limit); err != nil {
		return nil, err
	}
	if err := assertOffset(limit); err != nil {
		return nil, err
	}

	return &Pagination{limit, offset}, nil
}

func assertLimit(limit int) error {
	//WIP
	return nil
}

func assertOffset(offset int) error {
	//WIP
	return nil
}
