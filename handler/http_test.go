package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"

	"example/model"
	"example/service/mocks"
)

func TestHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockEmpService(ctrl)
	h := handler{service: m}
	_ = New(m)
	app := gofr.New()

	testcases := []struct {
		desc   string
		output interface{}
		err    error
		mock   []*gomock.Call
	}{
		{"success", []model.Employee{{ID: 1, Age: 21, Name: "Ram"}}, nil, []*gomock.Call{
			m.EXPECT().GetEmp(gomock.Any()).Return([]model.Employee{{ID: 1, Age: 21, Name: "Ram"}}, nil),
		}},
		{"Failure", nil, errors.Error("Connect Failed"), []*gomock.Call{
			m.EXPECT().GetEmp(gomock.Any()).Return(nil, errors.Error("Connect Failed")),
		}},
	}

	for _, tc := range testcases {
		r := httptest.NewRequest(http.MethodGet, "/emp", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		resp, err := h.Get(ctx)
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.output, resp, "Test Failed")
			assert.Equal(t, tc.err, err, "Test Failed")
		})
	}
}

func TestHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockEmpService(ctrl)
	h := handler{service: m}
	app := gofr.New()

	testcases := []struct {
		desc   string
		id     string
		output interface{}
		err    error
		mock   []*gomock.Call
	}{
		{"Failure", "31", model.Employee{}, errors.Error("Connect Failed"), []*gomock.Call{
			m.EXPECT().GetEmpByID(gomock.Any(), gomock.Any()).Return(model.Employee{}, errors.Error("Connect Failed")),
		}},
		{"Success", "2", model.Employee{ID: 2, Age: 22, Name: "Ram"}, nil, []*gomock.Call{
			m.EXPECT().GetEmpByID(gomock.Any(), gomock.Any()).Return(model.Employee{ID: 2, Age: 22, Name: "Ram"}, nil),
		}},
		{"ID_Empty", "", nil, errors.InvalidParam{Param: []string{"id"}}, nil},
		{"ID_Invalid", "jeh", nil, errors.Error("Failed to Convert id"), nil},
	}

	for i, tc := range testcases {
		r := httptest.NewRequest(http.MethodGet, "/emp/{id}", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		t.Run(tc.desc, func(t *testing.T) {
			ctx.SetPathParams(map[string]string{
				"id": tc.id,
			})
			resp, err := h.GetByID(ctx)
			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.output, resp)
			}
			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.err, err)
			}
		})
	}
}

func TestHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockEmpService(ctrl)
	h := handler{service: m}
	app := gofr.New()

	testcases := []struct {
		desc   string
		id     string
		req    []byte
		output interface{}
		err    error
		mock   []*gomock.Call
	}{
		{"ID EMPTY", "", []byte(`("id":1,"age":21,"name":"Ram")`), nil,
			errors.InvalidParam{Param: []string{"id"}}, nil},
		{"ID INVALID", "sd", []byte(`("id":2,"age":22,"name":"sai")`), nil,
			errors.Error("Failed to Convert to id"), nil},
		{"", "2", []byte(``), nil, errors.InvalidParam{Param: []string{"body"}}, nil},
		{desc: "Failure", id: "3", req: []byte(`{"id":3,"age":23,"name":"gopal"}`), output: nil,
			err: errors.Error("Connect Failed"), mock: []*gomock.Call{m.EXPECT().UpdateEmp(gomock.Any(), gomock.Any()).
			Return(model.Employee{},
				errors.Error("Connect Failed"))}},
		{desc: "Success", id: "4", req: []byte(`{"id":4,"age":24,"name":"harish"}`),
			output: model.Employee{ID: 4, Age: 24, Name: "harish"}, mock: []*gomock.Call{
			m.EXPECT().UpdateEmp(gomock.Any(), gomock.Any()).Return(model.Employee{ID: 4, Age: 24, Name: "harish"}, nil)}},
	}

	for i, tc := range testcases {
		r := httptest.NewRequest(http.MethodPut, "/emp/{id}", bytes.NewBuffer(tc.req))
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		cxt := gofr.NewContext(res, req, app)

		t.Run(tc.desc, func(t *testing.T) {
			cxt.SetPathParams(map[string]string{
				"id": tc.id,
			})
			resp, err := h.Update(cxt)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.err, err)
			}
		})
	}
}

func TestHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockEmpService(ctrl)
	h := handler{service: m}
	app := gofr.New()

	testcases := []struct {
		desc   string
		req    []byte
		output interface{}
		err    error
		mock   []*gomock.Call
	}{
		{"Unmarshal error", []byte(``), nil, errors.InvalidParam{Param: []string{"body"}}, nil},
		{"Failure", []byte(`{"id":1,"age":20,"name":"sai"}`), nil, errors.Error("Connect Failed"),
			[]*gomock.Call{m.EXPECT().CreateEmp(gomock.Any(), gomock.Any()).Return(model.Employee{},
				errors.Error("Connect Failed")),
			}},
		{desc: "Success", req: []byte(`{"id":2,"age":21,"name":"ram"}`), output: model.Employee{ID: 2, Age: 21, Name: "ram"},
			mock: []*gomock.Call{m.EXPECT().CreateEmp(gomock.Any(), gomock.Any()).
				Return(model.Employee{ID: 2, Age: 21, Name: "ram"}, nil)}},
	}

	for i, tc := range testcases {
		r := httptest.NewRequest(http.MethodPost, "/emp", bytes.NewReader(tc.req))
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		cxt := gofr.NewContext(res, req, app)

		t.Run(tc.desc, func(t *testing.T) {
			resp, err := h.Create(cxt)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed.Expected %v but Got %v", i+1, tc.err, err)
			}
		})
	}
}
