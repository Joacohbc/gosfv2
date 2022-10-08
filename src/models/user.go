package models

import (
	"context"
	"database/sql"
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

var Users users

type users struct{}

func (u users) GetAllUsers(ctx context.Context) ([]User, error) {
	conn := db.GetConn()

	rows, err := conn.QueryContext(ctx, `SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u users) NewUser(ctx context.Context, user User) error {
	conn := db.GetConn()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO users (username, password) VALUES (?, aes_encrypt(?,?))`, user.Username, user.Password, user.Password)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u users) FindUser(ctx context.Context, username string, password string) (bool, error) {
	conn := db.GetConn()

	row, err := conn.QueryContext(ctx, `SELECT id FROM users WHERE username = ? AND password = aes_encrypt(?, ?)`, username, password, password)
	if db.IsSQLError(err) {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return row.Next(), nil
}

func (u users) FindUserByName(ctx context.Context, username string) (bool, error) {
	conn := db.GetConn()

	row, err := conn.QueryContext(ctx, `SELECT id FROM users WHERE username = ?`, username)
	if db.IsSQLError(err) {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return row.Next(), nil
}

func (u users) GetIdByName(ctx context.Context, username string) (int, error) {
	conn := db.GetConn()

	var id int
	err := conn.QueryRowContext(ctx, `SELECT id FROM users WHERE username = ?`, username).Scan(&id)
	if db.IsSQLError(err) {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return id, nil
}
