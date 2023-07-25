package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Obtiene el ID como PathParam y lo convierte en uint
// Retorna un error 400 (Bad Request) si el ID no es un número
func getIdFromURL(c echo.Context, param string) (uint, error) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Invalid Id from URL")
	}

	return uint(id), nil
}
