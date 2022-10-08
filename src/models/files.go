package models

import (
	"context"
	"gosfV2/src/models/db"
	"gosfV2/src/models/env"
	"io"
	"log"
	"os"
	"path/filepath"
)

func init() {
	conn := db.GetConn()

	rows, err := conn.Query(`
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTO_INCREMENT, 
		filename TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		modified DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		id_user INTEGER NOT NULL,
		shared BOOLEAN DEFAULT FALSE,
		UNIQUE (filename, id_user),
		FOREIGN KEY (id_user) REFERENCES users(id)
	);`)

	if err != nil {
		log.Fatal("Error to create table files:", err)
	}
	rows.Close()

	if err := os.MkdirAll(env.Config.FilesDirectory, 0644); err != nil {
		log.Fatal("Error to create directory files:", err)
	}
}

type File struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Shared   bool   `json:"shared"`
}

var Files fileFuncs

type fileFuncs struct{}

func (f fileFuncs) MakeDirForUser(username string) error {
	return os.MkdirAll(filepath.Join(env.Config.FilesDirectory, username), 0744)
}

func (f fileFuncs) GetPathUser(username, file string) string {
	return filepath.Join(env.Config.FilesDirectory, username, file)
}

func (f fileFuncs) GetAllFilesByUsername(ctx context.Context, username string) ([]File, error) {

	id, err := Users.GetIdByName(ctx, username)
	if err != nil {
		return nil, err
	}

	conn := db.GetConn()

	rs, err := conn.Query(`SELECT id, filename, shared FROM files WHERE id_user=?`, id)
	if err != nil {
		return nil, err
	}

	var files []File
	for rs.Next() {
		var f File
		if err := rs.Scan(&f.ID, &f.Filename, &f.Shared); err != nil {
			return nil, err
		}
		f.Filename = filepath.Base(f.Filename)
		files = append(files, f)
	}

	return files, nil
}

func (f fileFuncs) GetFileByUsername(ctx context.Context, username, filename string) (File, error) {

	id, err := Users.GetIdByName(ctx, username)
	if err != nil {
		return File{}, err
	}

	conn := db.GetConn()

	var file File
	err = conn.QueryRow(`SELECT id, filename, shared FROM files WHERE id_user=? AND filename=?`, id, filename).Scan(&file.ID, &file.Filename, &file.Shared)
	if err != nil {
		return File{}, err
	}

	return file, nil
}

func (f fileFuncs) CreateFile(ctx context.Context, filename string, username string, src io.Reader) error {

	// Obtengo el ID del usuario
	id, err := Users.GetIdByName(ctx, username)
	if err != nil {
		return err
	}

	conn := db.GetConn()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Genero el nuevo registro en la Tabla
	_, err = tx.Exec(`INSERT INTO files (filename, id_user) VALUES (?, ?)`, filename, id)
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

	// Si todo salio bien hago el Commit a la BD
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (f fileFuncs) DeleteFile(ctx context.Context, username, filename string) error {

	id, err := Users.GetIdByName(ctx, username)
	if err != nil {
		return err
	}

	conn := db.GetConn()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Genero el nuevo registro en la Tabla
	_, err = tx.Exec(`DELETE FROM files WHERE filename=? AND id_user=?`, filename, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Creo el Directorio de archivos del usuario
	path := f.GetPathUser(username, filename)
	if err := os.Remove(path); err != nil {
		tx.Rollback()
		return err
	}

	// Si todo salio bien hago el Commit a la BD
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (f fileFuncs) SetShared(ctx context.Context, username, filename string, shared bool) error {

	file, err := f.GetFileByUsername(ctx, username, filename)
	if err != nil {
		return err
	}

	conn := db.GetConn()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var sql string
	if shared {
		sql = `UPDATE files SET shared=true WHERE filename=? AND id_user=?`
	} else {
		sql = `UPDATE files SET shared=false WHERE filename=? AND id_user=?`
	}

	_, err = tx.Exec(sql, file.Filename, file.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Si todo salio bien hago el Commit a la BD
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
