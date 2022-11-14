package handlers

import (
	"fmt"
	"gosfV2/src/auth"
	"gosfV2/src/dtos"
	"gosfV2/src/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

var Files fileController

type fileController struct{}

// Maneja los errores de los archivos, si el error ErrFileNotFound
// o si es un error desconocido (base de datos), devuelve un error 500
func HandleFileError(err error) error {
	if err == models.ErrFileNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

// Obtiene le path-variable "id" del URL para retornar el archivo
// si el usuario el dueño del archivo
func (f *fileController) GetFile(c echo.Context) error {
	idNum, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	path, err := models.Files(c).GetFilePath(uint(idNum), c.Get("user_id").(uint))
	if err != nil {
		return HandleFileError(err)
	}

	return c.File(path)
}

// Obtiene le path-variable "id" del URL para retornar el archivo
// si el usuario el dueño del archivo
func (f *fileController) GetInfo(c echo.Context) error {
	idNum, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	file, err := models.Files(c).GetByIdFromUser(uint(idNum), c.Get("user_id").(uint))
	if err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, file)
}

// Obtiene todos los archivos del usuario
func (f *fileController) GetAllFiles(c echo.Context) error {
	files, err := models.Files(c).GetAllFromUser(c.Get("user_id").(uint))
	if err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, dtos.ToFileListDTO(files))
}

// Obtiene le path-variable "id" del URL para retornar el archivo
// si el archivo esta compartido con el usuario
func (f *fileController) GetSharedFile(c echo.Context) error {
	idNum, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	file, err := models.Files(c).GetById(uint(idNum))
	if err != nil {
		return HandleFileError(err)
	}

	if file.UserID == c.Get("user_id").(uint) {
		return c.File(models.Files(c).GetPath(file.Filename, file.User.Username))
	}

	// Si esta compartido lo envió directamente
	if file.Shared {
		return c.File(models.Files(c).GetPath(file.Filename, file.User.Username))
	}

	sharedWithMe, err := models.Files(c).IsSharedWith(file.ID, c.Get("user_id").(uint))
	if err != nil {
		return HandleFileError(err)
	}

	// Si el el usuario no es el dueño del archivo, por seguridad
	// le digo que el archivo no existe
	if !sharedWithMe {
		return echo.NewHTTPError(http.StatusNotFound, models.ErrFileNotFound.Error())
	}

	return c.File(models.Files(c).GetPath(file.Filename, file.User.Username))
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

		if err := models.Files(c).Create(f.Filename, c.Get("user_id").(uint), src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := src.Close(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("%d was file/s uploaded successfully", len(mf.File["files"])),
	})
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

	idFile := uint(idNum)
	if _, err = models.Files(c).GetByIdFromUser(idFile, c.Get("user_id").(uint)); err != nil {
		return HandleFileError(err)
	}

	if err := models.Files(c).SetShared(idFile, file.Shared); err != nil {
		return HandleFileError(err)
	}

	if err := models.Files(c).Rename(idFile, file.Filename); err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("File %d updated successfully", idNum),
	})
}

func (f *fileController) DeleteFile(c echo.Context) error {
	idNum, err := strconv.Atoi(c.Param("id"))

	// El QueryParam "force" es opcional, si no viene se asume false
	force := c.QueryParam("force")

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	file, err := models.Files(c).GetByIdFromUser(uint(idNum), c.Get("user_id").(uint))
	if err != nil {
		return HandleFileError(err)
	}

	if len(file.SharedWith) != 0 {
		if force == "yes" {
			for _, user := range file.SharedWith {
				if err := models.Files(c).RemoveUserFromFile(user.ID, file.ID); err != nil {
					return HandleFileError(err)
				}
			}
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "File is shared with other users")
		}
	}

	if err := models.Files(c).Delete(uint(idNum)); err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("File %d deleted successfully", idNum),
	})
}

func (f *fileController) AddUserToFile(c echo.Context) error {

	var userId, fileId uint
	{
		idFileNum, err := strconv.Atoi(c.Param("idFile"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
		}
		fileId = uint(idFileNum)

		idUserNum, err := strconv.Atoi(c.Param("idUser"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
		}
		userId = uint(idUserNum)
	}

	file, err := models.Files(c).GetById(fileId)
	if err != nil {
		return HandleFileError(err)
	}

	// Verifico que exista el usuario con el userId
	if _, err = models.Users(c).FindUserById(userId); err != nil {
		return auth.HandleUserError(err)
	}

	if userId == c.Get("user_id").(uint) {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "You can't share a file with yourself",
		})
	}

	// Si el el usuario no es el dueño del archivo, por seguridad
	// le digo que el archivo no existe
	if file.UserID != c.Get("user_id").(uint) {
		return echo.NewHTTPError(http.StatusNotFound, models.ErrFileNotFound.Error())
	}

	if err := models.Files(c).AddUserToFile(userId, fileId); err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("User %d added to file %d successfully", userId, fileId),
	})
}

func (f *fileController) RemoveUserFromFile(c echo.Context) error {
	var userId, fileId uint
	{
		idFileNum, err := strconv.Atoi(c.Param("idFile"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
		}
		fileId = uint(idFileNum)

		idUserNum, err := strconv.Atoi(c.Param("idUser"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
		}
		userId = uint(idUserNum)
	}

	file, err := models.Files(c).GetById(fileId)
	if err != nil {
		return HandleFileError(err)
	}

	// Verifico que exista el usuario con el userId
	if _, err = models.Users(c).FindUserById(userId); err != nil {
		return auth.HandleUserError(err)
	}

	// Si el el usuario no es el dueño del archivo, por seguridad
	// le digo que el archivo no existe
	if file.UserID != c.Get("user_id").(uint) {
		return echo.NewHTTPError(http.StatusNotFound, models.ErrFileNotFound.Error())
	}

	if err := models.Files(c).RemoveUserFromFile(userId, fileId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("User %d removed from file %d successfully", userId, fileId),
	})
}
