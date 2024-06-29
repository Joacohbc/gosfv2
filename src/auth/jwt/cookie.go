package jwt

import (
	"gosfV2/src/models/env"
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

		config := &http.Cookie{
			Name:     cookieName,
			Value:    token,
			Expires:  time.Now().Add(tokenDuration),
			MaxAge:   int(tokenDuration.Seconds()),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}

		// Si estamos en modo desarrollo, no se establece la cookie como segura y se establece el SameSite como None
		// para poder probar la aplicación en local sin problemas
		if env.Config.DevMode {
			config.Secure = false
			config.SameSite = http.SameSiteDefaultMode
		}

		c.SetCookie(config)
	}
}

// Borra la Cookie en la respuesta (si se pasara el QueryParam de la Cookie)
func DeleteTokenCookieFromQuery(c echo.Context, token string) {
	if c.QueryParam(cookieQueryName) != "" {
		c.SetCookie(&http.Cookie{
			Name:     cookieName,
			MaxAge:   -1,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})
	}
}
