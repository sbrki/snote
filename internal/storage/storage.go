package storage

type Storage interface {
	LoadNote(id string) (error, *Note)
	SaveNote(note *Note) error
}
