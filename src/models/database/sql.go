package database

import (
	"context"
	"fmt"
	"gosfV2/src/ent"
	"gosfV2/src/models/env"
	"log"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"
)

var (
	mySqlDB *ent.Client
	redisDb *redis.Client
)

func GetMySQL() *ent.Client {
	return mySqlDB
}

func GetRedis() *redis.Client {
	return redisDb
}

func init() {

	// Se conecta a MySQL, 10 intentos
	for i := 1; i <= 10; i++ {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", env.Config.SQL.User, env.Config.SQL.Password, env.Config.SQL.Host, env.Config.SQL.Port, env.Config.SQL.Name, env.Config.SQL.Charset)

		client, err := ent.Open("mysql", dsn)
		if err != nil {
			log.Printf("failed opening connection to mysql: %v", err)
			continue
		}

		// Run the auto migration tool.
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Printf("failed creating schema resources: %v", err)
			continue
		}

		mySqlDB = client
		break
	}

	if mySqlDB == nil {
		log.Fatal("Error to connect to MySQL")
	}

	// Se conecta a Redis, 10 intentos
	for i := 1; i <= 10; i++ {
		temp := redis.NewClient(&redis.Options{
			Addr:     env.Config.Redis.Host + ":" + strconv.Itoa(env.Config.Redis.Port),
			Password: env.Config.Redis.Password,
			DB:       env.Config.Redis.DB,
		})

		if err := temp.Ping(context.Background()).Err(); err != nil {
			log.Printf("Error to connect to Redis (%s:%d - %s - %d): %s\n", env.Config.Redis.Host, env.Config.Redis.Port, env.Config.Redis.Password, env.Config.Redis.DB, err.Error())
			log.Printf("Retrying the MySQL connection in 20 seconds (time %d)..\n", i)
			time.Sleep(20 * time.Second)
			continue
		}

		redisDb = temp
		break
	}

	if redisDb == nil {
		log.Fatal("Error to connect to Redis")
	}
}
