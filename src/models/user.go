package models

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gosfV2/src/models/database"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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
	ErrUserNotFound           = errors.New("user/s not found")
	ErrIconTooLarge           = errors.New("icon too large (max 512 x 512)")
	ErrIconFormatNotSupported = errors.New("invalid image format, only jpeg/png/gif and png is supported")
)

type UserInterface interface {
	GetAllUsers() ([]User, error)
	NewUser(user User) error
	FindUserByName(username string) (User, error)
	FindUserById(id uint) (User, error)
	ExistUserByName(username string) (bool, error)
	Rename(id uint, newName string) error
	ChangePassword(id uint, newPassword string) error
	Delete(id uint) error
	UploadIcon(id uint, src io.Reader) error
	DeleteIcon(id uint) error
	GetIcon(id uint) string
	ManageError(err error) error
}

type usersBD struct {
	DB      *sqlx.DB
	Context context.Context
}

func Users(c echo.Context) UserInterface {
	return usersBD{DB: database.GetMySQL(), Context: c.Request().Context()}
}

func UsersC(ctx context.Context) UserInterface {
	return usersBD{DB: database.GetMySQL(), Context: ctx}
}

func (u usersBD) GetAllUsers() ([]User, error) {
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

func (u usersBD) ManageError(err error) error {
	if database.IsNotFound(err) {
		return ErrUserNotFound
	}
	return err
}

// Crea un nuevo usuario
func (u usersBD) NewUser(user User) error {
	_, err := u.DB.ExecContext(u.Context, "INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	return err
}

// Devuelve el usuario con el nombre de usuario
// Si no encuentra ninguno, devuelve ErrUserNotFound
func (u usersBD) FindUserByName(username string) (User, error) {

	var user User
	err := u.DB.GetContext(u.Context, &user, "SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return User{}, u.ManageError(err)
	}

	return user, nil
}

// Devuelve el usuario con el nombre de usuario
// Si no encuentra ninguno, devuelve ErrUserNotFound
func (u usersBD) FindUserById(id uint) (User, error) {

	var user User
	err := u.DB.GetContext(u.Context, &user, "SELECT * FROM users WHERE user_id = ?", id)
	if err != nil {
		return User{}, u.ManageError(err)
	}

	return user, nil
}

// Devuelve el usuario con el nombre de usuario
// Si no encuentra ninguno, no devuelve error (solo retorna false)
func (u usersBD) ExistUserByName(username string) (bool, error) {
	_, err := u.FindUserByName(username)
	if err != nil {
		if err == ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u usersBD) Rename(id uint, newName string) error {
	_, err := u.DB.ExecContext(u.Context, "UPDATE users SET username = ? WHERE user_id = ?", newName, id)
	return err
}

func (u usersBD) ChangePassword(id uint, newPassword string) error {
	_, err := u.DB.ExecContext(u.Context, "UPDATE users SET password = ? WHERE user_id = ?", newPassword, id)
	return err
}

func (u usersBD) Delete(id uint) error {
	_, err := u.DB.ExecContext(u.Context, "DELETE FROM users WHERE user_id = ?", id)
	return err
}

func (u usersBD) UploadIcon(id uint, src io.Reader) error {

	// Leo el archivo (para poder usar el contenido del Reader varias veces)
	blob, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	img, _, err := image.DecodeConfig(bytes.NewReader(blob))
	if err != nil {
		return ErrIconFormatNotSupported
	}

	if img.Width > 512 || img.Height > 512 {
		return ErrIconTooLarge
	}

	// Guardo el archivo
	file, err := os.Create(u.GetIcon(id))
	if err != nil {
		return err
	}

	// Cierro el archivo
	if _, err = io.Copy(file, bytes.NewReader(blob)); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

func (u usersBD) GetIcon(id uint) string {
	return filepath.Join(UserIconDir, fmt.Sprint(id)+"icon")
}

func (u usersBD) DeleteIcon(id uint) error {
	err := os.Remove(u.GetIcon(id))
	if err != nil {
		return err
	}

	return nil
}
