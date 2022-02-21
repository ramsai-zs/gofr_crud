package employees

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"example/datastore"
	"example/model"
)

type service struct {
	store datastore.EmpStore
}

func New(s datastore.EmpStore) service {
	return service{store: s}
}

func (s service) GetEmp(ctx *gofr.Context) ([]model.Employee, error) {
	resp, err := s.store.EmpGet(ctx)

	if err != nil {
		return nil, errors.Error("Connect Failed")
	}
	return resp, nil
}

func (s service) GetEmpByID(ctx *gofr.Context, id int) (model.Employee, error) {
	resp, err := s.store.EmpGetByID(ctx, id)

	if err != nil {
		return model.Employee{}, errors.Error("Connect Failed")
	}

	return resp, nil
}

func (s service) CreateEmp(ctx *gofr.Context, employee model.Employee) (model.Employee, error) {
	resp, err := s.store.EmpCreate(ctx, employee)

	if err != nil {
		return model.Employee{}, errors.Error("Connect Failed")
	}

	return resp, err
}

func (s service) UpdateEmp(ctx *gofr.Context, employee model.Employee) (model.Employee, error) {
	resp, err := s.store.EmpUpdate(ctx, employee)

	if err != nil {
		return model.Employee{}, errors.Error("Connect Failed")
	}

	return resp, err
}
