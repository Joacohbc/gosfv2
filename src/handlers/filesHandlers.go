package handlers

import (
	"gosfV2/src/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

var Files fileController

type fileController struct{}

func (f fileController) handleFileError(err error) error {
	if err == models.ErrFileNotFound {
		return f.handleFileError(err)
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

func (f *fileController) DeleteFile(c echo.Context) error {
	idNum, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	if err := models.Files(c).Delete(uint(idNum)); err != nil {
		return f.handleFileError(err)
	}

	return c.String(http.StatusOK, "File deleted successfully")
}

func (f *fileController) GetFile(c echo.Context) error {
	idNum, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	file, err := models.Files(c).GetById(uint(idNum))
	if err != nil {
		return f.handleFileError(err)
	}

	return c.File(models.Files(c).GetPath(file.Filename))
}

func (f *fileController) UploadFile(c echo.Context) error {

	mf, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, f := range mf.File["files"] {

		src, err := f.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := models.Files(c).Create(f.Filename, src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := src.Close(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.String(http.StatusOK, "File/s uploaded successfully")
}

func (f *fileController) GetAllFiles(c echo.Context) error {

	files, err := models.Files(c).GetAll()
	if err != nil {
		return f.handleFileError(err)
	}

	return c.JSON(http.StatusOK, files)
}

func (f *fileController) UpdateFile(c echo.Context) error {

	idNum, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	var file models.File
	if err := c.Bind(&file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := models.Files(c).SetShared(uint(idNum), file.Shared); err != nil {
		return f.handleFileError(err)
	}

	if err := models.Files(c).Rename(uint(idNum), file.Filename); err != nil {
		return f.handleFileError(err)
	}

	return c.String(http.StatusOK, "File updated successfully")
}

func (f *fileController) GetSharedFile(c echo.Context) error {
	idNum, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	file, err := models.Files(c).GetFileNotOfOwner(uint(idNum))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if file.Owner.Username == c.Get("username").(string) {
		return c.File(models.Files(c).GetPath(file.Filename))
	}

	if !file.Shared {
		return echo.NewHTTPError(http.StatusForbidden, "File not shared")
	}

	sharedWithMe, err := models.Files(c).IsSharedMe(file.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !sharedWithMe {
		return echo.NewHTTPError(http.StatusForbidden, "File not shared with you")
	}

	return c.File(models.Files(c).GetPathFromUser(file.Owner.Username, file.Filename))
}

func (f *fileController) AddUserToFile(c echo.Context) error {
	idFileNum, err := strconv.Atoi(c.Param("idFile"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	idUserNum, err := strconv.Atoi(c.Param("idUser"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}

	if err := models.Files(c).AddSharedWith(uint(idUserNum), uint(idFileNum)); err != nil {
		return f.handleFileError(err)
	}

	return c.String(http.StatusOK, "User added successfully")
}

func (f *fileController) RemoveUserFromFile(c echo.Context) error {
	idFile := c.Param("idFile")
	idUser := c.Param("idUser")

	idFileNum, err := strconv.Atoi(idFile)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	idUserNum, err := strconv.Atoi(idUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}

	if err := models.Files(c).RemoveSharedWith(uint(idUserNum), uint(idFileNum)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "User removed successfully")
}
