package main

import (
	"github.com/iqbalgany/golang-todos/controller"
	"github.com/iqbalgany/golang-todos/database"
	"github.com/labstack/echo"
)



func main()  {
	db := database.InitDB()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	controller.NewGetAllTodosController(e, db)
	controller.NewCreateTodoController(e, db)
	controller.NewDeleteTodoController(e, db)
	controller.NewUpdateTodoController(e, db)
	controller.NewCheckTodoController(e, db)

	e.Start(":8080")
}