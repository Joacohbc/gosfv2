package env

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type config struct {
	SQL             sql
	Redis           redis
	JWTKey          string `json:"jwt_key"`
	JWTMinutes      int    `json:"jwt_minutes"`
	Port            int    `json:"port"`
	LogDirPath      string `json:"log_dir_path"`
	FilesDirectory  string `json:"files_directory"`
	StaticFiles     string `json:"static_files"`
	MaxTokenPerUser int    `json:"max_token_per_user"`
}

type sql struct {
	Host     string `json:"db_host_sql"`
	User     string `json:"db_user_sql"`
	Password string `json:"db_password_sql"`
	Name     string `json:"db_name_sql"`
	Port     int    `json:"db_port_sql"`
	Charset  string `json:"db_charset_sql"`
}

type redis struct {
	Host     string `json:"redis_host"`
	Port     int    `json:"redis_port"`
	Password string `json:"redis_password"`
	DB       int    `json:"redis_db"`
}

var Config config

func init() {

	configFile := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	b, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatal("Error reading config.json:", err)
	}

	if err := json.Unmarshal(b, &Config); err != nil {
		log.Fatal("Error unmarshaling config.json:", err)
	}

	if err := json.Unmarshal(b, &Config.SQL); err != nil {
		log.Fatal("Error unmarshaling config.json:", err)
	}

	if err := json.Unmarshal(b, &Config.Redis); err != nil {
		log.Fatal("Error unmarshaling config.json:", err)
	}

	if _, err := os.Stat(Config.StaticFiles); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Static files directory does not exist: " + Config.StaticFiles)
		}
	}
}
