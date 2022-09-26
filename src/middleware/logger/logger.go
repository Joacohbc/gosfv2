package logger

import (
	"gosfV2/src/models/env"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func MyLoggerConfig() echo.MiddlewareFunc {

	err := os.MkdirAll(env.Config.LogDirPath, os.ModeDir)
	if err != nil {
		log.Fatal("Error creating log directory: ", err)
	}

	time := time.Now().Format("2006-01-02_15-04-05-999999999")

	file, err := os.Create(filepath.Join(env.Config.LogDirPath, time+".log"))
	if err != nil {
		log.Fatal("Error creating log file: ", err)
	}

	config := middleware.LoggerConfig{
		Format:           "${time_custom} => ${remote_ip} - ${method} ${uri} | ${status} - ${latency_human}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           io.MultiWriter(file, os.Stdout),
	}

	return middleware.LoggerWithConfig(config)
}
