package auth

import (
	"gosfV2/src/auth/jwt"
	"gosfV2/src/models"
	"net/http"

	"github.com/labstack/echo"
)

const (
	tokenContextKey  = "user_token"
	userIdContextKey = "user_id"
	userContextKey   = "user"
	claimsContextKey = "claims"
)

type authMiddleware struct {
}

var Middlewares authMiddleware

// Obtiene el ID del usuario logueado del contexto
func (a authMiddleware) GetUserId(c echo.Context) uint {
	return c.Get(userIdContextKey).(uint)
}

// Obtiene el token del usuario logueado del contexto
func (a authMiddleware) GetUserToken(c echo.Context) string {
	return c.Get(tokenContextKey).(string)
}

// Obtiene los claims del usuario logueado del contexto
func (a authMiddleware) GetUserClaims(c echo.Context) *jwt.UserClaims {
	return c.Get(claimsContextKey).(*jwt.UserClaims)
}

// Obtiene el usuario logueado del contexto
func (a authMiddleware) GetUser(c echo.Context) models.User {
	return c.Get(userContextKey).(models.User)
}

// Verifica que el usuario tenga un token válido
// Ya sea en el Header, en el QueryParam o en el Cookie
// Y si todo es correcto, se agrega el token al contexto, el ID del usuario y los claims
func (a authMiddleware) JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token string

		// Verifica si la petición viene con el token
		token, err := jwt.GetTokenFromRequest(c)
		if err != nil {
			return err
		}

		claims, err := jwt.ValidJWT(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		// Set user in the echo context
		c.Set(userIdContextKey, claims.UserID)
		c.Set(claimsContextKey, claims)
		c.Set(tokenContextKey, token)
		return next(c)
	}
}

// Valida los datos de un usuario, si el usuario no existe o la contraseña es incorrecta
// Y si todo es correcto, se agrega el usuario al contexto
func (a authMiddleware) UserCredencialMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

		c.Set(userContextKey, dbUser)
		return next(c)
	}
}
