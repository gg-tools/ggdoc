package internal

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryItem[T any] struct {
	IsDir    bool
	Name     string
	Path     string
	Item     *T
	Children []DirectoryItem[T]
}

func DirectoryTree[T any](dirPath string, f func(string, fs.FileInfo) (*T, bool)) ([]DirectoryItem[T], error) {
	var items []DirectoryItem[T]

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == dirPath {
			return nil
		}

		relPath, _ := filepath.Rel(dirPath, path)
		parts := strings.Split(relPath, string(filepath.Separator))
		name := parts[0]
		t, ok := f(path, info)
		if !ok {
			return nil
		}
		item := DirectoryItem[T]{
			Name:  name,
			IsDir: info.IsDir(),
			Path:  path,
			Item:  t,
		}
		if info.IsDir() {
			item.Children, err = DirectoryTree[T](path, f)
			if err != nil {
				return err
			}
			items = append(items, item)
			return filepath.SkipDir
		} else {
			items = append(items, item)
			return nil
		}
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}
