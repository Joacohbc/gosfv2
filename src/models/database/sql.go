package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gosfV2/src/models/env"
	"os"
	"strconv"
	"time"

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

	for i := 1; i <= 10; i++ {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", env.Config.DBUserSQL, env.Config.DBPasswordSQL, env.Config.DBHostSQL, env.Config.BDPortSQL, env.Config.DBNameSQL, env.Config.BDCharsetSQL)

		var err error
		temp, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			fmt.Printf("Error to connect to MySQL (%s): %s\n", dsn, err.Error())
			fmt.Printf("Retrying the MySQL connection in 10 seconds (time %d)..\n", i)
			time.Sleep(10 * time.Second)
			continue
		}

		temp.SetConnMaxLifetime(0)
		temp.SetMaxIdleConns(50)
		temp.SetMaxOpenConns(50)

		mySqlDB = temp
		break
	}

	if mySqlDB == nil {
		fmt.Println("Error to connect to MySQL")
		os.Exit(1)
	}

	for i := 1; i <= 10; i++ {
		temp := redis.NewClient(&redis.Options{
			Addr:     env.Config.RedisHost + ":" + strconv.Itoa(env.Config.RedisPort),
			Password: env.Config.RedisPassword,
			DB:       env.Config.RedisDB,
		})

		if err := temp.Ping(context.Background()).Err(); err != nil {
			fmt.Printf("Error to connect to Redis (%s - %s - %d): %s\n", redisDb.Options().Addr, redisDb.Options().Password, redisDb.Options().DB, err.Error())
			fmt.Printf("Retrying the MySQL connection in 10 seconds (time %d)..\n", i)
			time.Sleep(10 * time.Second)
			continue
		}

		redisDb = temp
		break
	}

	if redisDb == nil {
		fmt.Println("Error to connect to Redis")
		os.Exit(1)
	}

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

}
