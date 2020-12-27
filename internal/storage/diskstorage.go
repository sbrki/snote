package storage

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"hash"
	"io"
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
	os.MkdirAll(path.Join(storagePath, "blobs"), 0700)
	return &DiskStorage{storagePath}
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

func (ds *DiskStorage) CreateBlob(src io.Reader, hashAlg hash.Hash) (id string, err error) {
	// since src is underneath a multipart stream, calculating it's hash would consume it,
	// so it needs to be copied into a memory buffer first. Probably dangerous for huge files.
	buffer := new(bytes.Buffer)

	if _, err := io.Copy(buffer, src); err != nil {
		return id, err
	}

	hashAlg.Write(buffer.Bytes())
	id = hex.EncodeToString(hashAlg.Sum(nil))

	// with the calculated file checksum (id), write the buffer to a file
	filename := path.Join(ds.path, "blobs", id)
	dst, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	defer dst.Close()
	if err != nil {
		return id, err
	}
	dst.Write(buffer.Bytes())

	return id, nil
}

func (ds *DiskStorage) GetBlobPath(id string) (string, error) {
	filename := path.Join(ds.path, "blobs", id)

	// check if the filename exists
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return "", err
		}
	}
	return filename, nil
}
