package models

import (
	"errors"
	"gosfV2/src/models/db"
	"gosfV2/src/models/env"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func init() {
	bd := db.GetBd()

	if err := bd.AutoMigrate(&File{}); err != nil {
		log.Fatal("Error to create table files:", err)
	}

	if err := os.MkdirAll(env.Config.FilesDirectory, 0644); err != nil {
		log.Fatal("Error to create directory files:", err)
	}
}

type File struct {
	gorm.Model
	Filename   string `json:"filename" gorm:"not null; index:uk1,unique"`
	Shared     bool   `json:"shared" gorm:"default:false"`
	OwnerID    uint   `json:"owner_id" gorm:"not null; index:uk1,unique"`
	Owner      User   `json:"owner"`
	SharedWith []User `json:"shared_with" gorm:"many2many:file_users"`
}

type FileDTO struct {
	ID         uint      `json:"id"`
	Filename   string    `json:"filename"`
	Shared     bool      `json:"shared"`
	SharedWith []UserDTO `json:"shared_with" gorm:"many2many:file_users"`
}

var (
	ErrFileNotFound = errors.New("file/s not found")
)

func Files(c echo.Context) FileFuncs {
	return FileFuncs{BD: db.GetBdCtx(c.Request().Context()), Owner: c.Get("username").(string)}
}

type FileFuncs struct {
	Owner string
	BD    *gorm.DB
}

func (f FileFuncs) getUser() User {
	return User{Username: f.Owner}
}

func (f FileFuncs) MakeUserDir() error {
	return os.MkdirAll(filepath.Join(env.Config.FilesDirectory, f.Owner), 0744)
}

func (f FileFuncs) GetPath(file string) string {
	return filepath.Join(env.Config.FilesDirectory, f.Owner, file)
}

func (f FileFuncs) GetPathFromUser(user, file string) string {
	return filepath.Join(env.Config.FilesDirectory, user, file)
}

func (f FileFuncs) GetAll() ([]FileDTO, error) {

	var files []FileDTO
	err := f.BD.Model(&File{}).
		Joins("JOIN users ON files.owner_id = users.id").
		Where("users.username = ?", f.Owner).
		Find(&files).Error

	if err != nil {
		if db.IsNotFound(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return files, nil
}

func (f FileFuncs) GetById(id uint) (File, error) {
	var file File

	err := f.BD.Preload("Owner").
		Preload("SharedWith").
		Joins("JOIN users ON files.owner_id = users.id").
		Where("users.username = ? AND files.id = ?", f.Owner, id).
		First(&file, id).Error

	if err != nil {
		if db.IsNotFound(err) {
			return file, ErrFileNotFound
		}
		return file, err
	}

	return file, nil
}

func (f FileFuncs) GetShareFile(id uint) (File, error) {
	var file File

	err := f.BD.Preload("Owner").
		Preload("SharedWith").
		First(&file, id).Error

	if err != nil {
		if db.IsNotFound(err) {
			return File{}, ErrFileNotFound
		}
		return File{}, err
	}
	return file, nil
}

func (f FileFuncs) Create(filename string, src io.Reader) error {

	// Inicio una transaccion
	tx := f.BD.Begin()

	file := File{
		Filename: filename,
		Owner:    f.getUser(),
		Shared:   false,
	}

	if err := tx.Create(&file).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Creo el Directorio de archivos del usuario
	if err := f.MakeUserDir(); err != nil {
		tx.Rollback()
		return err
	}

	// Creo el archivo (en Local)
	path := f.GetPath(file.Filename)
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

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (f FileFuncs) Delete(fileId uint) error {

	tx := f.BD.Begin()

	file, err := f.GetById(fileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&file).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := os.Remove(f.GetPath(file.Filename)); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (f FileFuncs) SetShared(fileId uint, shared bool) error {

	file, err := f.GetById(fileId)
	if err != nil {
		return err
	}

	file.Shared = shared
	if err := f.BD.Save(&file).Error; err != nil {
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

	if err := os.Rename(f.GetPath(file.Filename), f.GetPath(newName)); err != nil {
		return err
	}

	file.Filename = newName
	if err := f.BD.Save(&file).Error; err != nil {
		return err
	}

	return nil
}

// func (f FileFuncs) IsSharedWith(ctx context.Context, user string, fileId uint) (bool, error) {
// 	tx := db.GetBdCtx(ctx)

// 	err := tx.Model(&File{Model: gorm.Model{ID: fileId}}).
// 		Association("SharedWith").
// 		Find(&User{Username: user})

// 	if err != nil {
// 		if db.IsNotFound(err) {
// 			return false, nil
// 		}
// 		return false, err
// 	}

// 	return true, nil
// }

// func (f FileFuncs) AddSharedWith(ctx context.Context, user string, fileId uint) error {
// 	tx := db.GetBdCtx(ctx)

// 	err := tx.Model(&File{Model: gorm.Model{ID: fileId}}).
// 		Association("SharedWith").
// 		Append(User{Username: user})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (f FileFuncs) RemoveSharedWith(ctx context.Context, user string, fileId uint) error {
// 	tx := db.GetBdCtx(ctx)

// 	err := tx.Model(&File{Model: gorm.Model{ID: fileId}}).
// 		Association("SharedWith").
// 		Delete((User{Username: user}))

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
