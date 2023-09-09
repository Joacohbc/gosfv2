package routes

import (
	"gosfV2/src/handlers"

	"github.com/labstack/echo/v4"
)

var User userRoutes

type userRoutes struct{}

// Agrego los Endpoints de User
func (a *userRoutes) AddRoutesToGroup(group *echo.Group) {
	users := group.Group("/users")

	me := users.Group("/me")
	me.GET("", handlers.Users.GetUser)
	me.PUT("", handlers.Users.RenameUser)
	me.PUT("/password", handlers.Users.ChangePassword)
	me.DELETE("", handlers.Users.DeleteUser)
	me.GET("/icon", handlers.Users.GetIcon)
	me.POST("/icon", handlers.Users.UploadIcon)
	me.DELETE("/icon", handlers.Users.DeleteIcon)

	icon := users.Group("/icon")
	icon.GET("/:userId", handlers.Users.GetIconFromUser)
}
