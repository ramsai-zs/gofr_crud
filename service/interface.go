package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"example/model"
)

type EmpService interface {
	GetEmp(ctx *gofr.Context) ([]model.Employee, error)
	GetEmpByID(ctx *gofr.Context, id int) (model.Employee, error)
	CreateEmp(ctx *gofr.Context, employee model.Employee) (model.Employee, error)
	UpdateEmp(ctx *gofr.Context, employee model.Employee) (model.Employee, error)
}
