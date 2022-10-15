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

const (

	// Bearer es el tipo de Header de autenticación que se utiliza en JWT
	authScheme string = "Bearer"

	// token es el nombre de la Cookie donde se buscara el JWT
	cookieName string = "token"
)

func getTokenFromHeader(c echo.Context) (string, error) {

	// Obtengo las 2 partes del Header, el tipo[0] y el contenido[1]
	auth := strings.Fields(c.Request().Header.Get(echo.HeaderAuthorization))
	if len(auth) != 2 {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header")
	}

	// Si el tipo de Heaeder de autorización es diferente a Bearer, se devuelve un error
	if auth[0] != authScheme {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header, must be a "+authScheme+" token")
	}

	// Si el Token esta vacio
	if auth[1] == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	}

	return auth[1], nil
}

func getTokenFromCookie(c echo.Context) (string, error) {
	ck, err := c.Cookie(cookieName)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization cookie, must be a "+cookieName+" cookie")
	}

	if ck.Value == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	}

	return ck.Value, nil
}

// Genera un JWT para el usuario que se está logueando.
// El error que se genera es un error de tipo echo.HTTPError
func generateJWTForUser(c echo.Context) (string, error) {

	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	exist, err := models.Users.ExistUser(c.Request().Context(), user.Username, user.Password)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !exist {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(env.Config.JWTHours))),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gosfV2",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString([]byte(env.Config.JWTKey))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

	return claims, nil
}
