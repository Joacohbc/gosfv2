package handlers

import (
	"gosfV2/src/models/files"
	"net/http"

	"github.com/labstack/echo"
)

func UploadFile(c echo.Context) error {

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

		if err := files.CreateFile(c.Request().Context(), f.Filename, username, src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := src.Close(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.String(http.StatusOK, "File/s uploaded successfully")
}

func DeleteFile(c echo.Context) error {
	username := c.Get("username").(string)
	fileName := c.Param("filename")

	if err := files.DeleteFile(c.Request().Context(), username, fileName); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "File deleted successfully")
}

func GetAllFiles(c echo.Context) error {
	username := c.Get("username").(string)
	files, err := files.GetAllFilesByUsername(c.Request().Context(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, files)
}

func GetFile(c echo.Context) error {
	username := c.Get("username").(string)
	fileName := c.Param("filename")

	file, err := files.GetFileByUsername(c.Request().Context(), username, fileName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.File(files.GetPathUser(username, file.Filename))
}
