package models

import (
	"errors"
	"sync"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

// Note is a struct that represents a note
type Note struct {
	Content string `json:"content"`
}

var NoteManager NoteInterface = &noteBD{mapNote: sync.Map{}}

type NoteInterface interface {
	GetNote(userId uint) (Note, error)
	CreateNote(userId uint, content string) (Note, error)
	UpdateNote(userId uint, content string) (Note, error)
	DeleteNote(userId uint) error
}

type noteBD struct {
	mapNote sync.Map
}

func (n *noteBD) GetNote(userId uint) (Note, error) {
	value, ok := n.mapNote.Load(userId)
	if !ok {
		return Note{}, ErrNoteNotFound
	}

	return value.(Note), nil
}

func (n *noteBD) CreateNote(userId uint, content string) (Note, error) {
	n.mapNote.Store(userId, Note{Content: content})
	return Note{}, nil
}

func (n *noteBD) UpdateNote(userId uint, content string) (Note, error) {
	n.mapNote.Store(userId, Note{Content: content})
	return Note{}, nil
}

func (n *noteBD) DeleteNote(userId uint) error {
	n.mapNote.Delete(userId)
	return nil
}
