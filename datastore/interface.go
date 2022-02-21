package datastore

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"example/model"
)

type EmpStore interface {
	EmpGet(ctx *gofr.Context) ([]model.Employee, error)
	EmpGetByID(ctx *gofr.Context, id int) (model.Employee, error)
	EmpCreate(ctx *gofr.Context, employee model.Employee) (model.Employee, error)
	EmpUpdate(ctx *gofr.Context, employee model.Employee) (model.Employee, error)
}
