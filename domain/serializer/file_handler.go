package serializer

import (
	"os"
	"path/filepath"
)

type FileHandler struct {
	baseDir string
}

func MakeFileHandler(baseDir string) *FileHandler {
	return &FileHandler{baseDir: baseDir}
}

func (fh *FileHandler) SaveToFile(filename string, data []byte) error {
	fullPath := filepath.Join(fh.baseDir, filename)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, data, 0664)
}

func (fh *FileHandler) LoadFromFile(filename string) ([]byte, error) {
	fullPath := filepath.Join(fh.baseDir, filename)
	return os.ReadFile(fullPath)
}

func (fh *FileHandler) SaveObject(filename string, obj Serializable, serializer *Serializer) error {
	data, err := serializer.Serialize(obj)
	if err != nil {
		return err
	}
	return fh.SaveToFile(filename, data)
}

func (fh *FileHandler) LoadObject(filename string, obj Serializable, serializer *Serializer) error {
	data, err := fh.LoadFromFile(filename)
	if err != nil {
		return err
	}
	return serializer.Deserialize(data, obj)
}
