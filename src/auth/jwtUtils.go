package auth

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
	"golang.org/x/crypto/bcrypt"
)

const (
	// Bearer es el tipo de Header de autenticación que se utiliza en JWT
	authScheme string = "Bearer"

	// api-token es el nombre del QueryParam donde se buscara el JWT
	queryName string = "api-token"

	cookieName string = "token"
)

type UserClaims struct {
	UserID   uint
	Username string
	IP       string
	Location string
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

// Obtiene el Token del usuario del Header
// El error que se genera es un error de tipo echo.HTTPError
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

// Obtiene el Token del usuario del QueryParam
// El error que se genera es un error de tipo echo.HTTPError
func getTokenAsQueryParam(c echo.Context) (string, error) {
	token := c.QueryParam(queryName)
	if strings.TrimSpace(token) == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	}
	return token, nil
}

// Obtiene la Localización del IP que se le pase
func getLocation(ip string) string {
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
func generateJWTForUser(userId uint, username string, ip string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserID:   userId,
		Username: username,
		IP:       ip,
		Location: getLocation(ip),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gosfV2",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString([]byte(env.Config.JWTKey))
	if err != nil {
		return "", err
	}

	if err := NewTokenManager().AddToken(userId, tokenString); err != nil {
		return "", err
	}

	return tokenString, nil
}

// Valida el JWT que se le pase
func validJWT(tokenString string) (*UserClaims, error) {
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
