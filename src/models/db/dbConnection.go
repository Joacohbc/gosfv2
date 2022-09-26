package db

import (
	"context"
	"database/sql"
	"fmt"
	"gosfV2/src/models/env"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Dejar la conexion en Singleton
var poolConn *sql.DB

func GetConn() *sql.DB {
	return poolConn
}

// Retorna false si el err no es nil o sql.ErrNoRows
func NotErrNoRows(err error) bool {
	if err == sql.ErrNoRows {
		return false
	}

	if err == nil {
		return true
	}

	return false
}

func Ping(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := poolConn.PingContext(ctx); err != nil {
		return false
	}

	return true
}

func init() {
	var err error
	poolConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", env.Config.DBUser, env.Config.DBPassword, env.Config.DBHost, env.Config.DBName))
	if err != nil {
		log.Fatal("Error to connect to database:", err)
	}

	poolConn.SetConnMaxLifetime(time.Minute * 3)
	poolConn.SetMaxOpenConns(10)
	poolConn.SetMaxIdleConns(10)

	if err := poolConn.Ping(); err != nil {
		log.Fatal("Error to connect to database:", err)
	}
}
