package routes

import (
	"gosfV2/src/auth"
	"net/http"

	"github.com/labstack/echo"
)

var Auth authRoutes

type authRoutes struct{}

// Agrego los Endpoints de Auth
func (a *authRoutes) AddRoutes(e *echo.Echo) {
	e.POST("/login", auth.LoginHandler)
	e.POST("/register", auth.RegisterHandler)
	e.POST("/logout", auth.LogoutHandler)
	e.GET("/auth", func(c echo.Context) error {
		return c.String(http.StatusOK, "You are authenticated")
	})
}
