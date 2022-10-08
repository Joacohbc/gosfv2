package routes

import (
	"gosfV2/src/handlers"

	"github.com/labstack/echo"
)

var Files filesRoutes

type filesRoutes struct{}

// Agrego un nuevo Group a la API llamada Files y
// agrega los Endpoints de Files
func (f *filesRoutes) AddRoutes(e *echo.Echo) {
	files := e.Group("/files")
	files.POST("/", handlers.Files.UploadFile)
	files.GET("/", handlers.Files.GetFiles)
	files.DELETE("/:filename", handlers.Files.DeleteFile)
}

// Agrego los Endpoints de Files al grupo de Endpoints
func (f *filesRoutes) AddRoutesToGroup(group *echo.Group) {
	files := group.Group("/files")
	files.POST("/", handlers.Files.UploadFile)
	files.GET("/", handlers.Files.GetFiles)
	files.DELETE("/", handlers.Files.DeleteFile)
}
