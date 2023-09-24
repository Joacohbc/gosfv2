package handlers

import (
	"gosfV2/src/auth"
	"gosfV2/src/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

var Notes notesController

type notesController struct{}

func (n *notesController) GetNotes(c echo.Context) error {
	userId := auth.Middlewares.GetUserId(c)
	note, err := models.NoteManager.GetNote(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return jsonDTO(c, http.StatusOK, note)
}

func (n *notesController) CreateNotes(c echo.Context) error {
	userId := auth.Middlewares.GetUserId(c)

	var note models.Note
	if err := c.Bind(&note); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	models.NoteManager.CreateNote(userId, note.Content)
	return jsonDTO(c, http.StatusOK, note)
}

// func (n *notesController) UpdateNotes(c echo.Context) error {
// 	userId := auth.Middlewares.GetUserId(c)

// 	var note models.Note
// 	if err := c.Bind(&note); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	note, err := models.NoteManager.UpdateNote(userId, note.Content)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusNotFound, err.Error())
// 	}

// 	return jsonDTO(c, http.StatusOK, note)
// }

// func (n *notesController) DeleteNotes(c echo.Context) error {
// 	userId := auth.Middlewares.GetUserId(c)

// 	note, err := models.NoteManager.GetNote(userId)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusNotFound, err.Error())
// 	}

// 	if err := models.NoteManager.DeleteNote(userId); err != nil {
// 		return echo.NewHTTPError(http.StatusNotFound, err.Error())
// 	}

// 	return jsonDTO(c, http.StatusOK, note)
// }
