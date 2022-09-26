package auth

import (
	"gosfV2/src/models/env"
	"gosfV2/src/models/users"
	"strings"

	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Username string
}

func JWTMiddlewareConfigured() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(env.Config.JWTKey),
		Claims:     &UserClaims{},
	})
}

func RegisterHandler(c echo.Context) error {
	user := new(users.User)

	err := c.Bind(user)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" {
		return c.String(http.StatusBadRequest, "Username or password is empty")
	}

	exist, err := users.FindUserByName(user.Username)
	if err != nil {
		c.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if exist {
		return c.String(http.StatusUnauthorized, "A user with that username already exists")
	}

	err = users.NewUser(*user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "You have been registered successfully")
}

func LoginHandler(c echo.Context) error {
	user := new(users.User)

	err := c.Bind(user)
	if err != nil {
		c.Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	exist, err := users.FindUser(user.Username, user.Password)
	if err != nil {
		c.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if !exist {
		return c.String(http.StatusUnauthorized, "Invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(env.Config.JWTHours))),
		},
		Username: user.Username,
	})

	tokenString, err := token.SignedString([]byte(env.Config.JWTKey))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, tokenString)
}
