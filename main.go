package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"example/datastore/employee"
	"example/handler"
	"example/middleware"
	"example/service/employees"
)

func main() {
	app := gofr.New()
	app.Server.UseMiddleware(middleware.Oauth)

	store := employee.New()
	service := employees.New(store)
	h := handler.New(service)

	app.GET("/emp", h.Get)
	app.GET("/emp/{id}", h.GetByID)
	app.PUT("/emp/{id}", h.Update)
	app.POST("/emp", h.Create)

	app.Server.HTTP.Port = 9090
	app.EnableSwaggerUI()
	app.Start()
}
