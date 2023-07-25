package utils

import (
	"github.com/labstack/echo/v4"
)

// ToJSON convierte un mensaje a un objeto JSON
func ToJSON(msg string) echo.Map {
	return echo.Map{
		"message": msg,
	}
}
