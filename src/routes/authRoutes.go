package routes

import (
	"gosfV2/src/auth"

	"github.com/labstack/echo"
)

var Auth authRoutes

type authRoutes struct{}

// Agrego los Endpoints de Auth
func (a *authRoutes) AddAuthRoutes(group *echo.Group) {
	// Endpoints de Auth
	group.POST("/register", auth.RegisterUser)

	// Endpoints de Auth con Middleware de Credenciales
	group.POST("/login", auth.Login, auth.UserCredencialMiddleware)
	group.POST("/restore", auth.DeleteAllTokens, auth.UserCredencialMiddleware)

	// Endpoints de Auth con Middleware de JWT
	group.GET("/refresh", auth.RefreshToken, auth.JWTAuthMiddleware)
	group.GET("/verify", auth.VerifyAuth, auth.JWTAuthMiddleware)
	group.DELETE("/logout", auth.Logout, auth.JWTAuthMiddleware)
}
