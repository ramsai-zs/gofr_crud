package employee

import (
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"example/model"
)

type store struct{}

func New() store {
	return store{}
}

func (s store) EmpGet(ctx *gofr.Context) ([]model.Employee, error) {
	var emp []model.Employee

	rows, err := ctx.DB().DB.Query("select * from employee")
	if err != nil {
		return nil, errors.DB{Err: errors.Error("Internal DB error")}
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var e model.Employee

		err := rows.Scan(&e.ID, &e.Age, &e.Name)

		if err != nil {
			return nil, errors.Error("Scan Error")
		}

		emp = append(emp, e)
	}

	return emp, nil
}

func (s store) EmpGetByID(ctx *gofr.Context, id int) (model.Employee, error) {
	var e model.Employee

	row := ctx.DB().DB.QueryRow("select * from employee where id = $1", id)

	err := row.Scan(&e.ID, &e.Age, &e.Name)

	if err != nil {
		return model.Employee{}, errors.Error("Scan Error")
	}

	return e, nil
}

func (s store) EmpCreate(ctx *gofr.Context, employee model.Employee) (model.Employee, error) {
	_, err := ctx.DB().DB.Exec("insert into employee(id,age,name) VALUES ($1,$2,$3)",
		employee.ID, employee.Age, employee.Name)

	if err != nil {
		return model.Employee{}, errors.Error("Internal DB Error")
	}

	return employee, nil
}

func (s store) EmpUpdate(ctx *gofr.Context, employee model.Employee) (model.Employee, error) {
	_, err := ctx.DB().DB.Exec("update employee set age = $1,name = $2 where id = $3",
		employee.Age, employee.Name, employee.ID)

	if err != nil {
		return model.Employee{}, errors.Error("Internal DB Error")
	}

	return employee, nil
}
