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

	files.GET("/", handlers.Files.GetAllFiles)
	files.GET("/:filename", handlers.Files.GetFile)
	files.POST("/", handlers.Files.UploadFile)
	files.DELETE("/:filename", handlers.Files.DeleteFile)

	files.GET("/share", handlers.Files.GetShareFile)
	files.PUT("/share/:filename", handlers.Files.ShareFile)
}
