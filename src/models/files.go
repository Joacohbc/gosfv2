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
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func init() {

	database.GetBd().MustExec(`
	CREATE TABLE IF NOT EXISTS files (
		file_id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		filename TEXT NOT NULL UNIQUE,
		shared BOOLEAN NOT NULL DEFAULT false,
		update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		user_id INT UNSIGNED NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)

	database.GetBd().MustExec(`
	CREATE TABLE IF NOT EXISTS file_users (
		user_id INT UNSIGNED NOT NULL PRIMARY KEY,
		file_id INT UNSIGNED NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (file_id) REFERENCES files(file_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)

}

// File: Representa un archivo de un usuario
type File struct {
	ID         uint         `json:"id" db:"file_id"`
	CreateAt   time.Time    `json:"create_at" db:"created_at"`
	UpdateAt   sql.NullTime `json:"update_at" db:"update_at"`
	Filename   string       `json:"filename"`
	Shared     bool         `json:"shared"`
	UserID     uint         `json:"owner_id" db:"user_id"`
	User       User         `json:"owner" bd:"user"`
	SharedWith []User       `json:"shared_with,omitempty"`
}

var (
	ErrFileNotFound = errors.New("file/s not found")
)

type FileFuncs struct {
	Context context.Context
	BD      *sqlx.DB
	Cancel  context.CancelFunc
}

func Files(c echo.Context) FileFuncs {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return FileFuncs{BD: database.GetBd(), Context: ctx, Cancel: cancel}
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

func (f FileFuncs) GetAll(userId uint) ([]File, error) {
	var files []File

	err := f.BD.SelectContext(f.Context, &files, `SELECT * FROM files WHERE user_id = ?`, userId)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return files, nil
}

func (f FileFuncs) GetByIdFromUser(fileId, userId uint) (File, error) {
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
	WHERE f.file_id = ? AND f.user_id = ?`, fileId, userId)
	if err != nil {
		if database.IsNotFound(err) {
			return File{}, ErrFileNotFound
		}
		return File{}, err
	}

	user, err := UsersC(f.Context).FindUserById(file.UserID)
	if err != nil {
		return File{}, err
	}
	file.User = user

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

	if _, err := tx.Exec("DELETE FROM users WHERE file_id = ?", fileId); err != nil {
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

func (f FileFuncs) GetById(fileId uint) (File, error) {
	var file File

	if err := f.BD.GetContext(f.Context, &file, `
	SELECT 
		f.*,
		u.user_id "user.user_id",
		u.username "user.username",
		u.update_at "user.update_at",
		u.created_at "user.created_at"
	FROM files f
	JOIN users u ON f.user_id  = u.user_id
	WHERE f.file_id = ?;`, fileId); err != nil {
		if database.IsNotFound(err) {
			return File{}, ErrFileNotFound
		}
		return File{}, err
	}

	return file, nil
}

func (f FileFuncs) IsSharedWith(fileId, userId uint) (bool, error) {

	err := f.BD.QueryRowContext(f.Context, "SELECT f.file_id FROM files f JOIN file_users fu ON f.file_id = fu.file_id WHERE f.file_id = ? AND f.user_id = ?", fileId, userId).Err()
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
