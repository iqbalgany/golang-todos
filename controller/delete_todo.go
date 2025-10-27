package controller

import (
	"database/sql"
	"net/http"

	"github.com/iqbalgany/golang-todos/models"
	"github.com/labstack/echo"
)

func NewDeleteTodoController(e *echo.Echo, db *sql.DB)  {
	e.DELETE("/todos/:id", func(ctx echo.Context) error {
		user := ctx.Get("USER").(models.AuthClaimJwt)
		id := ctx.Param("id")

		permissionFound := false
		for _, scope := range user.UserScopes  {
			if scope == "todos:delete" {
				permissionFound = true
				break
			}
		}
		if !permissionFound {
			return ctx.String(http.StatusForbidden, "Forbidden")
		}


		_, err := db.Exec(
			"DELETE FROM todos WHERE id = ?",
			id,
			user.UserId,
		)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		return ctx.String(http.StatusOK, "OK")
	})
}