package handlers

import (
	"gosfV2/src/auth"
	"gosfV2/src/dtos"
	"gosfV2/src/models"
	"net/http"

	"github.com/labstack/echo"
)

var Users userController

type userController struct{}

func (u userController) RenameUser(c echo.Context) error {

	var user dtos.UserDTO
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
	}

	exist, err := models.Users(c).ExistUserByName(user.Username)
	if err != nil {
		return auth.HandleUserError(err)
	}

	if exist {
		return echo.NewHTTPError(http.StatusBadRequest, "Username already exists")
	}

	if err := models.Users(c).Rename(c.Get("user_id").(uint), user.Username); err != nil {
		return auth.HandleUserError(err)
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"message": "User renamed successfully",
	})
}

func (u userController) ChangePassword(c echo.Context) error {

	type PasswordDTO struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var password PasswordDTO

	if err := c.Bind(&password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
	}

	user, err := models.Users(c).FindUserById(c.Get("user_id").(uint))
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
	if err := models.Users(c).ChangePassword(c.Get("user_id").(uint), password.NewPassword); err != nil {
		return auth.HandleUserError(err)
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"message": "Password changed successfully",
	})
}

func (u userController) DeleteUser(c echo.Context) error {

	files, err := models.Files(c).GetFilesFromUser(c.Get("user_id").(uint))
	if err != nil {
		return auth.HandleUserError(err)
	}

	if len(files) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "You can't delete your account because you have files")
	}

	if err := models.Users(c).Delete(c.Get("user_id").(uint)); err != nil {
		return auth.HandleUserError(err)
	}

	if err := auth.NewTokenManager().RemoveUserTokens(c.Get("user_id").(uint)); err != nil {
		return auth.HandlerTokenError(err)
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"message": "User deleted successfully",
	})
}

func (u userController) GetUser(c echo.Context) error {
	user, err := models.Users(c).FindUserById(c.Get("user_id").(uint))
	if err != nil {
		return auth.HandleUserError(err)
	}

	return c.JSON(http.StatusOK, dtos.ToUserDTO(user))
}
