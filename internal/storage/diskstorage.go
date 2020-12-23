package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type DiskStorage struct {
	path string
}

func NewDiskStorage(path string) *DiskStorage {
	// create path if it doesn't exists
	os.MkdirAll(path, 0700)
	return &DiskStorage{path}
}

func (ds *DiskStorage) LoadNote(id string) (error, *Note) {
	filename := path.Join(ds.path, id+".json")
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err, nil
	}
	note := new(Note)
	json.Unmarshal(b, note)
	return nil, note
}

func (ds *DiskStorage) SaveNote(note *Note) error {
	filename := path.Join(ds.path, note.ID+".json")
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
