package database

import (
	"database/sql"
	"errors"
	"fmt"
	"gosfV2/src/models/env"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func GetBd() *sqlx.DB {
	return db
}

// Retorna true si es sql.ErrNoRows)
func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)

}

func init() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", env.Config.DBUser, env.Config.DBPassword, env.Config.DBHost, env.Config.BDPort, env.Config.DBName, env.Config.BDCharset)

	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("Error to connect to database:", err)
		os.Exit(1)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
}
