package store

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/X-AROK/urlcut/internal/util"
)

type fileRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

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

func (fw *fileWriter) WriteRecord(r *fileRecord) error {
	return fw.encoder.Encode(r)
}

func (fw *fileWriter) Close() error {
	return fw.file.Close()
}

type FileStore struct {
	mx     sync.Mutex
	values map[string]url.URL
	writer *fileWriter
	lastID int
}

func (fs *FileStore) Add(v url.URL) (string, error) {
	id := util.GenerateID(8)

	record := &fileRecord{
		UUID:        strconv.Itoa(fs.lastID + 1),
		ShortURL:    id,
		OriginalURL: v.Addr,
	}
	err := fs.writer.WriteRecord(record)

	if err != nil {
		return "", err
	}
	fs.lastID += 1

	fs.mx.Lock()
	fs.values[id] = v
	fs.mx.Unlock()

	return id, nil
}

func (fs *FileStore) Get(id string) (url.URL, error) {
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
	var record *fileRecord
	fs.mx.Lock()
	for {
		err := decoder.Decode(&record)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		url := url.NewURL(record.OriginalURL)
		fs.values[record.ShortURL] = url
	}
	fs.mx.Unlock()

	if record != nil {
		lastID, err := strconv.Atoi(record.UUID)
		if err != nil {
			return err
		}
		fs.lastID = lastID
	}

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
		values: make(map[string]url.URL),
		writer: writer,
	}
	if err := fs.parse(fname); err != nil {
		return nil, err
	}

	return fs, nil
}
