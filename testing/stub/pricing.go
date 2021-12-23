package stub

import . "github.com/amelendres/go-catalog"

type StubDiscountRepo struct {
	discounts []Discount
	wantErr   error
}

func (r *StubDiscountRepo) Find(search SearchCriteria) (discounts []Discount, err error) {
	if r.wantErr != nil {
		return nil, r.wantErr
	}
	return r.discounts, nil
}

func NewStubDiscountRepo(d []Discount, wantErr error) *StubDiscountRepo {
	return &StubDiscountRepo{d, wantErr}
}
