package jwt

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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
func SetTokenCookieFromQuery(c echo.Context, token string) {
	if c.QueryParam(cookieQueryName) != "" {
		c.SetCookie(&http.Cookie{
			Name:     cookieName,
			Value:    token,
			Expires:  time.Now().Add(tokenDuration),
			MaxAge:   int(tokenDuration.Seconds()),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
	}
}

// Borra la Cookie en la respuesta (si se pasara el QueryParam de la Cookie)
func DeleteTokenCookieFromQuery(c echo.Context, token string) {
	if c.QueryParam(cookieQueryName) != "" {
		c.SetCookie(&http.Cookie{
			Name:    cookieName,
			MaxAge:  -1,
			Expires: time.Now().Add(-1 * time.Minute),
		})
	}
}
