package handlers

import (
	"gosfV2/src/auth"
	"gosfV2/src/dtos"
	"gosfV2/src/models"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
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

// Obtiene el Id del URL para retornar el archivo
// Si el usuario logueado es el dueño del archivo
//
// PathParams:
// - Id de Archivo | Uint
func (fc *fileController) GetFile(c echo.Context) error {
	idFile, err := getIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files().GetByIdFromUser(idFile, auth.Middlewares.GetUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	return c.Inline(models.Files().GetPath(file.ID, file.Filename), file.Filename)
}

// Obtiene el Id del URL para retornar la información del archivo
// Si el usuario logueado es el dueño del archivo
//
// PathParams:
// - Id de Archivo | Uint
func (fc *fileController) GetInfo(c echo.Context) error {
	idNum, err := getIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files().GetByIdFromUser(idNum, auth.Middlewares.GetUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	return jsonDTO(c, http.StatusOK, file)
}

// Obtiene todos los archivos del usuario (Su información)
func (fc *fileController) GetAllFiles(c echo.Context) error {
	files, err := models.Files().GetFilesFromUser(auth.Middlewares.GetUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	return jsonDTO(c, http.StatusOK, files)
}

// Obtiene el Id del URL para retornar el archivo
// SI el archivo esta compartido con el Usuario logueado
//
// PathParams:
// - Id de Archivo | Uint
func (fc *fileController) GetSharedFile(c echo.Context) error {
	idFile, err := getIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files().GetById(idFile)
	if err != nil {
		return HandleFileError(err)
	}

	idCurrentUser := auth.Middlewares.GetUserId(c)

	// Si el usuario es el dueño del archivo, lo envío directamente
	if file.Edges.Owner.ID == idCurrentUser {
		return c.File(models.Files().GetPath(file.ID, file.Filename))
	}

	// Si esta compartido lo envió directamente
	if file.IsShared {
		return c.File(models.Files().GetPath(file.ID, file.Filename))
	}

	sharedWithMe, err := models.Files().IsSharedWith(file.ID, idCurrentUser)
	if err != nil {
		return HandleFileError(err)
	}

	if !sharedWithMe {
		return echo.NewHTTPError(http.StatusForbidden, "The file is not shared with you, request access to the owner")
	}

	return c.File(models.Files().GetPath(file.ID, file.Filename))
}

// Obtiene el Id del URL para retornar la informacion del archivo
// SI el archivo esta compartido con el Usuario logueado
//
// PathParams:
// - Id de Archivo | Uint
func (fc *fileController) GetSharedFileInfo(c echo.Context) error {
	idFile, err := getIdFromURL(c, "id")
	if err != nil {
		return err
	}

	file, err := models.Files().GetById(idFile)
	if err != nil {
		return HandleFileError(err)
	}

	idCurrentUser := auth.Middlewares.GetUserId(c)

	// Si el usuario es el dueño del archivo, lo envío directamente
	if file.Edges.Owner.ID == idCurrentUser {
		return jsonDTO(c, http.StatusOK, file)
	}

	// Si esta compartido lo envió directamente
	if file.IsShared {
		return jsonDTO(c, http.StatusOK, file)
	}

	sharedWithMe, err := models.Files().IsSharedWith(file.ID, idCurrentUser)
	if err != nil {
		return HandleFileError(err)
	}

	// Si el el usuario no es el dueño del archivo, por seguridad
	// le digo que el archivo no existe
	if !sharedWithMe {
		return echo.NewHTTPError(http.StatusForbidden, "The file is not shared with you, request access to the owner")
	}

	return jsonDTO(c, http.StatusOK, file)
}

// Obtiene todos los archivos compartidos con el usuario logueado
func (fc *fileController) GetAllShareFiles(c echo.Context) error {
	files, err := models.Files().GetFilesShared(auth.Middlewares.GetUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	return jsonDTO(c, http.StatusOK, files)
}

// Obtiene los archivos subidos desde el Body de la Request (Formulario)
// y los guarda en la base de datos (y en file system) para el usuario logueado
func (fc *fileController) UploadFile(c echo.Context) error {
	mf, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	filesIds := make([]uint, len(mf.File["files"]))

	for _, file := range mf.File["files"] {
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		file, err := models.Files().Create(file.Filename, auth.Middlewares.GetUserId(c), src)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := src.Close(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		filesIds = append(filesIds, file.ID)
	}

	filesCreated, err := models.Files().GetByIds(filesIds)
	if err != nil {
		return HandleFileError(err)
	}

	return jsonDTO(c, http.StatusCreated, filesCreated)
}

// Obtiene los archivos subidos desde el Body (Formulario)
// y los guarda en la base de datos para el usuario logueado
//
// PathParams:
// - Id de Archivo | Uint
func (fc *fileController) UpdateFile(c echo.Context) error {

	idFile, err := getIdFromURL(c, "id")
	if err != nil {
		return err
	}

	var file dtos.FileDTO
	if err := c.Bind(&file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file data")
	}

	actual, err := models.Files().GetByIdFromUser(idFile, auth.Middlewares.GetUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	// Si el nombre del archivo es diferente al actual, lo renombro
	if file.Filename != nil && actual.Filename != *file.Filename {
		if filepath.Ext(*file.Filename) != filepath.Ext(actual.Filename) {
			return echo.NewHTTPError(http.StatusBadRequest, "The extension of the file cannot be changed")
		}

		if _, err := models.Files().Rename(idFile, *file.Filename); err != nil {
			return HandleFileError(err)
		}
	}

	if file.Shared != nil && actual.IsShared != *file.Shared {
		if _, err := models.Files().SetShared(idFile, *file.Shared); err != nil {
			return HandleFileError(err)
		}
	}

	fileUpdated, err := models.Files().GetById(idFile)
	if err != nil {
		return HandleFileError(err)
	}

	return jsonDTO(c, http.StatusOK, fileUpdated)
}

// Elimina el archivo de la base de datos y del file system
// SI el usuario logueado es el dueño del archivo
//
// PathParams:
// - Id de Archivo | Uint
//
// QueryParams:
// - Force | String
func (fc *fileController) DeleteFile(c echo.Context) error {
	idNum, err := getIdFromURL(c, "id")
	if err != nil {
		return err
	}

	fileDeleted, err := models.Files().GetByIdFromUser(idNum, auth.Middlewares.GetUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	// El QueryParam "force" es opcional, si no viene se asume false
	force := c.QueryParam("force")

	// Verifico si el archivo esta compartido con otros usuarios
	if len(fileDeleted.Edges.SharedWith) != 0 {

		// Si viene el QueryParam "force" y es true, elimino el archivo (aunque este compartido)
		if force == "yes" {
			for _, user := range fileDeleted.Edges.SharedWith {
				if err := models.Files().RemoveUserFromFile(user.ID, fileDeleted.ID); err != nil {
					return HandleFileError(err)
				}
			}
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "File is shared with other users")
		}
	}

	if _, err := models.Files().Delete(uint(idNum)); err != nil {
		return HandleFileError(err)
	}

	return jsonDTO(c, http.StatusOK, fileDeleted)
}

// Agrega un usuario a la lista de usuarios con los que se comparte el archivo
//
// PathParams:
// - Id de Archivo | Uint
// - Id de Usuario | Uint
func (fc *fileController) AddUserToFile(c echo.Context) error {

	fileId, err := getIdFromURL(c, "idFile")
	if err != nil {
		return err
	}

	userId, err := getIdFromURL(c, "idUser")
	if err != nil {
		return err
	}

	file, err := models.Files().GetByIdFromUser(fileId, auth.Middlewares.GetUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	if file.Edges.Owner.ID == userId {
		return echo.NewHTTPError(http.StatusBadRequest, "The user is the owner of the file")
	}

	// Verifico que exista el usuario con el userId
	if _, err := models.Users().FindUserById(userId); err != nil {
		return auth.HandleUserError(err)
	}

	ok, err := models.Files().IsSharedWith(fileId, userId)
	if err != nil {
		return HandleFileError(err)
	}

	if ok {
		return echo.NewHTTPError(http.StatusBadRequest, "The File is already shared with that user")
	}

	if err := models.Files().AddUserToFile(userId, fileId); err != nil {
		return HandleFileError(err)
	}

	fileUpdated, err := models.Files().GetById(fileId)
	if err != nil {
		return HandleFileError(err)
	}

	return jsonDTO(c, http.StatusOK, fileUpdated)
}

// Remueve un usuario a la lista de usuarios con los que se comparte el archivo
//
// PathParams:
// - Id de Archivo | Uint
// - Id de Usuario | Uint
func (fc *fileController) RemoveUserFromFile(c echo.Context) error {
	fileId, err := getIdFromURL(c, "idFile")
	if err != nil {
		return err
	}

	userId, err := getIdFromURL(c, "idUser")
	if err != nil {
		return err
	}

	// Verifico que exista el archivo con el fileId
	file, err := models.Files().GetByIdFromUser(fileId, auth.Middlewares.GetUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	// Verifico que exista el usuario con el userId
	if _, err := models.Users().FindUserById(userId); err != nil {
		return auth.HandleUserError(err)
	}

	ok, err := models.Files().IsSharedWith(fileId, userId)
	if err != nil {
		return HandleFileError(err)
	}

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "The File is not shared with that user")
	}

	// Si el el usuario no es el dueño del archivo, por seguridad
	// le digo que el archivo no existe
	if file.Edges.Owner.ID != auth.Middlewares.GetUserId(c) {
		return echo.NewHTTPError(http.StatusNotFound, models.ErrFileNotFound.Error())
	}

	if err := models.Files().RemoveUserFromFile(userId, fileId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	fileUpdated, err := models.Files().GetByIdFromUser(fileId, auth.Middlewares.GetUserId(c))
	if err != nil {
		return HandleFileError(err)
	}

	return jsonDTO(c, http.StatusOK, fileUpdated)
}
