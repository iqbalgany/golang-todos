package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

func NewAssignScopeController(e *echo.Echo, db *sql.DB)  {
	e.POST("/users/:userId/scopes/:scopeId/assign", func(ctx echo.Context) error {
		userId := ctx.Param("userId")
		scopeId := ctx.Param("scopeId")

		row := db.QueryRow("SELECT id FROM user_roles WHERE user_id = ? AND scope_id = ?", userId, scopeId)
		if row.Err() != nil {
			return ctx.String(http.StatusInternalServerError, row.Err().Error())
		}

		var retrivedId int 
		err := row.Scan(&retrivedId)
		if err == nil {
			return ctx.String(http.StatusBadRequest, "duplicate assignment found")
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		_, err = db.Exec(
			"INSERT INTO user_roles (user_id, scope_id) VALUES (?, ?)",
			userId,
			scopeId,
		)

		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		return ctx.String(http.StatusOK, "OK")
	})
}