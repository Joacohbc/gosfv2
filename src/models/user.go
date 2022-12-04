package models

import (
	"context"
	"database/sql"
	"errors"
	"gosfV2/src/models/database"
	"regexp"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type User struct {
	ID       uint         `json:"id" db:"user_id"`
	Username string       `json:"username"`
	Password string       `json:"password"`
	CreateAt time.Time    `json:"create_at" db:"created_at"`
	UpdateAt sql.NullTime `json:"update_at" db:"update_at"`
}

// Equals: Compara dos usuarios (su ID y su nombre de usuario)
func (u *User) Equals(u2 User) bool {
	return u.ID == u2.ID && u.Username == u2.Username
}

func (u User) Validate() error {

	if len(u.Username) > 255 {
		return errors.New("invalid username, the name is too long (max 255)")
	}

	matched, err := regexp.MatchString("[a-zA-Z0-9]+", u.Username)
	if !matched || err != nil {
		return errors.New("invalid username, the name must be alphanumeric")
	}

	matched, err = regexp.MatchString(".*[a-z]+.*.*[A-Z]+.*.*[0-9]+.*.*[#!?{}!?*^%@#$]+.*", u.Password)
	if !matched || err != nil {
		return errors.New(`invalid password, the password must have at least one lowercase, one uppercase, one number and one special character ("#", "!", "?", "{", "}", "!", "?", "*", "^", "%", "@", "#", "$")`)
	}

	return nil
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
