package main

import (
	"context"
	"gosfV2/src/auth"
	"gosfV2/src/middleware/logger"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	e := echo.New()
	e.Use(logger.MyLoggerConfig())
	e.Use(middleware.Recover())

	group := e.Group("/auth/api")
	group.Use(auth.JWTMiddlewareConfigured())
	{
		group.GET("/", func(ctx echo.Context) error {
			return ctx.String(http.StatusOK, "You logged sucesfully!")
		})
	}

	e.POST("/login", auth.LoginHandler)
	e.POST("/register", auth.RegisterHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		if err := e.Shutdown(context.Background()); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	e.Start(":8080")
}
