package main

import (
	"context"
	"flag"
	"gosfV2/src/auth"
	"gosfV2/src/middleware/logger"
	"gosfV2/src/routes"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "8080", "Port to listen on")
	flag.Parse()
}

func main() {
	e := echo.New()

	e.Use(logger.RequestLoggerConfig())
	e.Use(middleware.Recover())
	e.Use(auth.JWTAuthMiddleware)
	e.Logger = logger.Logger(log.DEBUG)

	// Test Endpoint
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	routes.Auth.AddRoutes(e)

	api := e.Group("/api")
	routes.Files.AddRoutesToGroup(api)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	e.Start(":" + port)
}
