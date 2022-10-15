package models

import (
	"context"
	"errors"
	"gosfV2/src/models/db"
	"gosfV2/src/models/env"
	"io"
	"log"
	"os"
	"path/filepath"

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
	Files           fileFuncs
	ErrFileNotFound = errors.New("file/s not found")
)

type fileFuncs struct{}

func (f fileFuncs) MakeDirForUser(username string) error {
	return os.MkdirAll(filepath.Join(env.Config.FilesDirectory, username), 0744)
}

func (f fileFuncs) GetPathUser(username, file string) string {
	return filepath.Join(env.Config.FilesDirectory, username, file)
}

func (f fileFuncs) GetAllFilesByUsername(ctx context.Context, username string) ([]FileDTO, error) {
	bd := db.GetBdCtx(ctx)

	var files []FileDTO
	err := bd.Model(&File{}).
		Joins("JOIN users ON files.owner_id = users.id").
		// Select("files.id, files.filename, files.shared").
		Where("users.username = ?", username).
		Find(&files).Error

	if err != nil {
		if db.IsNotFound(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return files, nil
}

func (f fileFuncs) GetFileByUsername(ctx context.Context, username, filename string) (File, error) {
	bd := db.GetBdCtx(ctx)

	var file File
	err := bd.Joins("JOIN users ON files.owner_id = users.id").
		Where("users.username = ? AND files.filename = ?", username, filename).
		First(&file).Error

	if err != nil {
		if db.IsNotFound(err) {
			return file, ErrFileNotFound
		}
		return file, err
	}

	return file, nil
}

func (f fileFuncs) GetFileById(ctx context.Context, id string) (File, error) {
	bd := db.GetBdCtx(ctx)

	var file File
	if err := bd.First(&file, id).Error; err != nil {
		if db.IsNotFound(err) {
			return file, ErrFileNotFound
		}
		return file, err
	}

	return file, nil
}

func (f fileFuncs) CreateFile(ctx context.Context, filename string, username string, src io.Reader) error {

	// Inicio una transaccion
	tx := db.GetBdCtx(ctx).Begin()

	err := tx.Create(&File{
		Filename: filename,
		Owner:    User{Username: username},
		Shared:   false,
	}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	// Creo el Directorio de archivos del usuario
	if err := f.MakeDirForUser(username); err != nil {
		tx.Rollback()
		return err
	}

	// Creo el archivo (en Local)
	path := f.GetPathUser(username, filename)
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

func (f fileFuncs) DeleteFile(ctx context.Context, username, filename string) error {

	// Inicio una transaccion
	tx := db.GetBdCtx(ctx).Begin()

	res := tx.Unscoped().
		Where(File{
			Filename: filename,
			Owner:    User{Username: username}}).
		Delete(&File{})

	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	if res.RowsAffected == 0 {
		tx.Rollback()
		return ErrFileNotFound
	}

	// Creo el archivo (en Local)
	path := f.GetPathUser(username, filename)
	if err := os.Remove(path); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (f fileFuncs) SetShared(ctx context.Context, username, filename string, shared bool) error {
	res := db.GetBdCtx(ctx).Model(&File{}).
		Where(File{
			Filename: filename,
			Owner:    User{Username: username}}).
		Update("shared", shared)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrFileNotFound
	}

	return nil
}
