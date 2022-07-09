package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FilePathWithInfo struct {
	Path string
	Info os.FileInfo
}

// https://gist.github.com/mustafaydemir/c90db8fcefeb4eb89696e6ccb5b28685
func ScanRecursive(dir string, ignore []string) ([]FilePathWithInfo, []FilePathWithInfo, error) {
	var dirs []FilePathWithInfo
	var files []FilePathWithInfo

	err := filepath.Walk(dir, func(path string, _f os.FileInfo, err error) error {
		skip := false

		for _, i := range ignore {
			if strings.Contains(path, i) {
				skip = true
			}
		}

		if !skip {
			f, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("stating a file: %v", err)
			}

			if fMode := f.Mode(); fMode.IsDir() {
				dirs = append(dirs, FilePathWithInfo{path, f})
			} else if fMode.IsRegular() {
				files = append(files, FilePathWithInfo{path, f})
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("walking filesystem directory: %v", err)
	}

	return dirs, files, nil
}
