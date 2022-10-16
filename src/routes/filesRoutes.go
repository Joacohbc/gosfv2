package routes

import (
	"gosfV2/src/handlers"

	"github.com/labstack/echo"
)

var Files filesRoutes

type filesRoutes struct{}

// Agrego los Endpoints de Files al grupo de Endpoints
func (f *filesRoutes) AddRoutesToGroup(group *echo.Group) {
	files := group.Group("/files")

	// Consultas
	files.GET("/", handlers.Files.GetAllFiles)
	files.GET("/:id", handlers.Files.GetFile)

	// Creación
	files.POST("/", handlers.Files.UploadFile)

	// Borrar
	files.DELETE("/:id", handlers.Files.DeleteFile)

	// Modicacion
	files.PUT("/:id", handlers.Files.UpdateFile)

	// Opciones de Share
	files.GET("/share/:id", handlers.Files.GetSharedFile)
}
