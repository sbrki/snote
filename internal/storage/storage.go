package storage

import (
	"bytes"
)

// Storage interface represents storage for both notes and user-uploaded blobs.
// In order to add a new storage type to snote, this interface has to be
// implemented. Performance wise this interface is not optimal,as it aims to be
// a good ratio of performance and ease of adding new implementations.
// (should be reasonably fast even for ~10K notes, ~10K blobs, ~1K tags).
type Storage interface {
	// LoadNote fetches a Note from storage.
	LoadNote(id string) (*Note, error)
	// SaveNote saves a Note to storage.
	SaveNote(note *Note) error
	// DeleteNote removes a Note from storage.
	DeleteNote(id string) error
	// GetAllNoteIDs fetches IDs of all  notes currently in storage.
	// It is primarily used to display all notes.
	GetAllNoteIDs() ([]string, error)

	// LoadBlobPath fetches file path of the stored blob.
	// This path has to be accessible by the server, as blobs
	// are sent to client as files from disk.
	//
	// This approach is potentially problematic for storages that offer trivial
	// file fetching (s3 and variants) as the implemetation first has to fetch the
	// file from storage to local disk, and then the server sends it
	// to the client (waste/duplication of I/O).
	// TODO: meditate on resolving this inefficiency.
	LoadBlobPath(id string) (string, error)
	// SaveBlob saves a blob to storage.
	SaveBlob(id string, data bytes.Buffer) error
	// DeleteBlob deletes a blob from storage.
	DeleteBlob(id string) error
	// GetAllBlobIDs fetches IDs of all blobs currently in storage.
	GetAllBlobIDs() ([]string, error)

	// SetNoteTags sets the tags of a particular note. Note struct itself has no knowledge
	// of its tags, they are ment to be stored in a seperate structure/index.
	// SetNoteTags is also used to remove tags. For example if a note had tags ["a","b"],
	// and SetNoteTags was called with ["b", "c"], the resulting tags would be ["b", "c"].
	SetNoteTags(id string, tags []string) error
	// GetAllNoteTags fetches all tags along with all the corresponding note IDs.
	GetAllNoteTags() (map[string][]string, error)
}
