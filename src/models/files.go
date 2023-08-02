package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gosfV2/src/models/database"
	"gosfV2/src/models/env"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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
	User       User   `json:"owner" bd:"user"`                      // "user" es un prefix para usar sqlx.Get y sqlx.Select con una sola consulta
	SharedWith []User `json:"sharedWith,omitempty" db:"sharedWith"` // "sharedWith" es un prefix para usar sqlx.Get y sqlx.Select con una sola consulta
}

func (f *File) Validate() error {
	match, err := regexp.MatchString(`[a-zA-Z0-9_-]+(\.)[a-z]+`, f.Filename)
	if !match || err != nil {
		return errors.New(`invalid filename format, only letters, numbers, "-" and "_" are allowed`)
	}

	return nil
}

func (f *File) GetPath() string {
	return filepath.Join(env.Config.FilesDirectory, fmt.Sprint(f.ID)+filepath.Ext(f.Filename))
}

var (
	// Error: "file/s not found"
	ErrFileNotFound = errors.New("file/s not found")
)

type FileInterface interface {

	// Copia el archivo del Reader al sistema de archivos
	// y guarda la la ruta de archivo en la base de datos
	Create(filename string, userId uint, src io.Reader) (uint, error)

	// Elimina un archivo del sistema de archivos y de la ruta de la base de datos
	Delete(fileId uint) error

	// Cambia el estado de un archivo a compartido o no compartido
	SetShared(fileId uint, shared bool) error

	// Cambia el nombre de un archivo
	Rename(fileId uint, newName string) error

	// Obtiene todos los archivos de un usuario
	GetFilesFromUser(userId uint) ([]File, error)

	// Obtiene un archivo de un usuario
	GetByIdFromUser(fileId, userId uint) (File, error)

	// Obtiene un archivo de un usuario por su nombre
	GetByFilenameFromUser(filename string, userId uint) (File, error)

	// Obtiene todos los usuarios que tienen acceso a un archivo
	GetUsersFromFile(filedId uint) ([]User, error)

	// Obtiene un archivo por su ID (sin importar el usuario)
	GetById(fileId uint) (File, error)

	GetByIds(fileIds []uint) ([]File, error)

	// Determina si un usuario tiene acceso a un archivo
	IsSharedWith(fileId, userId uint) (bool, error)

	// Agrega un usuario a un archivo
	AddUserToFile(userId, fileId uint) error

	// Elimina un usuario de un archivo
	RemoveUserFromFile(userId, fileId uint) error

	// Obtiene todos los archivos compartidos con un usuario
	GetFilesShared(userId uint) ([]File, error)

	ManageError(err error) error
}

type fileBD struct {
	BD      *sqlx.DB
	Context context.Context
}

func Files(c echo.Context) FileInterface {
	return fileBD{BD: database.GetMySQL(), Context: c.Request().Context()}
}

func (f fileBD) ManageError(err error) error {
	if database.IsNotFound(err) {
		return ErrFileNotFound
	}
	return err
}

//
// Consultas
//

func (f fileBD) GetFilesFromUser(userId uint) ([]File, error) {
	var files []File

	err := f.BD.SelectContext(f.Context, &files, `
	SELECT 
		f.* 
	FROM files f 
	WHERE user_id = ?
	ORDER BY f.filename`, userId)
	if err != nil {
		return nil, f.ManageError(err)
	}

	return files, nil
}

func (f fileBD) GetUsersFromFile(filedId uint) ([]User, error) {

	var users []User
	err := f.BD.Select(&users, `
	SELECT 
		u.user_id, 
		u.username 
	FROM users u
	JOIN file_users fu ON u.user_id = fu.user_id 
	WHERE fu.file_id = ?;`, filedId)

	if err != nil {
		return nil, f.ManageError(err)
	}

	return users, nil
}

func (f fileBD) GetById(fileId uint) (File, error) {
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
		return File{}, f.ManageError(err)
	}

	users, err := f.GetUsersFromFile(file.ID)
	if err != nil {
		return File{}, err
	}
	file.SharedWith = users

	return file, nil
}

