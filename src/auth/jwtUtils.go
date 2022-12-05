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
	"golang.org/x/crypto/bcrypt"
)

const (
	// Bearer es el tipo de Header de autenticación que se utiliza en JWT
	authScheme string = "Bearer"

	// api-token es el nombre del QueryParam donde se buscara el JWT
	queryName string = "api-token"
)

type UserClaims struct {
	UserID   uint
	Username string
	jwt.RegisteredClaims
}

// Inscrita el password con AES y retorna la cadena encriptada
func GeneratePassword(password *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*password = string(hash)
	return nil
}

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

func getTokenFromHeader(c echo.Context) (string, error) {

	// Obtengo las 2 partes del Header, el tipo[0] y el contenido[1]
	auth := strings.Fields(c.Request().Header.Get(echo.HeaderAuthorization))
	if len(auth) != 2 {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header")
	}

	// Si el tipo de Header de autorización es diferente a Bearer, se devuelve un error
	if auth[0] != authScheme {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header, must be a "+authScheme+" token")
	}

	// Si el Token esta vació
	if auth[1] == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	}

	return auth[1], nil
}

func getTokenAsQueryParam(c echo.Context) (string, error) {
	token := c.Param(queryName)
	if strings.TrimSpace(token) == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	}
	return token, nil
}

// Genera un JWT para el usuario que se está logueado.
// El error que se genera es un error de tipo echo.HTTPError
func generateJWTForUser(c echo.Context) (string, error) {

	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dbUser, err := models.Users(c).FindUserByName(user.Username)
	if err != nil {
		if err == models.ErrUserNotFound {
			return "", echo.NewHTTPError(http.StatusNotFound, "Invalid username or password")
		}
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ok, err := CheckPassword(user.Password, dbUser.Password)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserID:   dbUser.ID,
		Username: dbUser.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gosfV2",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString([]byte(env.Config.JWTKey))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := NewTokenManager().AddToken(dbUser.ID, tokenString); err != nil {
		return "", HandlerTokenError(err)
	}

	return tokenString, nil
}

func ValidJWT(tokenString string) (*UserClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(env.Config.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*UserClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, jwt.ErrTokenExpired
	}

	exist, err := NewTokenManager().ExistsToken(claims.UserID, tokenString)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("the token is not valid for the current session")
	}

	return claims, nil
}
