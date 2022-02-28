package employees

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"example/datastore/mocks"
	"example/model"
)

func TestService_GetEmp(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockEmpStore(ctrl)
	s := service{store: m}
	app := gofr.New()
	_ = New(m)

	testcases := []struct {
		desc   string
		output []model.Employee
		err    error
		mock   []*gomock.Call
	}{
		{"success", []model.Employee{{ID: 2, Age: 21, Name: "Ram"}}, nil, []*gomock.Call{
			m.EXPECT().EmpGet(gomock.Any()).Return([]model.Employee{{ID: 2, Age: 21, Name: "Ram"}}, nil),
		}},
		{desc: "failure", output: nil, err: errors.Error("Connect Failed"), mock: []*gomock.Call{
			m.EXPECT().EmpGet(gomock.Any()).Return(nil, errors.Error("Connect Failed"))}},
	}

	for i, tc := range testcases {
		tc := tc
		ctx := gofr.NewContext(nil, nil, app)
		ctx.Context = context.Background()

		t.Run(tc.desc, func(t *testing.T) {
			resp, err := s.GetEmp(ctx)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.err, err)
			}
		})
	}
}

func TestService_GetEmpByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockEmpStore(ctrl)
	s := service{store: m}
	app := gofr.New()

	testcases := []struct {
		desc   string
		id     int
		output model.Employee
		err    error
		mock   []*gomock.Call
	}{
		{desc: "Success", id: 1, output: model.Employee{ID: 1, Age: 21, Name: "Ram"}, mock: []*gomock.Call{m.EXPECT().
			EmpGetByID(gomock.Any(), gomock.Any()).Return(model.Employee{ID: 1, Age: 21, Name: "Ram"}, nil),
		}},
		{"Failure", 2, model.Employee{}, errors.Error("Connect Failed"), []*gomock.Call{m.EXPECT().
			EmpGetByID(gomock.Any(), gomock.Any()).Return(model.Employee{}, errors.Error("Connect Failed")),
		}},
	}

	for i, tc := range testcases {
		tc := tc
		cxt := gofr.NewContext(nil, nil, app)

		t.Run(tc.desc, func(t *testing.T) {
			resp, err := s.GetEmpByID(cxt, tc.id)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.err, err)
			}
		})
	}
}

func TestService_CreateEmp(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockEmpStore(ctrl)
	s := service{store: m}
	app := gofr.New()

	testcases := []struct {
		desc   string
		input  model.Employee
		output model.Employee
		err    error
		mock   []*gomock.Call
	}{
		{desc: "Success", input: model.Employee{ID: 1, Age: 20, Name: "Ram"}, output: model.Employee{ID: 1, Age: 20, Name: "Ram"},
			mock: []*gomock.Call{m.EXPECT().EmpCreate(gomock.Any(), gomock.Any()).Return(model.Employee{ID: 1, Age: 20, Name: "Ram"}, nil)}},
		{"Failure", model.Employee{2, 20, "Sai"}, model.Employee{},
			errors.Error("Connect Failed"), []*gomock.Call{m.EXPECT().EmpCreate(gomock.Any(), gomock.Any()).
				Return(model.Employee{}, errors.Error("Connect Failed"))}},
	}

	for i, tc := range testcases {
		tc := tc
		cxt := gofr.NewContext(nil, nil, app)

		t.Run(tc.desc, func(t *testing.T) {
			resp, err := s.CreateEmp(cxt, tc.input)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %d]Failed.Expected %v but Got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %d]Failed.Expected %v but Got %v", i+1, tc.err, err)
			}
		})
	}
}

func TestService_UpdateEmp(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockEmpStore(ctrl)
	s := service{store: m}
	app := gofr.New()

	testcases := []struct {
		desc   string
		id     int
		input  model.Employee
		output model.Employee
		err    error
		mock   []*gomock.Call
	}{
		{"success", 1, model.Employee{1, 21, "Ram"},
			model.Employee{1, 21, "Ram"}, nil, []*gomock.Call{m.EXPECT().
				EmpUpdate(gomock.Any(), gomock.Any()).Return(model.Employee{1, 21, "Ram"}, nil),
			}},
		{"Failure", 2, model.Employee{2, 22, "sai"}, model.Employee{},
			errors.Error("Connect Failed"), []*gomock.Call{m.EXPECT().EmpUpdate(gomock.Any(), gomock.Any()).
				Return(model.Employee{}, errors.Error("Connect Failed"))}},
	}

	for i, tc := range testcases {
		tc := tc
		cxt := gofr.NewContext(nil, nil, app)

		t.Run(tc.desc, func(t *testing.T) {
			resp, err := s.UpdateEmp(cxt, tc.input)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.output, resp)
			}
		})
	}
}
