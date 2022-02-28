package employee

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"example/model"
)

func TestStore_EmpGet(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Errorf("database error %s", err)
	}

	defer db.Close()

	g := gofr.Gofr{DataStore: datastore.DataStore{ORM: db}}
	ctx := gofr.NewContext(nil, nil, &g)
	ctx.Context = context.Background()
	dataStore := New()

	query := "select * from employee"

	row := sqlmock.NewRows([]string{"id", "age", "name"}).AddRow(1, 21, "Ram")
	scanError := sqlmock.NewRows([]string{"id", "age", "name", "err"}).AddRow(1, 21, "Ram", "error")

	testcases := []struct {
		desc   string
		output []model.Employee
		err    error
		Mock   []interface{}
	}{
		{desc: "Failure", err: errors.DB{Err: errors.Error("Internal DB error")}, Mock: []interface{}{
			mock.ExpectQuery(query).WithArgs().WillReturnError(errors.DB{Err: errors.Error("Internal DB error")}),
		}},
		{"success", []model.Employee{{1, 21, "Ram"}}, nil, []interface{}{
			mock.ExpectQuery(query).WithArgs().WillReturnRows(row),
		}},
		{"ScanError", nil, errors.Error("Scan Error"), []interface{}{
			mock.ExpectQuery(query).WithArgs().WillReturnRows(scanError),
		}},
	}
	for i, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			tc := tc
			resp, err := dataStore.EmpGet(ctx)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed. Expected %v but got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed. Expected %v but got %v", i+1, tc.output, resp)
			}
		})
	}
}

func TestStore_EmpGetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Errorf("Database error %v", err)
	}

	defer db.Close()

	g := gofr.Gofr{DataStore: datastore.DataStore{ORM: db}}
	cxt := gofr.NewContext(nil, nil, &g)
	cxt.Context = context.Background()
	dataStore := New()

	query := "select * from employee where id = $1"

	row := sqlmock.NewRows([]string{"id", "age", "name"}).AddRow(1, 21, "Ram")
	scanError := sqlmock.NewRows([]string{"id", "age", "name", "err"}).AddRow(1, 21, "Ram", "error")

	testcases := []struct {
		desc   string
		id     int
		output model.Employee
		err    error
		mock   []interface{}
	}{
		{desc: "success", id: 1, output: model.Employee{ID: 1, Age: 21, Name: "Ram"}, mock: []interface{}{
			mock.ExpectQuery(query).WithArgs(1).WillReturnRows(row),
		}},
		{"scanError", 3, model.Employee{}, errors.Error("Scan Error"), []interface{}{
			mock.ExpectQuery(query).WithArgs(3).WillReturnRows(scanError),
		}},
	}

	for i, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			resp, err := dataStore.EmpGetByID(cxt, tc.id)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed. Expected %v but got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed. Expected %v but got %v", i+1, tc.output, resp)
			}
		})
	}
}

func TestStore_EmpCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Errorf("Database error")
	}

	defer db.Close()

	g := gofr.Gofr{DataStore: datastore.DataStore{ORM: db}}
	cxt := gofr.NewContext(nil, nil, &g)
	cxt.Context = context.Background()
	dataStore := New()

	exec := "insert into employee(id,age,name) VALUES ($1,$2,$3)"
	execError := "insert into employee(id,age,name,err) VALUES ($1,$2,$3,$4)"

	testcases := []struct {
		desc   string
		input  model.Employee
		output model.Employee
		err    error
		mock   []interface{}
	}{
		{desc: "Success", input: model.Employee{ID: 1, Age: 21, Name: "Ram"}, output: model.Employee{ID: 1, Age: 21, Name: "Ram"},
			mock: []interface{}{mock.ExpectExec(exec).WithArgs(1, 21, "Ram").WillReturnResult(sqlmock.NewResult(1, 1))}},
		{desc: "Failure", input: model.Employee{ID: 2, Age: 22, Name: "Sai"}, err: errors.Error("Internal DB Error"),
			mock: []interface{}{mock.ExpectExec(execError).WithArgs(2, 22, "Sai", err).
				WillReturnError(errors.Error("Internal DB Error")),
			}},
	}

	for i, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			resp, err := dataStore.EmpCreate(cxt, tc.input)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed. Expected %v but got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed. Expected %v but got %v", i+1, tc.err, err)
			}
		})
	}
}

func TestStore_EmpUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Errorf("Database error")
	}

	defer db.Close()

	g := gofr.Gofr{DataStore: datastore.DataStore{ORM: db}}
	cxt := gofr.NewContext(nil, nil, &g)
	cxt.Context = context.Background()
	dataStore := New()

	update := "update employee set age = $1,name = $2 where id = $3"
	updateError := "update set age = $1,name = $2,err = $3 where id = $4"

	testcases := []struct {
		desc   string
		input  model.Employee
		output model.Employee
		err    error
		mock   []interface{}
	}{
		{desc: "success", input: model.Employee{ID: 1, Age: 21, Name: "Ram"}, output: model.Employee{ID: 1, Age: 21, Name: "Ram"},
			mock: []interface{}{mock.ExpectExec(update).WithArgs(21, "Ram", 1).WillReturnResult(sqlmock.NewResult(1, 1))}},
		{desc: "Failure", input: model.Employee{ID: 2, Age: 22, Name: "Sai"}, err: errors.Error("Internal DB Error"),
			mock: []interface{}{mock.ExpectExec(updateError).WithArgs(2, 22, "Sai", 2, err).WillReturnError(errors.Error("Internal DB Error"))},
		},
	}

	for i, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			resp, err := dataStore.EmpUpdate(cxt, tc.input)

			if !reflect.DeepEqual(tc.output, resp) {
				t.Errorf("[Test %v]Failed. Expected %v but got %v", i+1, tc.output, resp)
			}

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("[Test %v]Failed. Expected %v but got %v", i+1, tc.err, err)
			}
		})
	}
}
