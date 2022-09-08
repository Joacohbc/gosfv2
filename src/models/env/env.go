package env

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	JWTKey        string `json:"jwt_key"`
	JWTHours      int    `json:"jwt_hours"`
	UsersFilePath string `json:"users_file_path"`
	LogDirPath    string `json:"log_dir_path"`
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
