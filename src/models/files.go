package models

import (
	"context"
	"database/sql"
	"errors"
	"gosfV2/src/models/database"
	"gosfV2/src/models/env"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

// File: Representa un archivo de un usuario
type File struct {
	ID       uint         `json:"id" db:"file_id"`
	CreateAt time.Time    `json:"create_at" db:"created_at"`
	UpdateAt sql.NullTime `json:"update_at" db:"update_at"`
	Filename string       `json:"filename"`

	// Shared: Indica que esta compartido con TODOS los usuarios
	Shared     bool   `json:"shared"`
	UserID     uint   `json:"owner_id" db:"user_id"`
	User       User   `json:"owner" bd:"user"`                        // "user" es un prefix para usar sqlx.Get y sqlx.Select con una sola consulta
	SharedWith []User `json:"shared_with,omitempty" db:"shared_with"` // "shared_with" es un prefix para usar sqlx.Get y sqlx.Select con una sola consulta
}

func (f *File) Validate() error {
	match, err := regexp.MatchString(`[a-zA-Z0-9_-]+(\.)[a-z]+`, "Image.png")
	if !match || err != nil {
		return errors.New(`invalid filename format, only letters, numbers, "-" and "_" are allowed`)
	}

	return nil
}

var (
	// Error: "file/s not found"
	ErrFileNotFound = errors.New("file/s not found")
)

type FileFuncs struct {
	BD      *sqlx.DB
	Context context.Context
}

func Files(c echo.Context) FileFuncs {
	return FileFuncs{BD: database.GetBd(), Context: c.Request().Context()}
}

func (f FileFuncs) MakeUserDir(user string) error {
	return os.MkdirAll(filepath.Join(env.Config.FilesDirectory, user), 0744)
}

func (f FileFuncs) GetPath(file, user string) string {
	return filepath.Join(env.Config.FilesDirectory, user, file)
}

//
// Consultas
//

func (f FileFuncs) GetAllFromUser(userId uint) ([]File, error) {
	var files []File

	err := f.BD.SelectContext(f.Context, &files, `
	SELECT 
		f.* 
	FROM files f 
	WHERE user_id = ?`, userId)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return files, nil
}

func (f FileFuncs) GetFilepath(fileId, userId uint) (string, error) {

	var path struct {
		Filename string `db:"filename"`
		Username string `db:"username"`
	}

	err := f.BD.GetContext(f.Context, &path, `
	SELECT 
		f.filename,
		u.username
	FROM files f
	JOIN users u ON f.user_id  = u.user_id
	WHERE f.file_id = ? AND f.user_id = ?`, fileId, userId)
	if err != nil {
		if database.IsNotFound(err) {
			return "", ErrFileNotFound
		}
		return "", err
	}

	return f.GetPath(path.Filename, path.Username), nil
}

func (f FileFuncs) GetByIdFromUser(fileId, userId uint) (File, error) {
	var file File

	file, err := f.GetById(fileId)
	if err != nil {
		return file, err
	}

	if file.UserID != userId {
		return file, ErrFileNotFound
	}

	return file, nil
}

//
// Inserts
//

func (f FileFuncs) Create(filename string, userId uint, src io.Reader) error {

	// Inicio una transacción
	tx, err := f.BD.BeginTxx(f.Context, nil)
	if err != nil {
		return err
	}

	// Creo el archivo
	if _, err := tx.ExecContext(f.Context, `INSERT INTO files (filename,user_id) VALUES (?,?);`, filename, userId); err != nil {
		tx.Rollback()
		return err
	}

	user, err := UsersC(f.Context).FindUserById(userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Creo el Directorio de archivos del usuario
	if err := f.MakeUserDir(user.Username); err != nil {
		tx.Rollback()
		return err
	}

	// Creo el archivo (en Local)
	path := f.GetPath(filename, user.Username)
	fileOpen, err := os.Create(path)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Copio el archivo Remoto en el Local
	if _, err = io.Copy(fileOpen, src); err != nil {
		tx.Rollback()
		return err
	}

	if err := fileOpen.Close(); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

//
// Deletes
//

func (f FileFuncs) Delete(fileId uint) error {

	// Inicio una transacción
	tx, err := f.BD.BeginTxx(f.Context, nil)
	if err != nil {
		return err
	}

	file, err := f.GetById(fileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	user, err := UsersC(f.Context).FindUserById(file.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM files WHERE file_id = ?", fileId); err != nil {
		tx.Rollback()
		return err
	}

	if err := os.Remove(f.GetPath(file.Filename, user.Username)); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

//
// Updates
//

func (f FileFuncs) SetShared(fileId uint, shared bool) error {

	file, err := f.GetById(fileId)
	if err != nil {
		return err
	}

	file.Shared = shared
	if _, err := f.BD.Exec("UPDATE files SET shared = ? WHERE file_id = ?", shared, fileId); err != nil {
		return err
	}

	return nil
}

func (f FileFuncs) Rename(fileId uint, newName string) error {

	if strings.TrimSpace(newName) == "" {
		return errors.New("filename can't be empty")
	}

	file, err := f.GetById(fileId)
	if err != nil {
		return err
	}

	user, err := UsersC(f.Context).FindUserById(file.UserID)
	if err != nil {
		return err
	}

	if err := os.Rename(f.GetPath(file.Filename, user.Username), f.GetPath(newName, user.Username)); err != nil {
		return err
	}

	if _, err := f.BD.Exec("UPDATE files SET filename = ? WHERE file_id = ?", newName, fileId); err != nil {
		return err
	}

	return nil
}

//
// Shared Options
//

func (f FileFuncs) GetAllUserFromFile(filedId uint) ([]User, error) {

	var users []User
	err := f.BD.Select(&users, `
	SELECT 
		u.user_id, 
		u.username 
	FROM users u
	JOIN file_users fu ON u.user_id = fu.user_id 
	WHERE fu.file_id = ?;`, filedId)

	if err != nil {
		if database.IsNotFound(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return users, nil
}

func (f FileFuncs) GetById(fileId uint) (File, error) {
	var file File

	err := f.BD.GetContext(f.Context, &file, `
	SELECT 
		f.*,
		u.user_id "user.user_id",
		u.username "user.username",
		u.update_at "user.update_at",
		u.created_at "user.created_at"
	FROM files f
	JOIN users u ON f.user_id  = u.user_id
	WHERE f.file_id = ?;`, fileId)
	if err != nil {
		if database.IsNotFound(err) {
			return File{}, ErrFileNotFound
		}
		return File{}, err
	}

	users, err := f.GetAllUserFromFile(file.ID)
	if err != nil {
		return File{}, err
	}
	file.SharedWith = users

	return file, nil
}

func (f FileFuncs) IsSharedWith(fileId, userId uint) (bool, error) {

	var foo uint
	err := f.BD.QueryRowContext(f.Context, `
	SELECT 
		f.file_id
	FROM files f 
	JOIN file_users fu ON f.file_id = fu.file_id 
	WHERE f.file_id = ? AND fu.user_id = ?`,
		fileId, userId).Scan(&foo)

	if err != nil {
		if database.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (f FileFuncs) AddUserToFile(userId, fileId uint) error {
	_, err := f.BD.ExecContext(f.Context, "INSERT INTO file_users (user_id, file_id) VALUES (?,?)", userId, fileId)
	if err != nil {
		return err
	}

	return nil
}

func (f FileFuncs) RemoveUserFromFile(userId, fileId uint) error {
	_, err := f.BD.ExecContext(f.Context, "DELETE FROM file_users WHERE user_id = ? AND file_id = ?", userId, fileId)
	if err != nil {
		return err
	}

	return nil
}
