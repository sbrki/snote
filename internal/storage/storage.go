package storage

import (
	"hash"
	"io"
)

type Storage interface {
	LoadNote(id string) (*Note, error)
	SaveNote(note *Note) error
	DeleteNote(id string) error
	GetAllNoteIDs() ([]string, error)

	CreateBlob(src io.Reader, hashAlg hash.Hash) (id string, err error)
	GetBlobPath(id string) (string, error)
	//DeleteBlob(id string) error
}
