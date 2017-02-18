package app

import (
	"os"
	"path/filepath"
)

// file struct
type file struct {
	filename  string
	directory bool
	content   []byte
}

func (f *file) path(dir string) string {
	return filepath.Join(dir, f.filename)
}

func (f *file) isDirectory() bool {
	return f.directory
}

func (f *file) createIfNotExist(dir string) error {

	path := f.path(dir)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		switch f.directory {
		case true:
			os.Mkdir(path, os.ModePerm)
		case false:
			out, err := os.Create(path)
			if err != nil {
				return err
			}
			defer out.Close()
			if _, err := out.Write(f.content); err != nil {
				return err
			}
		}
	}

	return nil
}
