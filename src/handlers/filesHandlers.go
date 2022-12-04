package handlers

import (
	"fmt"
	"gosfV2/src/auth"
	"gosfV2/src/dtos"
	"gosfV2/src/models"
	"gosfV2/src/utils"
	"net/http"
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
func (f *fileController) GetIdFromURL(c echo.Context, param string) (uint, error) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Invalid Id from URL")
	}

	return uint(id), nil
}

// Obtiene el Id del URL para retornar el archivo
// SI el usuario logueado es el dueño del archivo
//
// PathParams:
// - Id de Archivo | Uint
func (f *fileController) GetFile(c echo.Context) error {
	idFile, err := f.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	path, err := models.Files(c).GetFilepath(idFile, c.Get("user_id").(uint))
	if err != nil {
		return HandleFileError(err)
	}

	return c.File(path)
}

// Obtiene el Id del URL para retornar la Información del archivo
// SI el usuario logueado es el dueño del archivo
//
// PathParams:
// - Id de Archivo | Uint
func (f *fileController) GetInfo(c echo.Context) error {
	idNum, err := f.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetByIdFromUser(idNum, c.Get("user_id").(uint))
	if err != nil {
		return HandleFileError(err)
	}

	return c.JSON(http.StatusOK, dtos.ToFileDTO(file))
}

// Obtiene todos los archivos del usuario (Su información)
func (f *fileController) GetAllFiles(c echo.Context) error {
	files, err := models.Files(c).GetAllFromUser(c.Get("user_id").(uint))
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
func (f *fileController) GetSharedFile(c echo.Context) error {
	idFile, err := f.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetById(idFile)
	if err != nil {
		return HandleFileError(err)
	}

	idCurrentUser := c.Get("user_id").(uint)

	// Si el usuario es el dueño del archivo, lo envío directamente
	if file.UserID == idCurrentUser {
		return c.File(models.Files(c).GetPath(file.Filename, file.User.Username))
	}

	// Si esta compartido lo envió directamente
	if file.Shared {
		return c.File(models.Files(c).GetPath(file.Filename, file.User.Username))
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

	return c.File(models.Files(c).GetPath(file.Filename, file.User.Username))
}

// Obtiene los archivos subidos desde el Body de la Request (Formulario)
// y los guarda en la base de datos para el usuario logueado
//
// PathParams:
// - Id de Archivo | Uint
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

	return c.JSON(http.StatusOK, utils.ToJSON(fmt.Sprintf("%d was file/s uploaded successfully", len(mf.File["files"]))))
}

// Obtiene los archivos subidos desde el Body (Formulario)
// y los guarda en la base de datos para el usuario logueado
//
// PathParams:
// - Id de Archivo | Uint
func (f *fileController) UpdateFile(c echo.Context) error {

	idFile, err := f.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	var file models.File
	if err := c.Bind(&file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err = models.Files(c).GetByIdFromUser(idFile, c.Get("user_id").(uint)); err != nil {
		return HandleFileError(err)
	}

	if err := models.Files(c).SetShared(idFile, file.Shared); err != nil {
		return HandleFileError(err)
	}

	if err := models.Files(c).Rename(idFile, file.Filename); err != nil {
		return HandleFileError(err)
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
func (f *fileController) DeleteFile(c echo.Context) error {
	idNum, err := f.GetIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetByIdFromUser(idNum, c.Get("user_id").(uint))
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
func (f *fileController) AddUserToFile(c echo.Context) error {

	fileId, err := f.GetIdFromURL(c, "idFile")
	if err != nil {
		return err
	}

	userId, err := f.GetIdFromURL(c, "idUser")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetByIdFromUser(fileId, c.Get("user_id").(uint))
	if err != nil {
		return HandleFileError(err)
	}

	// Verifico que exista el usuario con el userId
	user, err := models.Users(c).FindUserById(userId)
	if err != nil {
		return auth.HandleUserError(err)
	}

	for _, u := range file.SharedWith {
		if u.Equals(user) {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
				"message": "The File is already shared with that user",
			})
		}
	}

	if userId == file.UserID {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "You can't share a file with yourself",
		})
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
func (f *fileController) RemoveUserFromFile(c echo.Context) error {
	fileId, err := f.GetIdFromURL(c, "idFile")
	if err != nil {
		return err
	}

	userId, err := f.GetIdFromURL(c, "idUser")
	if err != nil {
		return err
	}

	file, err := models.Files(c).GetByIdFromUser(fileId, c.Get("user_id").(uint))
	if err != nil {
		return HandleFileError(err)
	}

	// Verifico que exista el usuario con el userId
	user, err := models.Users(c).FindUserById(userId)
	if err != nil {
		return auth.HandleUserError(err)
	}

	flagFind := false
	for _, u := range file.SharedWith {
		if u.Equals(user) {
			flagFind = true
			break
		}
	}

	if !flagFind {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "The File is not shared with that user",
		})
	}

	// Si el el usuario no es el dueño del archivo, por seguridad
	// le digo que el archivo no existe
	if file.UserID != c.Get("user_id").(uint) {
		return echo.NewHTTPError(http.StatusNotFound, models.ErrFileNotFound.Error())
	}

	if err := models.Files(c).RemoveUserFromFile(userId, fileId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.ToJSON(fmt.Sprintf("User %d removed from file %d successfully", userId, fileId)))
}
