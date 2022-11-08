package indexer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/giantas/gotu/storage"
)

func Run(store storage.FileStorage) error {
	root, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fmt.Printf("walking %s\n", root)
	files, err := index(root)
	if err != nil {
		return err
	}
	fmt.Printf("indexing %s\n", root)
	err = store.CreateMany(files)
	if err != nil {
		return err
	}
	fmt.Printf("%d files indexed\n", len(files))
	return nil
}

func index(root string) ([]*storage.File, error) {
	files := make([]*storage.File, 0)

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if isHidden(path) {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		file := storage.File{
			Name: d.Name(),
			Path: path,
		}
		files = append(files, &file)
		return nil
	})
	return files, nil
}

func isHidden(filePath string) bool {
	// TODO: Adapt for Win systems
	filename := filepath.Base(filePath)
	return filename[0] == '.'
}
