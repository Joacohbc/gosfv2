package env

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	JWTKey     string `json:"jwt_key"`
	JWTHours   int    `json:"jwt_hours"`
	LogDirPath string `json:"log_dir_path"`
	DBHost     string `json:"db_host"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
}

var Config config

func init() {
	b, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error reading config.json:", err)
	}

	err = json.Unmarshal(b, &Config)
	if err != nil {
		log.Fatal("Error unmarshaling config.json:", err)
	}
}
