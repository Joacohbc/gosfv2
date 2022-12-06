package auth

import (
	"context"
	"gosfV2/src/models"
	"net/http"
	"strings"
	"time"

	"gosfV2/src/utils"

	"github.com/labstack/echo"
)

// Maneja los errores de los archivos, si el error ErrUserNotFound
// o si es un error desconocido (base de datos), devuelve un error 500
func HandleUserError(err error) error {
	if err == models.ErrUserNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

func LoginValid(c echo.Context) (models.User, error) {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return models.User{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dbUser, err := models.Users(c).FindUserByName(user.Username)
	if err != nil {
		if err == models.ErrUserNotFound {
			return models.User{}, echo.NewHTTPError(http.StatusNotFound, "Invalid username or password")
		}
		return models.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ok, err := CheckPassword(user.Password, dbUser.Password)
	if err != nil {
		return models.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !ok {
		return models.User{}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	return dbUser, nil
}

func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Path() == "/register" || c.Path() == "/login" || c.Path() == "/tokens" || strings.HasPrefix(c.Path(), "/static") {
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
		} else if c.QueryParam(queryName) != "" {
			t, err := getTokenAsQueryParam(c)
			if err != nil {
				return err
			}
			token = t

		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Not token provided")
		}

		claims, err := ValidJWT(token)
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

	if err := GeneratePassword(&user.Password); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = models.Users(c).NewUser(*user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.ToJSON("User created successfully"))
}

func LoginHandler(c echo.Context) error {

	user, err := LoginValid(c)
	if err != nil {
		return err
	}

	tokenString, err := generateJWTForUser(user.ID, user.Username, c.RealIP())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})
}

func RefreshHandler(c echo.Context) error {
	claims, err := ValidJWT(c.Get("token").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if err := NewTokenManager().RemoveToken(c.Get("user_id").(uint), c.Get("token").(string)); err != nil {
		return HandlerTokenError(err)
	}

	tokenString, err := generateJWTForUser(claims.UserID, claims.Username, c.RealIP())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})
}

func DeleteTokens(c echo.Context) error {
	user, err := LoginValid(c)
	if err != nil {
		return err
	}

	if err := NewTokenManager().RemoveUserTokens(user.ID); err != nil {
		return HandlerTokenError(err)
	}

	return c.JSON(http.StatusOK, utils.ToJSON("All tokens have been deleted successfully"))
}

func VerifyAuth(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.ToJSON("You are authenticated"))
}

func LogoutHandler(c echo.Context) error {
	if err := NewTokenManager().RemoveToken(c.Get("user_id").(uint), c.Get("token").(string)); err != nil {
		return HandlerTokenError(err)
	}

	return c.JSON(http.StatusOK, utils.ToJSON("You have been logged out successfully"))
}

func init() {
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			users, err := models.UsersC(context.Background()).GetAllUsers()
			if err != nil && err != models.ErrUserNotFound {
				panic(err)
			}

			for _, user := range users {
				tokens, err := NewTokenManager().GetTokens(user.ID)
				if err != nil && err != ErrTokenNotFound {
					panic(err)
				}

				for _, token := range tokens {
					if _, err := ValidJWT(token); err != nil {
						if err := NewTokenManager().RemoveToken(user.ID, token); err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}()
}
