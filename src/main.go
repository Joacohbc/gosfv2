package main

import (
	"context"
	"fmt"
	"gosfV2/src/auth"
	"gosfV2/src/middleware/logger"
	"gosfV2/src/models/database"
	"gosfV2/src/models/env"
	"gosfV2/src/routes"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {

	if r := recover(); r != nil {
		fmt.Println(r)
		os.Exit(1)
	}
	e := echo.New()

	e.Static("/static", env.Config.StaticFiles)

	e.Use(logger.RequestLoggerConfig())
	e.Use(middleware.Recover())
	e.Logger = logger.Logger(log.DEBUG)

	// Test Endpoint
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	tokens := e.Group("/auth")
	routes.Auth.AddAuthRoutes(tokens)

	api := e.Group("/api", auth.JWTAuthMiddleware)
	routes.Files.AddRoutesToGroup(api)
	routes.User.AddRoutesToGroup(api)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		if err := database.GetMySQL().Close(); err != nil {
			e.Logger.Error(err.Error())
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err.Error())
		}
	}()

	if err := e.Start(":" + strconv.Itoa(env.Config.Port)); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
