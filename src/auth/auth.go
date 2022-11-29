package auth

import (
	"gosfV2/src/models"
	"gosfV2/src/models/env"
	"net/http"
	"strings"
	"time"

	"gosfV2/src/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
)

type UserClaims struct {
	UserID   uint
	Username string
	jwt.RegisteredClaims
}

// Maneja los errores de los archivos, si el error ErrUserNotFound
// o si es un error desconocido (base de datos), devuelve un error 500
func HandleUserError(err error) error {
	if err == models.ErrUserNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Path() == "/register" || c.Path() == "/login" || c.Path() == "/logout" || strings.HasPrefix(c.Path(), "/static") {
			return next(c)
		}

		var token string

		// Si tiene un Header Authorization, se toma el token de ahí
		if c.Request().Header.Get(echo.HeaderAuthorization) != "" {
			t, err := getTokenFromHeader(c)
			if err != nil {
				return err
			}
			token = t

			// Si no tiene un Header Authorization, se busca el Token en el URL
		} else if c.Param(queryName) != "" {
			t, err := getTokenAsQueryParam(c)
			if err != nil {
				return err
			}
			token = t

			// Si no tiene el URL, se busca el Token en una Cookie
		} else if _, err := c.Cookie(cookieName); err == nil {
			t, err := getTokenFromCookie(c)
			if err != nil {
				return err
			}
			token = t

			// Si no tiene Token lo informo
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Not token provided")
		}

		claims, err := ValidJWT(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		// Set user in the echo context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("claims", claims)
		c.Set(cookieName, token)

		return next(c)
	}
}

func RegisterHandler(c echo.Context) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username or password is empty")
	}

	exist, err := models.Users(c).ExistUserByName(user.Username)
	if err != nil {
		return HandleUserError(err)
	}

	if exist {
		return echo.NewHTTPError(http.StatusBadRequest, "Username already exists")
	}

	if err := generatePassword(&user.Password); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = models.Users(c).NewUser(*user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.ToJSON("User created successfully"))
}

func LoginHandler(c echo.Context) error {
	if ck, err := c.Cookie(cookieName); err == nil {

		if _, err := ValidJWT(ck.Value); err != nil {
			c.SetCookie(&http.Cookie{
				Name:   cookieName,
				MaxAge: -1, // Poniendo -1 se borra la cookie
			})
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Cookie: "+err.Error())
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": ck.Value,
		})
	}

	tokenString, err := generateJWTForUser(c)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:    cookieName,
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * time.Duration(env.Config.JWTHours)),
		MaxAge:  3600 * env.Config.JWTHours, // El MaxAge se pide en segundos (3600s = 1 hora)
	})

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})
}

func LogoutHandler(c echo.Context) error {

	if _, err := c.Cookie(cookieName); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not logged in")
	}

	c.SetCookie(&http.Cookie{
		Name:   cookieName,
		MaxAge: -1, // Poniendo -1 se borra la cookie
	})

	return c.JSON(http.StatusOK, utils.ToJSON("You have been logged out successfully"))
}
