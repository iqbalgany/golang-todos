package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iqbalgany/golang-todos/models"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func NewLoginController(e *echo.Echo, db *sql.DB)  {
	e.POST("/auth/login", func(ctx echo.Context) error {
		var request LoginRequest
		json.NewDecoder(ctx.Request().Body).Decode(&request)

		row := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", request.Email)
		if row.Err() != nil {
			return ctx.String(http.StatusInternalServerError, row.Err().Error())
		}

		var retrivedId int
		var retrivedName, retrivedEmail, retrivedPassword string

		err := row.Scan(&retrivedId, &retrivedName, &retrivedEmail, &retrivedPassword)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ctx.String(http.StatusUnauthorized, "email is not registered")
			}
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		rows, err := db.Query(
			"SELECT scopes.name AS scope_name FROM users LEFT JOIN user_roles ON user_roles.user_id = users.id JOIN scopes ON scopes.id = user_roles.scope_id WHERE email = ?",
			retrivedEmail,
		)

		var scopes []string = make([]string, 0)
		for rows.Next() {
			var scope string

			rows.Scan(&scope)
			if err != nil {
				 ctx.String(http.StatusInternalServerError, err.Error())
			}

			scopes = append(scopes, scope)
		}


		
		err = bcrypt.CompareHashAndPassword([]byte(retrivedPassword), []byte(request.Password))
		if err != nil {
			return ctx.String(http.StatusUnauthorized, err.Error())
		}

		tokenClaim := models.AuthClaimJwt{
			UserId: retrivedId,
			UserName: retrivedName,
			UserEmail: retrivedEmail,
			UserScopes: scopes,
		}
		
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)
		tokenStr, err := token.SignedString([]byte("TEST"))
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		response := LoginResponse{
			AccessToken: tokenStr,
		}

		return ctx.JSON(http.StatusOK, response)
	})
}