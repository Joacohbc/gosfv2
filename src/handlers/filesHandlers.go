package handlers

import (
	"gosfV2/src/models"
	"net/http"

	"github.com/labstack/echo"
)

var Files fileController

type fileController struct{}

func (f *fileController) UploadFile(c echo.Context) error {

	mf, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	username := c.Get("username").(string)
	for _, f := range mf.File["files"] {

		src, err := f.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := models.Files.CreateFile(c.Request().Context(), f.Filename, username, src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := src.Close(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.String(http.StatusOK, "File/s uploaded successfully")
}

func (f *fileController) DeleteFile(c echo.Context) error {
	username := c.Get("username").(string)
	fileName := c.Param("filename")

	if err := models.Files.DeleteFile(c.Request().Context(), username, fileName); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "File deleted successfully")
}

func (f *fileController) GetFiles(c echo.Context) error {
	username := c.Get("username").(string)

	if filename := c.QueryParam("filename"); filename != "" {

		file, err := models.Files.GetFileByUsername(c.Request().Context(), username, filename)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, file)

	}

	files, err := models.Files.GetAllFilesByUsername(c.Request().Context(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, files)
}
