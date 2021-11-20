package cache

import (
	"os"
	"path/filepath"
	"sync"
)

type File struct {
	path string
	lock sync.RWMutex
}

func NewFile(path string) *File {
	return &File{path: path}
}
func (b *File) Write(data []byte) (err error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	os.MkdirAll(filepath.Dir(b.path), os.ModePerm)
	return os.WriteFile(b.path, data, os.ModePerm)
}
func (b *File) Read() ([]byte, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return os.ReadFile(b.path)
}
func (b *File) Remove() error {
	b.lock.Lock()
	defer b.lock.Unlock()
	return os.Remove(b.path)
}
