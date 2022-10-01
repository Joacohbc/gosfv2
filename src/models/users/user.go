package users

import (
	"gosfV2/src/models/db"
	"log"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	conn := db.GetConn()

	rows, err := conn.Query(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY, 
			username VARCHAR(255) UNIQUE NOT NULL, 
			password BLOB NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Error to create table users:", err)
	}
	rows.Close()
}
