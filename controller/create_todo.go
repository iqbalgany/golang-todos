package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/iqbalgany/golang-todos/models"
	"github.com/labstack/echo"
)

type CreateRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
}

func NewCreateTodoController(e *echo.Echo, db *sql.DB)  {
	e.POST("/todos", func(ctx echo.Context) error {
		user := ctx.Get("USER").(models.AuthClaimJwt)

		permissionFound := false
		for _, scope := range user.UserScopes  {
			if scope == "todos:create" {
				permissionFound = true
				break
			}
		}
		if !permissionFound {
			return ctx.String(http.StatusForbidden, "Forbidden")
		}


		var request CreateRequest
		json.NewDecoder(ctx.Request().Body).Decode(&request)

		_, err := db.Exec(
			"INSERT INTO todos (title, description, done, user_id) VALUES (?, ?, 0, ?)",
			request.Title,
			request.Description,
			user.UserId,
		)

		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		return ctx.String(http.StatusOK, "OK")
	})
}