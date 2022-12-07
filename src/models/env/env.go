package env

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	JWTKey          string `json:"jwt_key"`
	JWTMinutes      int    `json:"jwt_minutes"`
	LogDirPath      string `json:"log_dir_path"`
	DBHostSQL       string `json:"db_host_sql"`
	DBUserSQL       string `json:"db_user_sql"`
	DBPasswordSQL   string `json:"db_password_sql"`
	DBNameSQL       string `json:"db_name_sql"`
	BDPortSQL       int    `json:"db_port_sql"`
	BDCharsetSQL    string `json:"db_charset_sql"`
	RedisHost       string `json:"redis_host"`
	RedisPort       int    `json:"redis_port"`
	RedisPassword   string `json:"redis_password"`
	RedisDB         int    `json:"redis_db"`
	FilesDirectory  string `json:"files_directory"`
	StaticFiles     string `json:"static_files"`
	MaxTokenPerUser int    `json:"max_token_per_user"`
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
