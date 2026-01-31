package utils

import (
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// ToJSON convierte un mensaje a un objeto JSON
func ToJSON(msg string) echo.Map {
	return echo.Map{
		"message": msg,
	}
}

var dangerousExts = map[string]bool{
	".html": true,
	".htm":  true,
	".svg":  true,
	".xml":  true,
	".js":   true,
}

// ServeFileSafe sirve un archivo de manera segura, forzando la descarga para tipos peligrosos
// para prevenir Stored XSS.
func ServeFileSafe(c echo.Context, path string, filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	if dangerousExts[ext] {
		return c.Attachment(path, filename)
	}

	return c.Inline(path, filename)
}
