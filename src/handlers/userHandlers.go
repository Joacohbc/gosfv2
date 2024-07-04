package handlers

import (
	"fmt"
	"gosfV2/src/auth"
	"gosfV2/src/auth/jwt"
	"gosfV2/src/dtos"
	"gosfV2/src/models"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

var Users userController

type userController struct{}

// Obtiene el usuario actual
func (u userController) GetUser(c echo.Context) error {
	user, err := models.Users().FindUserById(auth.Middlewares.GetUserId(c))
	if err != nil {
		return auth.HandleUserError(err)
	}

	return jsonDTO(c, http.StatusOK, user)
}

func generateUserIcon(username string) string {
	size := 512
	backgroundColor := "ffffff"
	fontColor := "000000"

	userNameDivided := strings.Split(username, " ")

	firstLetter := strings.ToUpper(string(userNameDivided[0][0]))
	secondLetter := ""

	if len(userNameDivided) > 1 && len(userNameDivided[1]) > 0 {
		secondLetter = strings.ToUpper(string(userNameDivided[1][0]))
	}

	icon := fmt.Sprintf(`
    <svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 250 250">
        <g id="icon">
            <circle cx="125" cy="125" r="125" fill="#%s"/>
            <text x="50%%" y="54%%" fill="#%s" font-size="110" style="font-family:monospace;" dominant-baseline="middle" text-anchor="middle">%s</text>
        </g>
    </svg>`, size, size, backgroundColor, fontColor, firstLetter+secondLetter)

	// return []byte(icon)
	return icon
}

// Obtiene el icono del usuario actual
func (u userController) GetIcon(c echo.Context) error {
	path := models.Users().GetIcon(auth.Middlewares.GetUserId(c))
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {

			file, err := models.Users().FindUserById(auth.Middlewares.GetUserId(c))
			if err != nil {
				return c.File(models.DefaultIcon)
			}

			// reader := strings.NewReader(generateUserIcon(file.Username))
			return c.Blob(http.StatusOK, "image/svg+xml", []byte(generateUserIcon(file.Username)))
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

	path := models.Users().GetIcon(id)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {

			file, err := models.Users().FindUserById(id)
			if err != nil {
				return c.File(models.DefaultIcon)
			}

			// reader := strings.NewReader(generateUserIcon(file.Username))
			return c.Blob(http.StatusOK, "image/svg+xml", []byte(generateUserIcon(file.Username)))
		}
		return echo.NewHTTPError(http.StatusNotFound, "Icon not found")
	}

	return c.File(path)
}

// Cambia el nombre de usuario actual
//
// Body:
// - Username | String
func (u userController) RenameUser(c echo.Context) error {

	var user dtos.UserDTO
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
	}

	if *user.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username must not be empty")
	}

	exist, err := models.Users().ExistUserByName(*user.Username)
	if err != nil {
		return auth.HandleUserError(err)
	}

	if exist {
		return echo.NewHTTPError(http.StatusBadRequest, "Username already exists")
	}

	if err := models.Users().Rename(auth.Middlewares.GetUserId(c), *user.Username); err != nil {
		return auth.HandleUserError(err)
	}

	userUpdated, err := models.Users().FindUserById(auth.Middlewares.GetUserId(c))
	if err != nil {
		return auth.HandleUserError(err)
	}

	return jsonDTO(c, http.StatusOK, userUpdated)
}

// Cambia la contraseña del usuario actual
//
// Body:
// - OldPassword | String
// - NewPassword | String
func (u userController) ChangePassword(c echo.Context) error {

	var password struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.Bind(&password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
	}

	user, err := models.Users().FindUserById(auth.Middlewares.GetUserId(c))
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
	if err := models.Users().ChangePassword(auth.Middlewares.GetUserId(c), password.NewPassword); err != nil {
		return auth.HandleUserError(err)
	}

	return c.NoContent(http.StatusOK)
}

// Elimina el usuario actual
func (u userController) DeleteUser(c echo.Context) error {

	files, err := models.Files().GetFilesFromUser(auth.Middlewares.GetUserId(c))
	if err != nil {
		return auth.HandleUserError(err)
	}

	if len(files) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "You can't delete your account because you have files")
	}

	if err := models.Users().Delete(auth.Middlewares.GetUserId(c)); err != nil {
		return auth.HandleUserError(err)
	}

	if err := jwt.TokenManager.RemoveUserTokens(auth.Middlewares.GetUserId(c)); err != nil {
		return jwt.HandlerTokenError(err)
	}

	return c.NoContent(http.StatusOK)
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

	if err := models.Users().UploadIcon(auth.Middlewares.GetUserId(c), blob); err != nil {

		if err == models.ErrIconFormatNotSupported || err == models.ErrIconTooLarge {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// Elimina el icono del usuario actual
func (u userController) DeleteIcon(c echo.Context) error {

	id := auth.Middlewares.GetUserId(c)
	path := models.Users().GetIcon(id)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, "The user doesn't have an icon")
		}
	}

	err := models.Users().DeleteIcon(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
