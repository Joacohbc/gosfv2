package routes

import (
	"gosfV2/src/handlers"

	"github.com/labstack/echo/v4"
)

var Files filesRoutes

type filesRoutes struct{}

// Agrego los Endpoints de Files
func (f *filesRoutes) AddFilesRoutes(group *echo.Group) {
	files := group.Group("/files")

	// Consultas
	files.GET("", handlers.Files.GetAllFiles)
	files.GET("/:id", handlers.Files.GetFile)
	files.GET("/:id/info", handlers.Files.GetInfo)

	// Creación
	files.POST("", handlers.Files.UploadFile)

	// Borrar
	files.DELETE("/:id", handlers.Files.DeleteFile)

	// Modificación
	files.PUT("/:id", handlers.Files.UpdateFile)

	// Opciones de Share
	files.GET("/share", handlers.Files.GetAllShareFiles)
	files.GET("/share/:id", handlers.Files.GetSharedFile)
	files.GET("/share/:id/info", handlers.Files.GetSharedFileInfo)
	files.POST("/share/:idFile/user/:idUser", handlers.Files.AddUserToFile)
	files.DELETE("/share/:idFile/user/:idUser", handlers.Files.RemoveUserFromFile)
}
