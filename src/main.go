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
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const Version = "v1.0.0"

func main() {

	// Si hubo un error en el proceso de inicialización, se detiene el programa
	if err := recover(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	// Configuración de los archivos estáticos
	front := e.Group("/*")
	front.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  env.Config.StaticFiles,
		HTML5: true,
	}))
	front.Use(middleware.Secure())
	front.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieDomain:   "example.com",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))

	// Configuración de los middlewares
	e.Use(logger.RequestLoggerConfig())
	e.Use(middleware.Recover())
	// e.Use(middleware.CORS())

	// Configuración de las rutas
	tokens := e.Group("/auth")
	routes.Auth.AddAuthRoutes(tokens)

	api := e.Group("/api", auth.Middlewares.JWTAuthMiddleware)
	routes.Files.AddRoutesToGroup(api)
	routes.User.AddRoutesToGroup(api)

	// Configuración del cierre del servidor
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		fmt.Println("Shutting down server... (this may take a few seconds)")
		if err := database.GetMySQL().Close(); err != nil {
			e.Logger.Error(err.Error())
		}

		if err := database.GetRedis().Close(); err != nil {
			e.Logger.Error(err.Error())
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err.Error())
		}
	}()

	// Configuración de la salida del servidor
	e.HideBanner = true
	e.HidePort = true
	e.Logger = logger.Logger(log.DEBUG)
	e.Debug = false

	// Configuración de la salida de la información del servidor
	staticAbs, err := filepath.Abs(env.Config.StaticFiles)
	if err != nil {
		log.Fatal(err.Error())
	}

	filesAbs, err := filepath.Abs(env.Config.FilesDirectory)
	if err != nil {
		log.Fatal(err.Error())
	}

	logAbs, err := filepath.Abs(env.Config.LogDirPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(`
____  ___  ____  _____       ____  
/ ___|/ _ \/ ___||  ___|_   _|___ \ 
| |  _| | | \___ \| |_  \ \ / / __) |
| |_| | |_| |___) |  _|  \ V / / __/ 
\____|\___/|____/|_|     \_/ |_____|
Powered By Echo v4 with Go Language - ` + Version)
	fmt.Println()
	fmt.Println("Server's Configuration: ")
	fmt.Println("- Running on port " + strconv.Itoa(env.Config.Port))
	// fmt.Println("- Frontend is running on port " + strconv.Itoa(env.Config.Port+1))
	fmt.Println("- Frontend files are served from " + staticAbs)
	fmt.Println("- Files are stored in " + filesAbs)
	fmt.Println("- Logs are stored in " + logAbs)
	fmt.Println()
	fmt.Println("MySQL's Connection: ")
	fmt.Println("- Host: " + env.Config.SQL.Host)
	fmt.Println("- Port: " + strconv.Itoa(env.Config.SQL.Port))
	fmt.Println("- User: " + env.Config.SQL.User)
	fmt.Println("- Database: " + env.Config.SQL.Name)
	fmt.Println()
	fmt.Println("Redis's Connection: ")
	fmt.Println("- Host: " + env.Config.Redis.Host)
	fmt.Println("- Port: " + strconv.Itoa(env.Config.Redis.Port))
	fmt.Println("- Database: " + strconv.Itoa(env.Config.Redis.DB))
	fmt.Println()
	fmt.Println("Starting server...")
	fmt.Println("Press CTRL+C to stop the server")

	// go func() {
	// 	front := echo.New()
	// 	front.HideBanner = true
	// 	front.HidePort = true
	// 	// front.Static("/*", env.Config.StaticFiles)
	// 	front.Use(middleware.StaticWithConfig(middleware.StaticConfig{
	// 		Root:  env.Config.StaticFiles,
	// 		HTML5: true,
	// 	}))

	// 	if err := front.Start(":" + strconv.Itoa(env.Config.Port+1)); err != nil {
	// 		e.Logger.Fatal(err.Error())
	// 	}
	// }()

	// Inicio del servidor
	if err := e.Start(":" + strconv.Itoa(env.Config.Port)); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
