package users

import (
	"context"
	"database/sql"
	"gosfV2/src/models/db"
)

func GetAllUsers(ctx context.Context) ([]User, error) {
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

func NewUser(ctx context.Context, user User) error {
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

func FindUser(ctx context.Context, username string, password string) (bool, error) {
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

func FindUserByName(ctx context.Context, username string) (bool, error) {
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

func GetIdByName(ctx context.Context, username string) (int, error) {
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
