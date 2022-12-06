package routes

import (
	"gosfV2/src/auth"

	"github.com/labstack/echo"
)

var Auth authRoutes

type authRoutes struct{}

// Agrego los Endpoints de Auth
func (a *authRoutes) AddNoAuthRoutes(e *echo.Echo) {
	e.POST("/login", auth.LoginHandler)
	e.POST("/register", auth.RegisterHandler)
	e.DELETE("/tokens", auth.DeleteTokens)
}

func (a *authRoutes) AddTokenRoutes(group *echo.Group) {
	group.GET("/auth", auth.VerifyAuth)
	group.GET("/refresh", auth.RefreshHandler)
	group.POST("/logout", auth.LogoutHandler)
}
