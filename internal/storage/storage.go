package storage

type Storage interface {
	LoadNote(id string) (*Note, error)
	SaveNote(note *Note) error
	DeleteNote(id string) error
	GetAllNoteIDs() ([]string, error)

	//CreateBlob(srcHandle *multipart.FileHeader) error
	//GetBlobPath(id string) string
	//DeleteBlob(id string) error
}
