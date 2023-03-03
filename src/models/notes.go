package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Note is a struct that represents a note
type Note struct {
	ID          uint         `json:"id" bd:"note_id"`
	Content     string       `json:"content"`
	User        User         `json:"user" db:"user"`
	Attachments []File       `json:"files" db:"files"` // Relación 1:N
	CreateAt    time.Time    `json:"create_at" db:"created_at"`
	UpdateAt    sql.NullTime `json:"update_at" db:"update_at"`
}

type NoteInterface interface {
	GetNote(noteId uint) (Note, error)
	CreateNote(content string, userId uint) (Note, error)
	UpdateNote(noteIdId uint, content string) (Note, error)
	DeleteNote(noteIdId int) error
}

type noteBD struct {
	BD      *sqlx.DB
	Context context.Context
}

func (c noteBD) GetNote(noteId uint) (Note, error) {
	var note Note
	err := c.BD.GetContext(c.Context, &note, `
	SELECT
		n.*,
	FROM notes n 
	WHERE note_id = ?`, noteId)
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (c noteBD) CreateNote(content string, userId uint) (Note, error) {
	var note Note
	err := c.BD.GetContext(c.Context, &note, "INSERT INTO notes (content, user_id) VALUES (?, ?)", content, userId)
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (c noteBD) UpdateClipboard(noteId uint, content string) (Note, error) {
	var note Note
	err := c.BD.GetContext(c.Context, &note, "UPDATE notes SET content = ? WHERE id = ?", content, noteId)
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (c noteBD) DeleteClipboard(noteId int) error {
	_, err := c.BD.ExecContext(c.Context, "DELETE FROM notes WHERE note_id = ?", noteId)
	if err != nil {
		return err
	}

	return nil
}
