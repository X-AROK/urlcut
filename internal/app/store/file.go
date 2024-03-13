package store

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/X-AROK/urlcut/internal/util"
)

type fileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func newFileWriter(fname string) (*fileWriter, error) {
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &fileWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (fw *fileWriter) WriteRecord(url *url.URL) error {
	return fw.encoder.Encode(url)
}

func (fw *fileWriter) Close() error {
	return fw.file.Close()
}

type FileStore struct {
	mx     sync.Mutex
	values map[string]*url.URL
	writer *fileWriter
}

func (fs *FileStore) Add(ctx context.Context, url *url.URL) (string, error) {
	id := util.GenerateID(8)

	url.ShortURL = id
	err := fs.writer.WriteRecord(url)

	if err != nil {
		return "", err
	}

	fs.mx.Lock()
	fs.values[id] = url
	fs.mx.Unlock()

	return id, nil
}

func (fs *FileStore) Get(ctx context.Context, id string) (*url.URL, error) {
	fs.mx.Lock()
	v, ok := fs.values[id]
	fs.mx.Unlock()

	if !ok {
		return v, url.ErrorNotFound
	}
	return v, nil
}

func (fs *FileStore) parse(fname string) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewBuffer(data))
	var url *url.URL
	fs.mx.Lock()
	for {
		err := decoder.Decode(&url)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fs.values[url.ShortURL] = url
	}
	fs.mx.Unlock()

	return nil
}

func (fs *FileStore) Close() error {
	return fs.writer.Close()
}

func createDir(fname string) error {
	dir := filepath.Dir(fname)
	return os.MkdirAll(dir, 0644)
}

func NewFileStore(fname string) (*FileStore, error) {
	if err := createDir(fname); err != nil {
		return nil, err
	}

	writer, err := newFileWriter(fname)
	if err != nil {
		return nil, err
	}

	fs := &FileStore{
		values: make(map[string]*url.URL),
		writer: writer,
	}
	if err := fs.parse(fname); err != nil {
		return nil, err
	}

	return fs, nil
}
