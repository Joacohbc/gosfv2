package handlers

import (
	"fmt"
	"gosfV2/src/dtos"
	"gosfV2/src/ent"
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
	fmt.Println(reflect.TypeOf(model))
	if reflect.TypeOf(model) == reflect.TypeOf(&ent.File{}) {
		return c.JSON(code, dtos.ToFileDTO(model.(*ent.File)))
	}

	if reflect.TypeOf(model) == reflect.TypeOf([]*ent.File{}) {
		return c.JSON(code, dtos.ToFileListDTO(model.([]*ent.File)))
	}

	if reflect.TypeOf(model) == reflect.TypeOf(&ent.User{}) {
		return c.JSON(code, dtos.ToUserDTO(model.(*ent.User)))
	}

	if reflect.TypeOf(model) == reflect.TypeOf([]ent.User{}) {
		return c.JSON(code, dtos.ToUserListDTO(model.([]*ent.User)))
	}

	return c.JSON(code, model)
}
