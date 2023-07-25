package auth

import (
	"context"
	"errors"
	"gosfV2/src/auth/jwt"
	"gosfV2/src/models"
	"gosfV2/src/models/env"
	"net/http"
	"strings"
	"time"

	"gosfV2/src/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Compara el password con el hash de la base de datos
func CheckPassword(password, bdHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(bdHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Registra un nuevo usuario
func RegisterUser(c echo.Context) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

// Genera un token JWT para el usuario si el usuario y la contraseña son correctos
func Login(c echo.Context) error {
	user := Middlewares.GetUser(c)

	tokenString, err := jwt.GenerateJWT(user.ID, c.RealIP())
	if err != nil {
		return jwt.HandlerTokenError(err)
	}

	// Si se paso el Cookie en el QueryParam se guarda en el Cookie
	jwt.SetTokenCookieFromQuery(c, tokenString)

	return c.JSON(http.StatusOK, echo.Map{
		"token":    tokenString,
		"duration": env.Config.JWTMinutes,
	})
}

// Refresca el token del usuario y lo guarda en la base de datos
// y en el cookie del usuario (si se paso el cookie en el QueryParam)
func RefreshToken(c echo.Context) error {
	if err := jwt.TokenManager.RemoveToken(Middlewares.GetUserId(c), Middlewares.GetUserToken(c)); err != nil {
		c.Logger().Error(err)
		return jwt.HandlerTokenError(err)
	}

	claims := Middlewares.GetUserClaims(c)
	tokenString, err := jwt.GenerateJWT(claims.UserID, c.RealIP())
	if err != nil {
		c.Logger().Error(err)
		return jwt.HandlerTokenError(err)
	}

	jwt.SetTokenCookieFromQuery(c, tokenString)

	return c.JSON(http.StatusOK, echo.Map{
		"token":    tokenString,
		"duration": env.Config.JWTMinutes,
	})
}

// Elimina todos los tokens del usuario
// y en el cookie del usuario (si se paso el cookie en el QueryParam)
func DeleteAllTokens(c echo.Context) error {
	user := Middlewares.GetUser(c)

	if err := jwt.TokenManager.RemoveUserTokens(user.ID); err != nil {
		return jwt.HandlerTokenError(err)
	}

	return c.JSON(http.StatusOK, utils.ToJSON("All tokens have been deleted successfully"))
}

// Elimina el token de la base de datos y del Cookie (si se paso el cookie en el QueryParam)
func Logout(c echo.Context) error {
	if err := jwt.TokenManager.RemoveToken(Middlewares.GetUserId(c), Middlewares.GetUserToken(c)); err != nil {
		return jwt.HandlerTokenError(err)
	}

	jwt.DeleteTokenCookieFromQuery(c, Middlewares.GetUserToken(c))

	return c.JSON(http.StatusOK, utils.ToJSON("You have been logged out successfully"))
}

// Si se puede acceder a la ruta, se retorna un mensaje de que se esta autenticado
func VerifyAuth(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.ToJSON("You are authenticated"))
}

func init() {
	// Cada un minutos se verifica si los tokens son validos
	// Si no son validos se eliminan
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			users, err := models.UsersC(context.Background()).GetAllUsers()
			if err != nil && err != models.ErrUserNotFound {
				panic(err)
			}

			for _, user := range users {
				tokens, err := jwt.TokenManager.GetTokens(user.ID)
				if err != nil && err != jwt.ErrTokenNotFound {
					panic(err)
				}

				for _, token := range tokens {
					if _, err := jwt.ValidJWT(token); err != nil {
						if err := jwt.TokenManager.RemoveToken(user.ID, token); err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}()
}
