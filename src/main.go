package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"gosfV2/src/auth"
	"gosfV2/src/middleware/logger"
	"gosfV2/src/models/database"
	"gosfV2/src/models/env"
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
	flag.StringVar(&port, "port", "3000", "Port to listen on")
	flag.Parse()
}

type User struct {
	ID        int          `db:"id"`
	Username  string       `db:"username"`
	Password  string       `db:"password"`
	UpdateAt  sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	CreateAt  sql.NullTime `db:"created_at"`
}

func main() {

	if r := recover(); r != nil {
		fmt.Println(r)
		os.Exit(1)
	}

	e := echo.New()

	e.Static("/static", env.Config.StaticFiles)

	e.Use(logger.RequestLoggerConfig())
	e.Use(middleware.Recover())
	e.Use(auth.JWTAuthMiddleware)
	e.Logger = logger.Logger(log.DEBUG)

	// Test Endpoint
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := e.Group("/api")
	routes.Auth.AddNoAuthRoutes(e)
	routes.Files.AddRoutesToGroup(api)
	routes.User.AddRoutesToGroup(api)
	routes.Auth.AddTokenRoutes(api)

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

	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
