package lib

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func GetFilesAsByte(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes, nil
}

func Save(path string, file io.Reader) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ListFiles(directory string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	return paths, err
}
