package server

import (
	"github.com/HypercubeHC/go-metrics-alerting.git/internal/storage"
)

type FileStorage struct {
	*storage.MemStorage[string, []byte]
	onWrite func()
}

func NewFileStorage(strg *storage.MemStorage[string, []byte], onWrite func()) *FileStorage {
	return &FileStorage{MemStorage: strg, onWrite: onWrite}
}

func (fs *FileStorage) Set(name string, val []byte) {
	fs.MemStorage.Set(name, val)
	fs.onWrite()
}

func (fs *FileStorage) Update(name string, fn func(prev []byte) []byte) {
	fs.MemStorage.Update(name, fn)
	fs.onWrite()
}
