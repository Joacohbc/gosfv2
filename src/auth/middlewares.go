package auth

import (
	"gosfV2/src/models"
	"net/http"

	"github.com/labstack/echo"
)

// Valida los datos de un usuario, si el usuario no existe o la contraseña es incorrecta
func UserCredencialMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(models.User)

		if err := c.Bind(user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if user.Username == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Username must not be empty")
		}

		if user.Password == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Password must not be empty")
		}

		dbUser, err := models.Users(c).FindUserByName(user.Username)
		if err != nil {
			if err == models.ErrUserNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "Invalid username or password")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		ok, err := CheckPassword(user.Password, dbUser.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
		}

		c.Set("user", dbUser)
		return next(c)
	}
}

// Verifica que el usuario tenga un token válido
// Ya sea en el Header, en el QueryParam o en el Cookie
func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token string

		// Verifica si la petición viene con el token
		token, err := GetToken(c)
		if err != nil {
			return err
		}

		claims, err := validJWT(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		// Set user in the echo context
		c.Set("user_id", claims.UserID)
		c.Set("claims", claims)
		c.Set("token", token)
		return next(c)
	}
}
