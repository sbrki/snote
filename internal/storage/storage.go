package storage

import (
	"bytes"
)

type Storage interface {
	LoadNote(id string) (*Note, error)
	SaveNote(note *Note) error
	DeleteNote(id string) error
	GetAllNoteIDs() ([]string, error)

	SaveBlob(id string, data bytes.Buffer) error
	LoadBlobPath(id string) (string, error)
	//DeleteBlob(id string) error
}
