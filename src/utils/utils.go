package utils

import "github.com/labstack/echo"

// ToJSON convierte un mensaje a un objeto JSON
func ToJSON(msg string) echo.Map {
	return echo.Map{
		"message": msg,
	}
}
