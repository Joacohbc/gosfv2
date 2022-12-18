package routes

import (
	"gosfV2/src/auth"

	"github.com/labstack/echo"
)

var Auth authRoutes

type authRoutes struct{}

// Agrego los Endpoints de Auth
func (a *authRoutes) AddAuthRoutes(group *echo.Group) {
	group.POST("/register", auth.RegisterUser)
	group.POST("/login", auth.Login, auth.UserCredencialMiddleware)
	group.POST("/restore", auth.DeleteAllTokens, auth.UserCredencialMiddleware)

	group.GET("/refresh", auth.RefreshToken, auth.JWTAuthMiddleware)
	group.GET("/verify", auth.VerifyAuth, auth.JWTAuthMiddleware)
	group.DELETE("/logout", auth.Logout, auth.JWTAuthMiddleware)
}
