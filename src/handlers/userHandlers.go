package handlers

import (
	"gosfV2/src/auth"
	"gosfV2/src/auth/jwt"
	"gosfV2/src/dtos"
	"gosfV2/src/models"
	"gosfV2/src/utils"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

var Users userController

type userController struct{}

// Cambia el nombre de usuario actual
//
// Body:
// - Username | String
func (u userController) RenameUser(c echo.Context) error {

	var user dtos.UserDTO
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
	}

	if user.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username must not be empty")
	}

	exist, err := models.Users(c).ExistUserByName(user.Username)
	if err != nil {
		return auth.HandleUserError(err)
	}

	if exist {
		return echo.NewHTTPError(http.StatusBadRequest, "Username already exists")
	}

	if err := models.Users(c).Rename(auth.Middlewares.GetUserId(c), user.Username); err != nil {
		return auth.HandleUserError(err)
	}

	return c.JSON(http.StatusAccepted, utils.ToJSON("User renamed successfully"))
}

// Cambia la contraseña del usuario actual
//
// Body:
// - OldPassword | String
// - NewPassword | String
func (u userController) ChangePassword(c echo.Context) error {

	type PasswordDTO struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var password PasswordDTO

	if err := c.Bind(&password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
	}

	user, err := models.Users(c).FindUserById(auth.Middlewares.GetUserId(c))
	if err != nil {
		return auth.HandleUserError(err)
	}

	// Verifico que la contraseña vieja sea correcta
	ok, err := auth.CheckPassword(password.OldPassword, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Current password is incorrect")
	}

	// Genero el hash de la contraseña nueva
	if err := auth.GeneratePassword(&password.NewPassword); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Actualizo la contraseña
	if err := models.Users(c).ChangePassword(auth.Middlewares.GetUserId(c), password.NewPassword); err != nil {
		return auth.HandleUserError(err)
	}

	return c.JSON(http.StatusAccepted, utils.ToJSON("Password changed successfully"))
}

// Elimina el usuario actual
func (u userController) DeleteUser(c echo.Context) error {

	files, err := models.Files(c).GetFilesFromUser(auth.Middlewares.GetUserId(c))
	if err != nil {
		return auth.HandleUserError(err)
	}

	if len(files) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "You can't delete your account because you have files")
	}

	if err := models.Users(c).Delete(auth.Middlewares.GetUserId(c)); err != nil {
		return auth.HandleUserError(err)
	}

	if err := jwt.TokenManager.RemoveUserTokens(auth.Middlewares.GetUserId(c)); err != nil {
		return jwt.HandlerTokenError(err)
	}

	return c.JSON(http.StatusAccepted, utils.ToJSON("User deleted successfully"))
}

// Obtiene el usuario actual
func (u userController) GetUser(c echo.Context) error {
	user, err := models.Users(c).FindUserById(auth.Middlewares.GetUserId(c))
	if err != nil {
		return auth.HandleUserError(err)
	}

	return c.JSON(http.StatusOK, dtos.ToUserDTO(user))
}

// Obtiene el icono del usuario actual
func (u userController) GetIcon(c echo.Context) error {

	path := models.Users(c).GetIcon(auth.Middlewares.GetUserId(c))

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return c.File(models.DefaultIcon)
		}
		return echo.NewHTTPError(http.StatusNotFound, "Icon not found")
	}

	return c.File(path)
}

// Obtiene el icono del usuario con el id especificado
//
// Params:
// - userId | String
func (u userController) GetIconFromUser(c echo.Context) error {

	id, err := getIdFromURL(c, "userId")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}

	path := models.Users(c).GetIcon(id)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return c.File(models.DefaultIcon)
		}
		return echo.NewHTTPError(http.StatusNotFound, "Icon not found")
	}

	return c.File(path)
}

// Sube un icono para el usuario actual
//
// Body:
// - icon | File
func (u userController) UploadIcon(c echo.Context) error {

	file, err := c.FormFile("icon")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file: "+err.Error())
	}

	blob, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := models.Users(c).UploadIcon(auth.Middlewares.GetUserId(c), blob); err != nil {

		if err == models.ErrIconFormatNotSupported || err == models.ErrIconTooLarge {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.ToJSON("Icon uploaded successfully"))
}

// Elimina el icono del usuario actual
func (u userController) DeleteIcon(c echo.Context) error {

	id := auth.Middlewares.GetUserId(c)
	path := models.Users(c).GetIcon(id)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, "The user doesn't have an icon")
		}
	}

	err := models.Users(c).DeleteIcon(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.ToJSON("Icon removed successfully"))

}
