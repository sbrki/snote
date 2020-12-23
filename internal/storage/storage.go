package storage

type Storage interface {
	LoadNote(id string) (*Note, error)
	SaveNote(note *Note) error
}
