package auth

import (
	"errors"
	"gosfV2/src/models"
	"gosfV2/src/models/env"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
)

type UserClaims struct {
	Username string
	jwt.RegisteredClaims
}

func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Path() == "/register" || c.Path() == "/login" || c.Path() == "/logout" {
			return next(c)
		}

		var token string
		if strings.HasPrefix(c.Path(), "/api") {
			t, err := getTokenFromHeader(c)
			if err != nil {
				return errors.New(err.Error() + " (header)")
			}
			token = t

		} else if strings.HasPrefix(c.Path(), "/auth") {
			t, err := getTokenFromCookie(c)
			if err != nil {
				return errors.New(err.Error() + " (cookie)")
			}
			token = t
		}

		claims, err := ValidJWT(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		// Get user from the database
		dbUser, err := models.Users.ExistUserByName(c.Request().Context(), claims.Username)
		if err != nil {
			return err
		}

		if !dbUser {
			return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
		}

		// Set user in the echo context
		c.Set("username", claims.Username)
		c.Set("claims", claims)
		c.Set("token", token)

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

	exist, err := models.Users.ExistUserByName(c.Request().Context(), user.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if exist {
		return echo.NewHTTPError(http.StatusBadRequest, "Username already exists")
	}

	err = models.Users.NewUser(c.Request().Context(), *user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "You have been registered successfully")
}

func LoginHandler(c echo.Context) error {

	if ck, err := c.Cookie("token"); err == nil {

		if _, err := ValidJWT(ck.Value); err != nil {
			c.SetCookie(&http.Cookie{
				Name:   "token",
				MaxAge: -1, // Poniendo -1 se borra la cookie
			})
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Cookie: "+err.Error())
		}

		return c.String(http.StatusOK, ck.Value)
	}

	tokenString, err := generateJWTForUser(c)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * time.Duration(env.Config.JWTHours)),
		MaxAge:  3600 * env.Config.JWTHours, // El MaxAge se piede en segundos (3600s = 1 hora)
	})

	return c.String(http.StatusOK, tokenString)
}

func LogoutHandler(c echo.Context) error {

	if _, err := c.Cookie("token"); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not logged in")
	}

	c.SetCookie(&http.Cookie{
		Name:   "token",
		MaxAge: -1, // Poniendo -1 se borra la cookie
	})
	return c.String(http.StatusOK, "You have been logged out")
}
