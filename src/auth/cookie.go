package auth

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

const (
	cookieQueryName = "cookie"
)

// Obtiene el Valor de la Cookie de la petición
func GetTokenCookie(c echo.Context) (string, error) {
	ck, err := c.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return ck.Value, nil
}

// Setea la Cookie en la respuesta (si se pasara el QueryParam de la Cookie)
func SetTokenCookie(c echo.Context, token string) {
	if c.QueryParam(cookieQueryName) != "" {
		c.SetCookie(&http.Cookie{
			Name:    cookieName,
			Value:   token,
			Expires: time.Now().Add(TokenDuration),
			MaxAge:  int(TokenDuration.Seconds()),
			Path:    "/",
		})
	}
}

// Borra la Cookie en la respuesta (si se pasara el QueryParam de la Cookie)
func DeleteTokenCookie(c echo.Context, token string) {
	if c.QueryParam(cookieQueryName) != "" {
		c.SetCookie(&http.Cookie{
			Name:    cookieName,
			MaxAge:  -1,
			Expires: time.Now().Add(-1 * time.Minute),
		})
	}
}
