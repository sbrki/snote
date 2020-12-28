package storage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type DiskStorage struct {
	path string
}

func NewDiskStorage(storagePath string) *DiskStorage {
	// create path if it doesn't exists
	os.MkdirAll(storagePath, 0700)
	os.MkdirAll(path.Join(storagePath, "notes"), 0700)
	os.MkdirAll(path.Join(storagePath, "blobs"), 0700)
	return &DiskStorage{storagePath}
}

func (ds *DiskStorage) LoadNote(id string) (*Note, error) {
	filename := path.Join(ds.path, "notes", id+".json")
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	note := new(Note)
	json.Unmarshal(b, note)
	return note, nil
}

func (ds *DiskStorage) SaveNote(note *Note) error {
	filename := path.Join(ds.path, "notes", note.ID+".json")
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	defer f.Close()
	if err != nil {
		return err
	}
	json, err := json.Marshal(note)
	if err != nil {
		return err
	}
	f.Write(json)
	return nil
}

func (ds *DiskStorage) DeleteNote(id string) error {
	err := os.Remove(path.Join(ds.path, "notes", id+".json"))
	if err != nil {
		return err
	}
	return nil
}

func (ds *DiskStorage) GetAllNoteIDs() ([]string, error) {
	notesPath := path.Join(ds.path, "notes")
	IDs := make([]string, 0)
	files, err := ioutil.ReadDir(notesPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), ".json") {
			ID := strings.Split(file.Name(), ".json")[0]
			IDs = append(IDs, ID)
		}
	}
	return IDs, nil
}

func (ds *DiskStorage) SaveBlob(id string, data bytes.Buffer) error {
	filename := path.Join(ds.path, "blobs", id)

	dst, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	defer dst.Close()
	if err != nil {
		return err
	}
	dst.Write(data.Bytes())

	return nil
}

func (ds *DiskStorage) LoadBlobPath(id string) (string, error) {
	filename := path.Join(ds.path, "blobs", id)

	// check if the filename exists
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return "", err
		}
	}
	return filename, nil
}

func (ds *DiskStorage) DeleteBlob(id string) error {
	filename := path.Join(ds.path, "blobs", id)
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DiskStorage) GetAllBlobIDs() ([]string, error) {
	blobsPath := path.Join(ds.path, "blobs")
	result := make([]string, 0)

	files, err := ioutil.ReadDir(blobsPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			result = append(result, file.Name())
		}
	}
	return result, nil
}
