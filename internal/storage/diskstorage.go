package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type DiskStorage struct {
	path string
}

func NewDiskStorage(path string) *DiskStorage {
	// create path if it doesn't exists
	os.MkdirAll(path, 0700)
	return &DiskStorage{path}
}

func (ds *DiskStorage) LoadNote(id string) (*Note, error) {
	filename := path.Join(ds.path, id+".json")
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	note := new(Note)
	json.Unmarshal(b, note)
	return note, nil
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

func (ds *DiskStorage) DeleteNote(id string) error {
	err := os.Remove(path.Join(ds.path, id+".json"))
	if err != nil {
		return err
	}
	return nil
}

func (ds *DiskStorage) GetAllNoteIDs() ([]string, error) {
	IDs := make([]string, 0)
	files, err := ioutil.ReadDir(ds.path)
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
