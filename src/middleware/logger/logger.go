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
	echoLog "github.com/labstack/gommon/log"
)

var LogWriter io.Writer

func init() {
	err := os.MkdirAll(env.Config.LogDirPath, 0744)
	if err != nil {
		log.Fatal("Error creating log directory: ", err)
	}

	time := time.Now().Format("2006-01-02_15-04-05.00000")

	file, err := os.Create(filepath.Join(env.Config.LogDirPath, time+".log"))
	if err != nil {
		log.Fatal("Error creating log file: ", err)
	}

	LogWriter = io.MultiWriter(file, os.Stdout)
}

func RequestLoggerConfig() echo.MiddlewareFunc {
	config := middleware.LoggerConfig{
		Format:           "${time_custom} => ${remote_ip} - ${method} ${uri} | ${status} - ${latency_human}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           LogWriter,
	}
	return middleware.LoggerWithConfig(config)
}

func Logger(lvl echoLog.Lvl) echo.Logger {
	logger := echo.New().Logger
	logger.SetOutput(LogWriter)
	logger.SetLevel(lvl)
	if lvl == echoLog.DEBUG {
		logger.SetHeader("LOG > ${level} ${time_rfc3339_nano} FROM LINE ${line} IN ${short_file} >")
	} else {
		logger.SetHeader("${level} | ${time_rfc3339_nano} >")
	}
	return logger
}
