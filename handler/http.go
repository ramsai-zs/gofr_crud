package handler

import (
	"strconv"

	"example/model"
	"example/service"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type handler struct {
	service service.EmpService
}

// nolint:revive // handlers should not be used without proper initialization with required dependency
func New(h service.EmpService) handler {
	return handler{service: h}
}

func (h handler) Get(c *gofr.Context) (interface{}, error) {
	resp, err := h.service.GetEmp(c)

	if err != nil {
		return nil, errors.Error("Connect Failed")
	}

	return resp, nil
}

func (h handler) GetByID(c *gofr.Context) (interface{}, error) {
	i := c.PathParam("id")

	if i == "" {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)

	if err != nil {
		return nil, errors.Error("Failed to Convert id")
	}

	resp, err := h.service.GetEmpByID(c, id)

	if err != nil {
		return model.Employee{}, errors.Error("Connect Failed")
	}

	return resp, nil
}

func (h handler) Update(c *gofr.Context) (interface{}, error) {
	var emp model.Employee

	i := c.PathParam("id")

	if i == "" {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)

	if err != nil {
		return nil, errors.Error("Failed to Convert to id")
	}

	if err = c.Bind(&emp); err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	emp.ID = id

	resp, err := h.service.UpdateEmp(c, emp)
	if err != nil {
		return nil, errors.Error("Connect Failed")
	}

	return resp, nil
}

func (h handler) Create(c *gofr.Context) (interface{}, error) {
	var emp model.Employee

	if err := c.Bind(&emp); err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := h.service.CreateEmp(c, emp)

	if err != nil {
		return nil, errors.Error("Connect Failed")
	}

	return resp, nil
}
