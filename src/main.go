package main

import (
	"context"
	"flag"
	"gosfV2/src/auth"
	"gosfV2/src/middleware/logger"
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
	e.Logger = logger.Logger(log.INFO)

	// Test Endpoint
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// Auth Endpoints
	e.POST("/login", auth.LoginHandler)
	e.POST("/register", auth.RegisterHandler)
	e.GET("/logout", auth.LogoutHandler)

	// Authentificated Endpoints
	group := e.Group("/auth")
	{
		group.GET("/", func(ctx echo.Context) error {
			return ctx.String(http.StatusOK, "You logged sucesfully!")
		})
	}

	api := e.Group("/api")
	{
		api.GET("/", func(ctx echo.Context) error {
			return ctx.String(http.StatusOK, "You are authenticated!")
		})
	}

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
