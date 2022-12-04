package routes

import (
	"gosfV2/src/handlers"

	"github.com/labstack/echo"
)

var User userRoutes

type userRoutes struct{}

// Agrego los Endpoints de User
func (a *userRoutes) AddRoutesToGroup(group *echo.Group) {
	users := group.Group("/users")
	users.PUT("/rename", handlers.Users.RenameUser)
	users.PUT("/password", handlers.Users.ChangePassword)
}
