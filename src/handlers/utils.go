package handlers

import (
	"gosfV2/src/dtos"
	"gosfV2/src/models"
	"net/http"
	"reflect"
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

func jsonDTO(c echo.Context, code int, model interface{}) error {
	if reflect.TypeOf(model) == reflect.TypeOf(models.File{}) {
		return c.JSON(code, dtos.ToFileDTO(model.(models.File)))
	}

	if reflect.TypeOf(model) == reflect.TypeOf([]models.File{}) {
		return c.JSON(code, dtos.ToFileListDTO(model.([]models.File)))
	}

	if reflect.TypeOf(model) == reflect.TypeOf(models.User{}) {
		return c.JSON(code, dtos.ToUserDTO(model.(models.User)))
	}

	if reflect.TypeOf(model) == reflect.TypeOf([]models.User{}) {
		return c.JSON(code, dtos.ToUserListDTO(model.([]models.User)))
	}

	return c.JSON(code, model)
}
