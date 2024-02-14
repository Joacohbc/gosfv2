package models

import (
	"context"
	"errors"
	"fmt"
	"gosfV2/src/ent"
	"gosfV2/src/ent/file"
	"gosfV2/src/ent/user"
	"gosfV2/src/models/database"
	"gosfV2/src/models/env"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Error: "file/s not found"
	ErrFileNotFound = errors.New("file/s not found")
)

type FileInterface interface {

	// Copia el archivo del Reader al sistema de archivos
	// y guarda la la ruta de archivo en la base de datos
	Create(filename string, userId uint, src io.Reader) (*ent.File, error)

	MoveFileToDir(fileId, dirId uint) (*ent.File, error)

	// Elimina un archivo del sistema de archivos y de la ruta de la base de datos
	Delete(fileId uint) (*ent.File, error)

	// Cambia el estado de un archivo a compartido o no compartido
	SetShared(fileId uint, shared bool) (*ent.File, error)

	// Cambia el nombre de un archivo
	Rename(fileId uint, newName string) (*ent.File, error)

	// Obtiene todos los archivos de un usuario
	GetFilesFromUser(userId uint) ([]*ent.File, error)

	// Obtiene un archivo de un usuario
	GetByIdFromUser(fileId, userId uint) (*ent.File, error)

	// Obtiene un archivo por su ID (sin importar el usuario)
	GetById(fileId uint) (*ent.File, error)

	GetByIds(fileIds []uint) ([]*ent.File, error)

	// Determina si un usuario tiene acceso a un archivo
	IsSharedWith(fileId, userId uint) (bool, error)

	// Agrega un usuario a un archivo
	AddUserToFile(userId, fileId uint) error

	// Elimina un usuario de un archivo
	RemoveUserFromFile(userId, fileId uint) error

	// Obtiene todos los archivos compartidos con un usuario
	GetFilesShared(userId uint) ([]*ent.File, error)

	GetPath(id uint, filename string) string
}

type fileService struct {
	BD      *ent.Client
	Context context.Context
}

func Files() FileInterface {
	return fileService{BD: database.GetMySQL(), Context: context.Background()}
}

//
// Consultas
//

func (f fileService) GetFilesFromUser(userId uint) ([]*ent.File, error) {

	files, err := f.BD.File.Query().
		WithOwner().
		WithSharedWith().
		WithParent().
		WithChildren().
		Where(file.HasOwnerWith(user.ID(userId))).
		Order(ent.Asc(file.FieldFilename)).
		All(f.Context)

	if ent.IsNotFound(err) {
		return nil, ErrFileNotFound
	}

	return files, nil
}

func (f fileService) GetByIdFromUser(fileId, userId uint) (*ent.File, error) {

	file, err := f.BD.File.Query().
		WithOwner().
		WithSharedWith().
		WithParent().
		WithChildren().
		Where(file.ID(fileId), file.HasOwnerWith(user.ID(userId))).
		First(f.Context)

	if ent.IsNotFound(err) {
		return nil, ErrFileNotFound
	}

	return file, nil
}

func (f fileService) GetById(fileId uint) (*ent.File, error) {
	file, err := f.BD.File.Query().
		WithOwner().
		WithSharedWith().
		WithParent().
		WithChildren().
		Where(file.ID(fileId)).
		First(f.Context)

	if ent.IsNotFound(err) {
		return nil, ErrFileNotFound
	}

	return file, nil
}

func (f fileService) GetByIds(fileIds []uint) ([]*ent.File, error) {
	files, err := f.BD.File.Query().
		WithOwner().
		WithSharedWith().
		WithParent().
		WithChildren().
		Where(file.IDIn(fileIds...)).
		WithSharedWith().
		All(f.Context)

	if ent.IsNotFound(err) {
		return nil, ErrFileNotFound
	}

	return files, nil
}

//
// Inserts
//

func (f fileService) GetPath(id uint, filename string) string {
	return filepath.Join(env.Config.FilesDirectory, fmt.Sprint(id)+filepath.Ext(filename))
}

func (f fileService) Create(filename string, userId uint, src io.Reader) (*ent.File, error) {

	tx, err := f.BD.Tx(f.Context)
	if err != nil {
		return nil, err
	}

	file, err := tx.File.Create().
		SetFilename(filename).
		SetOwnerID(userId).
		Save(f.Context)

	if err != nil {
		return nil, err
	}

	path := f.GetPath(file.ID, file.Filename)
	fileOpen, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(fileOpen, src); err != nil {
		return nil, err
	}

	if err := fileOpen.Close(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return file, nil
}

//
// Deletes
//

func (f fileService) Delete(fileId uint) (*ent.File, error) {
	file, err := f.GetById(fileId)
	if err != nil {
		return nil, err
	}

	if err := f.BD.File.DeleteOneID(fileId).Exec(f.Context); err != nil {
		return nil, err
	}

	if err := os.Remove(f.GetPath(file.ID, file.Filename)); err != nil {
		return nil, err
	}

	return file, nil
}

//
// Updates
//

func (f fileService) SetShared(fileId uint, shared bool) (*ent.File, error) {
	exit, err := f.BD.File.Query().Where(file.ID(fileId)).Exist(f.Context)
	if err != nil || !exit {
		return nil, ErrFileNotFound
	}

	file, err := f.BD.File.UpdateOneID(fileId).SetIsShared(shared).Save(f.Context)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f fileService) Rename(fileId uint, newName string) (*ent.File, error) {
	if strings.TrimSpace(newName) == "" {
		return nil, errors.New("filename can't be empty")
	}

	exit, err := f.BD.File.Query().Where(file.ID(fileId)).Exist(f.Context)
	if err != nil || !exit {
		return nil, ErrFileNotFound
	}

	file, err := f.BD.File.UpdateOneID(fileId).
		SetFilename(newName).
		Save(f.Context)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f fileService) MoveFileToDir(fileId, dirId uint) (*ent.File, error) {

	exit, err := f.BD.File.Query().Where(file.ID(fileId), file.IsDir(true)).Exist(f.Context)
	if err != nil {
		return nil, ErrFileNotFound
	}

	if !exit {
		return nil, errors.New("file (" + fmt.Sprint(dirId) + ") is not a directory")
	}

	file, err := f.BD.File.UpdateOneID(fileId).
		SetParentID(dirId).
		Save(f.Context)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// Shared Options
func (f fileService) IsSharedWith(fileId, userId uint) (bool, error) {
	exists, err := f.BD.File.Query().
		WithOwner().
		WithSharedWith().
		WithParent().
		WithChildren().
		Where(file.ID(fileId), file.HasSharedWithWith(user.ID(userId))).
		Exist(f.Context)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (f fileService) AddUserToFile(userId, fileId uint) error {
	_, err := f.BD.File.UpdateOneID(fileId).
		AddSharedWithIDs(userId).
		Save(f.Context)
	if err != nil {
		return err
	}

	return nil
}

func (f fileService) RemoveUserFromFile(userId, fileId uint) error {
	_, err := f.BD.File.UpdateOneID(fileId).
		RemoveSharedWithIDs(userId).
		Save(f.Context)
	if err != nil {
		return err
	}

	return nil
}

func (f fileService) GetFilesShared(userId uint) ([]*ent.File, error) {
	files, err := f.BD.File.Query().
		WithOwner().
		WithSharedWith().
		WithParent().
		WithChildren().
		Where(file.HasSharedWithWith(user.ID(userId))).
		Order(ent.Asc(file.FieldFilename)).
		All(f.Context)

	if err != nil {
		return nil, err
	}

	return files, nil
}
