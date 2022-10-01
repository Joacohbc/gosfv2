package files

import (
	"context"
	"gosfV2/src/models/db"
	"gosfV2/src/models/env"
	"gosfV2/src/models/users"
	"io"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
}

func init() {
	conn := db.GetConn()

	rows, err := conn.Query(`
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTO_INCREMENT, 
		filename TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		modified DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		id_user INTEGER NOT NULL,
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

func MakeDirForUser(username string) error {
	return os.MkdirAll(filepath.Join(env.Config.FilesDirectory, username), 0644)
}

func GetPathUser(username, file string) string {
	return filepath.Join(env.Config.FilesDirectory, username, file)
}

func GetAllFilesByUsername(ctx context.Context, username string) ([]File, error) {

	id, err := users.GetIdByName(ctx, username)
	if err != nil {
		return nil, err
	}

	conn := db.GetConn()

	rs, err := conn.Query(`SELECT id, filename FROM files WHERE id_user=?`, id)
	if err != nil {
		return nil, err
	}

	var files []File
	for rs.Next() {
		var f File
		if err := rs.Scan(&f.ID, &f.Filename); err != nil {
			return nil, err
		}
		f.Filename = filepath.Base(f.Filename)
		files = append(files, f)
	}

	return files, nil
}

func GetFileByUsername(ctx context.Context, username, filename string) (File, error) {

	id, err := users.GetIdByName(ctx, username)
	if err != nil {
		return File{}, err
	}

	conn := db.GetConn()

	var file File
	err = conn.QueryRow(`SELECT id, filename FROM files WHERE id_user=? AND filename=?`, id, filename).Scan(&file.ID, &file.Filename)
	if err != nil {
		return File{}, err
	}

	return file, nil
}

func CreateFile(ctx context.Context, filename string, username string, src io.Reader) error {

	// Obtengo el ID del usuario
	id, err := users.GetIdByName(ctx, username)
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
	if err := MakeDirForUser(username); err != nil {
		tx.Rollback()
		return err
	}

	// Creo el archivo (en Local)
	path := GetPathUser(username, filename)
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

func DeleteFile(ctx context.Context, username, filename string) error {

	conn := db.GetConn()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Genero el nuevo registro en la Tabla
	_, err = tx.Exec(`DELETE FROM files WHERE filename=?`, filename)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Creo el Directorio de archivos del usuario
	path := GetPathUser(username, filename)
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
