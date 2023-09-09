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
	// Crea el directorio de logs si no existe
	if err := os.MkdirAll(env.Config.LogDirPath, 0744); err != nil {
		log.Fatal("Error creating log directory: ", err)
	}

	// Obtiene el tiempo actual
	time := time.Now().Format("2006-01-02_15-04-05.00000")

	// Crea el archivo de log
	file, err := os.Create(filepath.Join(env.Config.LogDirPath, time+".log"))
	if err != nil {
		log.Fatal("Error creating log file: ", err)
	}

	// Asigna el archivo de log como writer
	LogWriter = io.MultiWriter(file, os.Stdout)
}

// RequestLoggerConfig configura el logger de las peticiones
func RequestLoggerConfig() echo.MiddlewareFunc {
	config := middleware.LoggerConfig{
		Format:           "${time_custom} => ${remote_ip} - ${method} ${uri} | ${status} - ${latency_human}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           LogWriter,
	}
	return middleware.LoggerWithConfig(config)
}

// Logger crea un logger para el servidor
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
