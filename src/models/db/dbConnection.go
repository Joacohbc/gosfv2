package db

import (
	"context"
	"errors"
	"fmt"
	"gosfV2/src/models/env"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Dejar la conexion en Singleton
var db *gorm.DB

func GetBd() *gorm.DB {
	return db
}

func GetBdCtx(cxt context.Context) *gorm.DB {
	return db.WithContext(cxt)
}

// Retorna true si es gorm.ErrRecordNotFound)
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)

}

func init() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", env.Config.DBUser, env.Config.DBPassword, env.Config.DBHost, env.Config.BDPort, env.Config.DBName, env.Config.BDCharset)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// var err error
	// db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error to connect to database:", err)
	}
}
