package main

import (
	"context"
	"flag"
	"gosfV2/src/middleware/auth"
	"gosfV2/src/middleware/logger"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	e.Debug = false

	e.Use(logger.MyLoggerConfig())
	e.Use(middleware.Recover())

	// Auth Endpoints
	e.POST("/login", auth.LoginHandler)
	e.POST("/register", auth.RegisterHandler)

	// Authentificated Endpoints
	group := e.Group("/auth/api", auth.JWTMiddlewareConfigured())
	{
		group.GET("/", func(ctx echo.Context) error {
			return ctx.String(http.StatusOK, "You logged sucesfully!")
		})
	}

	// Test Endpoint
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

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