func (f fileBD) GetByIds(fileIds []uint) ([]File, error) {
	var files []File

	query, args, err := sqlx.In(`
	SELECT 
		f.*,
		u.user_id "user.user_id",
		u.username "user.username",
		u.update_at "user.update_at",
		u.created_at "user.created_at"
	FROM files f
	JOIN users u ON f.user_id  = u.user_id
	WHERE f.file_id IN (?)`, fileIds)
	if err != nil {
		return nil, f.ManageError(err)
	}

	err = f.BD.SelectContext(f.Context, &files, query, args...)
	if err != nil {
		return nil, f.ManageError(err)
	}

	return files, nil
}

func (f fileBD) GetByIdFromUser(fileId, userId uint) (File, error) {
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

func (f fileBD) GetByFilenameFromUser(filename string, userId uint) (File, error) {
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
	WHERE f.filename = ? AND u.user_id = ?;`, filename, userId)
	if err != nil {
		if database.IsNotFound(err) {
			return File{}, ErrFileNotFound
		}
		return File{}, err
	}

	users, err := f.GetUsersFromFile(file.ID)
	if err != nil {
		return File{}, err
	}
	file.SharedWith = users

	return file, nil
}

//
// Inserts
//

func (f fileBD) Create(filename string, userId uint, src io.Reader) (uint, error) {

	// Inicio una transacción
	tx, err := f.BD.BeginTxx(f.Context, nil)
	if err != nil {
		return 0, err
	}

	// Creo el archivo
	if _, err := tx.ExecContext(f.Context, `INSERT INTO files (filename, user_id) VALUES (?,?);`, filename, userId); err != nil {
		tx.Rollback()
		return 0, err
	}

	var fileId uint
	if err := tx.GetContext(f.Context, &fileId, `SELECT LAST_INSERT_ID();`); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Creo el archivo (en Local)
	path := (&File{ID: fileId, Filename: filename}).GetPath()
	fileOpen, err := os.Create(path)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Copio el archivo Remoto en el Local
	if _, err = io.Copy(fileOpen, src); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := fileOpen.Close(); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return fileId, nil
}

//
// Deletes
//

func (f fileBD) Delete(fileId uint) error {

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

	if _, err := tx.Exec("DELETE FROM files WHERE file_id = ?", fileId); err != nil {
		tx.Rollback()
		return err
	}

	if err := os.Remove(file.GetPath()); err != nil {
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

func (f fileBD) SetShared(fileId uint, shared bool) error {

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

func (f fileBD) Rename(fileId uint, newName string) error {

	if strings.TrimSpace(newName) == "" {
		return errors.New("filename can't be empty")
	}

	if _, err := f.GetById(fileId); err != nil {
		return err
	}

	if _, err := f.BD.Exec("UPDATE files SET filename = ? WHERE file_id = ?", newName, fileId); err != nil {
		return err
	}

	return nil
}

// Shared Options
func (f fileBD) IsSharedWith(fileId, userId uint) (bool, error) {

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

func (f fileBD) AddUserToFile(userId, fileId uint) error {
	_, err := f.BD.ExecContext(f.Context, "INSERT INTO file_users (user_id, file_id) VALUES (?,?)", userId, fileId)
	if err != nil {
		return err
	}

	return nil
}

func (f fileBD) RemoveUserFromFile(userId, fileId uint) error {
	_, err := f.BD.ExecContext(f.Context, "DELETE FROM file_users WHERE user_id = ? AND file_id = ?", userId, fileId)
	if err != nil {
		return err
	}

	return nil
}

func (f fileBD) GetFilesShared(userId uint) ([]File, error) {
	var files []File

	err := f.BD.SelectContext(f.Context, &files, `
	SELECT
		f.*
	FROM files f
	JOIN file_users fu ON f.file_id = fu.file_id
	WHERE fu.user_id = ?
	ORDER BY f.filename`, userId)
	if err != nil {
		if database.IsNotFound(err) {
			return []File{}, nil
		}
		return nil, err
	}

	return files, nil
}
