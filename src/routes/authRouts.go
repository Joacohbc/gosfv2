package routes

import (
	"gosfV2/src/auth"

	"github.com/labstack/echo"
)

var Auth authRoutes

type authRoutes struct{}

// Agrego los Endpoints de Auth
func (a *authRoutes) AddRoutes(e *echo.Echo) {
	e.POST("/login", auth.LoginHandler)
	e.POST("/register", auth.RegisterHandler)
	e.GET("/logout", auth.LogoutHandler)
}
