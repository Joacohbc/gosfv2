package utils

import "github.com/labstack/echo"

func ToJSON(msg string) echo.Map {
	return echo.Map{
		"message": msg,
	}
}
