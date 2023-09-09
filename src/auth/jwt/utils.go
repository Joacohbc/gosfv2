package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"gosfV2/src/models/env"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID   uint
	IP       string
	Location string
}

var tokenDuration time.Duration = time.Minute * time.Duration(env.Config.JWTMinutes)

func getExpirationDate() time.Time {
	return time.Now().Add(tokenDuration)
}

const (
	// Bearer es el tipo de Header de autenticación que se utiliza en JWT
	authScheme string = "Bearer"

	// api-token es el nombre del QueryParam donde se buscara el JWT
	queryName string = "api-token"

	// token es el nombre de la cookie donde se buscara el JWT
	cookieName string = "token"
)

// Obtiene el Token del usuario del Header
// El error que se genera es un error de tipo echo.HTTPError
func GetTokenFromHeader(c echo.Context) (string, error) {

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

// Obtiene el Token del usuario del QueryParam
// El error que se genera es un error de tipo echo.HTTPError
func GetTokenFromQueryParam(c echo.Context) (string, error) {
	token := c.QueryParam(queryName)
	if strings.TrimSpace(token) == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	}
	return token, nil
}

// Obtiene el Token del usuario del QueryParam, Header o Cookie
// El error que se genera es un error de tipo echo.HTTPError
func GetTokenFromRequest(c echo.Context) (string, error) {
	// Si tiene un Header Authorization, se toma el token de ahí
	if c.Request().Header.Get(echo.HeaderAuthorization) != "" {
		t, err := GetTokenFromHeader(c)
		if err != nil {
			return "", err
		}
		return t, nil

		// Si no tiene un Header Authorization, se busca el Token en el URL
	} else if c.QueryParam(queryName) != "" {
		t, err := GetTokenFromQueryParam(c)
		if err != nil {
			return "", err
		}
		return t, nil

		// Si no tiene un Header Authorization ni en el URL, se busca el Token en el Cookie
	} else if ck, err := GetTokenCookie(c); err == nil {
		return ck, nil

	} else {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Not token provided")
	}
}

// Obtiene la Localización del IP que se le pase
func GetLocation(ip string) string {
	res, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=country,regionName,city,status", ip))
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ""
	}

	var location struct {
		Country    string `json:"country"`
		RegionName string `json:"regionName"`
		City       string `json:"city"`
		Status     string `json:"status"`
	}

	if err := json.NewDecoder(res.Body).Decode(&location); err != nil {
		return ""
	}

	if location.Status != "success" {
		return ""
	}

	return fmt.Sprintf("%s, %s, %s", location.Country, location.RegionName, location.City)
}

// Genera un JWT para el usuario que se está logueado
func GenerateJWT(userId uint, ip string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserID:   userId,
		IP:       ip,
		Location: GetLocation(ip),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(getExpirationDate()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gosfV2",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString([]byte(env.Config.JWTKey))
	if err != nil {
		return "", err
	}

	if err := TokenManager.AddToken(userId, tokenString); err != nil {
		return "", err
	}

	return tokenString, nil
}

// Valida el JWT que se le pase y retorna los Claims JWT
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

	exist, err := TokenManager.ExistsToken(claims.UserID, tokenString)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("the token is not valid for the current session")
	}

	return claims, nil
}
