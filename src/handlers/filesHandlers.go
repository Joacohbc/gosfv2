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

func (f *fileController) GetFile(c echo.Context) error {
	username := c.Get("username").(string)
	filename := c.Param("filename")

	file, err := models.Files.GetFileByUsername(c.Request().Context(), username, filename)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, file)
}

func (f *fileController) GetAllFiles(c echo.Context) error {
	username := c.Get("username").(string)

	files, err := models.Files.GetAllFilesByUsername(c.Request().Context(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, files)
}

func (f *fileController) ShareFile(c echo.Context) error {
	username := c.Get("username").(string)
	filename := c.Param("filename")
	share := c.QueryParam("share")

	if share == "yes" {
		err := models.Files.SetShared(c.Request().Context(), username, filename, true)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "File shared successfully")
	} else {
		err := models.Files.SetShared(c.Request().Context(), username, filename, false)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "File not shared successfully")
	}
}

func (f *fileController) GetShareFile(c echo.Context) error {
	username := c.QueryParam("username")
	filename := c.QueryParam("filename")

	file, err := models.Files.GetFileByUsername(c.Request().Context(), username, filename)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !file.Shared {
		return echo.NewHTTPError(http.StatusForbidden, "File not shared")
	}

	return c.File(models.Files.GetPathUser(username, filename))
}
