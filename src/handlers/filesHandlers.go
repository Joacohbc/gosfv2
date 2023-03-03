package handlers

import (
	"fmt"
	"gosfV2/src/auth"
	"gosfV2/src/dtos"
	"gosfV2/src/models"
	"gosfV2/src/utils"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo"
)

var Files fileController

type fileController struct{}

// Maneja los errores de los archivos:
// - Si el error ErrFileNotFound devuelve un error 404
// - Si es un error desconocido (base de datos), devuelve un error 500
func HandleFileError(err error) error {
	if err == models.ErrFileNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

// Obtiene el ID como PathParam y lo convierte en uint
func (fc *fileController) GetIdFromURL(c echo.Context, param string) (uint, error) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Invalid Id from URL")
	}

	return uint(id), nil
}

// Obtiene el ID del usuario logueado del contexto
func (fc *fileController) getUserId(c echo.Context) uint {
	return c.Get("user_id").(uint)
}

// Obtiene el Id del URL para retornar el archivo
// Si el usuario logueado es el dueño del archivo
//
// PathParams:
// - Id de Archivo | Uint
func (fc *fileController) GetFile(c echo.Context) error {
	idFile, err := fc.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetByIdFromUser(idFile, fc.getUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	return c.Inline(file.GetPath(), file.Filename)
}

// Obtiene el Id del URL para retornar la Información del archivo
// SI el usuario logueado es el dueño del archivo
//
// PathParams:
// - Id de Archivo | Uint
func (fc *fileController) GetInfo(c echo.Context) error {
	idNum, err := fc.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetByIdFromUser(idNum, fc.getUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, dtos.ToFileDTO(file))
}

// Obtiene todos los archivos del usuario (Su información)
func (fc *fileController) GetAllFiles(c echo.Context) error {
	files, err := models.Files(c).GetFilesFromUser(fc.getUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, dtos.ToFileListDTO(files))
}

// Obtiene el Id del URL para retornar el archivo
// SI el archivo esta compartido con el Usuario logueado
//
// PathParams:
// - Id de Archivo | Uint
func (fc *fileController) GetSharedFile(c echo.Context) error {
	idFile, err := fc.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetById(idFile)
	if err != nil {
		return HandleFileError(err)
	}

	idCurrentUser := fc.getUserId(c)

	// Si el usuario es el dueño del archivo, lo envío directamente
	if file.UserID == idCurrentUser {
		return c.File(file.GetPath())
	}

	// Si esta compartido lo envió directamente
	if file.Shared {
		return c.File(file.GetPath())
	}

	sharedWithMe, err := models.Files(c).IsSharedWith(file.ID, idCurrentUser)
	if err != nil {
		return HandleFileError(err)
	}

	// Si el el usuario no es el dueño del archivo, por seguridad
	// le digo que el archivo no existe
	if !sharedWithMe {
		return echo.NewHTTPError(http.StatusNotFound, models.ErrFileNotFound.Error())
	}

	return c.File(file.GetPath())
}

// Obtiene todos los archivos del usuario (Su información)
func (fc *fileController) GetAllShareFiles(c echo.Context) error {
	files, err := models.Files(c).GetFilesShared(fc.getUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, dtos.ToFileListDTO(files))
}

// Obtiene los archivos subidos desde el Body de la Request (Formulario)
// y los guarda en la base de datos para el usuario logueado
func (fc *fileController) UploadFile(c echo.Context) error {
	mf, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, file := range mf.File["files"] {
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := models.Files(c).Create(file.Filename, fc.getUserId(c), src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := src.Close(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, utils.ToJSON(fmt.Sprintf("%d was file/s uploaded successfully", len(mf.File["files"]))))
}

// Obtiene los archivos subidos desde el Body (Formulario)
// y los guarda en la base de datos para el usuario logueado
//
// PathParams:
// - Id de Archivo | Uint
func (fc *fileController) UpdateFile(c echo.Context) error {

	idFile, err := fc.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	var file dtos.FileDTO
	if err := c.Bind(&file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file data")
	}

	actual, err := models.Files(c).GetByIdFromUser(idFile, fc.getUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	// Si el nombre del archivo es diferente al actual, lo renombro
	if actual.Filename != file.Filename {

		if filepath.Ext(file.Filename) != filepath.Ext(actual.Filename) {
			return echo.NewHTTPError(http.StatusBadRequest, "The extension of the file cannot be changed")
		}

		if err := models.Files(c).Rename(idFile, file.Filename); err != nil {
			return HandleFileError(err)
		}
	}

	if actual.Shared != file.Shared {
		if err := models.Files(c).SetShared(idFile, file.Shared); err != nil {
			return HandleFileError(err)
		}
	}

	return c.JSON(http.StatusOK, utils.ToJSON(fmt.Sprintf("File %d updated successfully", idFile)))
}

// Elimina el archivo de la base de datos y del disco
// SI el usuario logueado es el dueño del archivo
//
// PathParams:
// - Id de Archivo | Uint
//
// QueryParams:
// - Force | String
func (fc *fileController) DeleteFile(c echo.Context) error {
	idNum, err := fc.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetByIdFromUser(idNum, fc.getUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	// El QueryParam "force" es opcional, si no viene se asume false
	force := c.QueryParam("force")

	// Verifico si el archivo esta compartido con otros usuarios
	if len(file.SharedWith) != 0 {

		// Si viene el QueryParam "force" y es true, elimino el archivo (aunque este compartido)
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

	return c.JSON(http.StatusOK, utils.ToJSON(fmt.Sprintf("File %d deleted successfully", idNum)))
}

// Agrega un usuario a la lista de usuarios con los que se comparte el archivo
//
// PathParams:
// - Id de Archivo | Uint
// - Id de Usuario | Uint
func (fc *fileController) AddUserToFile(c echo.Context) error {

	fileId, err := fc.GetIdFromURL(c, "idFile")
	if err != nil {
		return err
	}

	userId, err := fc.GetIdFromURL(c, "idUser")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetByIdFromUser(fileId, fc.getUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	if file.UserID == userId {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ToJSON("The user is the owner of the file"))
	}

	// Verifico que exista el usuario con el userId
	if _, err := models.Users(c).FindUserById(userId); err != nil {
		return auth.HandleUserError(err)
	}

	ok, err := models.Files(c).IsSharedWith(fileId, userId)
	if err != nil {
		return HandleFileError(err)
	}

	if ok {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ToJSON("The File is already shared with that user"))
	}

	if err := models.Files(c).AddUserToFile(userId, fileId); err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, utils.ToJSON(fmt.Sprintf("User %d added to file %d successfully", userId, fileId)))
}

// Remueve un usuario a la lista de usuarios con los que se comparte el archivo
//
// PathParams:
// - Id de Archivo | Uint
// - Id de Usuario | Uint
func (fc *fileController) RemoveUserFromFile(c echo.Context) error {
	fileId, err := fc.GetIdFromURL(c, "idFile")
	if err != nil {
		return err
	}

	userId, err := fc.GetIdFromURL(c, "idUser")
	if err != nil {
		return err
	}

	// Verifico que exista el archivo con el fileId
	file, err := models.Files(c).GetByIdFromUser(fileId, fc.getUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	// Verifico que exista el usuario con el userId
	if _, err := models.Users(c).FindUserById(userId); err != nil {
		return auth.HandleUserError(err)
	}

	ok, err := models.Files(c).IsSharedWith(fileId, userId)
	if err != nil {
		return HandleFileError(err)
	}

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ToJSON("The File is not shared with that user"))
	}

	// Si el el usuario no es el dueño del archivo, por seguridad
	// le digo que el archivo no existe
	if file.UserID != fc.getUserId(c) {
		return echo.NewHTTPError(http.StatusNotFound, models.ErrFileNotFound.Error())
	}

	if err := models.Files(c).RemoveUserFromFile(userId, fileId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.ToJSON(fmt.Sprintf("User %d removed from file %d successfully", userId, fileId)))
}
