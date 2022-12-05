package database

import (
	"database/sql"
	"errors"
	"fmt"
	"gosfV2/src/models/env"
	"os"

	redis "github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	mySqlDB *sqlx.DB
	redisDb *redis.Client
)

func GetMySQL() *sqlx.DB {
	return mySqlDB
}

func GetRedis() *redis.Client {
	return redisDb
}

// Retorna true si es sql.ErrNoRows)
func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)

}

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", env.Config.DBUser, env.Config.DBPassword, env.Config.DBHost, env.Config.BDPort, env.Config.DBName, env.Config.BDCharset)

	var err error
	mySqlDB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("Error to connect to database:", err)
		os.Exit(1)
	}

	mySqlDB.SetConnMaxLifetime(0)
	mySqlDB.SetMaxIdleConns(50)
	mySqlDB.SetMaxOpenConns(50)

	mySqlDB.MustExec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)

	mySqlDB.MustExec(`
	CREATE TABLE IF NOT EXISTS files (
		file_id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		filename VARCHAR(255) NOT NULL UNIQUE,
		shared BOOLEAN NOT NULL DEFAULT false,
		update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		user_id INT UNSIGNED NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)

	mySqlDB.MustExec(`
	CREATE TABLE IF NOT EXISTS file_users (
		user_id INT UNSIGNED NOT NULL,
		file_id INT UNSIGNED NOT NULL,
		PRIMARY KEY(user_id, file_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (file_id) REFERENCES files(file_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)

	redisDb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
