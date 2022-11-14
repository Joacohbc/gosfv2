package models

import (
	"context"
	"database/sql"
	"errors"
	"gosfV2/src/models/database"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func init() {
	database.GetBd().MustExec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)

}

type User struct {
	ID       uint         `json:"id" db:"user_id"`
	Username string       `json:"username"`
	Password string       `json:"password"`
	CreateAt time.Time    `json:"create_at" db:"created_at"`
	UpdateAt sql.NullTime `json:"update_at" db:"update_at"`
}

var (
	ErrUserNotFound = errors.New("user/s not found")
)

type UsersFuncs struct {
	DB      *sqlx.DB
	Context context.Context
}

func Users(c echo.Context) UsersFuncs {
	return UsersFuncs{DB: database.GetBd(), Context: c.Request().Context()}
}

func UsersC(ctx context.Context) UsersFuncs {
	return UsersFuncs{DB: database.GetBd(), Context: ctx}
}

func (u UsersFuncs) GetAllUsers() ([]User, error) {
	var users []User
	err := u.DB.SelectContext(u.Context, &users, "SELECT * FROM users")
	if err != nil {
		if database.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return users, nil
}

// Crea un nuevo usuario
func (u UsersFuncs) NewUser(user User) error {
	_, err := u.DB.ExecContext(u.Context, "INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	return err
}

// Devuelve el usuario con el nombre de usuario
// Si no encuentra ninguno, devuelve ErrUserNotFound
func (u UsersFuncs) FindUserByName(username string) (User, error) {

	var user User
	err := u.DB.GetContext(u.Context, &user, "SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		if database.IsNotFound(err) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

// Devuelve el usuario con el nombre de usuario
// Si no encuentra ninguno, devuelve ErrUserNotFound
func (u UsersFuncs) FindUserById(id uint) (User, error) {

	var user User
	err := u.DB.GetContext(u.Context, &user, "SELECT * FROM users WHERE user_id = ?", id)
	if err != nil {
		if database.IsNotFound(err) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

// Devuelve el usuario con el nombre de usuario
// Si no encuentra ninguno, no devuelve error (solo retorna false)
func (u UsersFuncs) ExistUserByName(username string) (bool, error) {
	_, err := u.FindUserByName(username)
	if err != nil {
		if err == ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
