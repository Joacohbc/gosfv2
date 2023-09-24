package routes

import (
	"gosfV2/src/handlers"

	"github.com/labstack/echo/v4"
)

var Notes notesRoutes

type notesRoutes struct{}

// Agrego los Endpoints de Notes
func (a *notesRoutes) AddNotesRoutes(group *echo.Group) {
	notes := group.Group("/notes")
	notes.GET("", handlers.Notes.GetNotes)
	notes.POST("", handlers.Notes.CreateNotes)
	// notes.PUT("", handlers.Notes.UpdateNotes)
	// notes.DELETE("", handlers.Notes.DeleteNotes)
}
