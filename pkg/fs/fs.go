package fs

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func UserPath(path string) string {
	if u, err := user.Current(); err == nil {
		if path == "~" {
			return u.HomeDir
		} else if strings.HasPrefix(path, "~/") {
			return filepath.Join(u.HomeDir, path[2:])
		}
	}
	return path
}

func FilesInCurrentDir() ([]string, error) {
	var files []string
	wd, err := os.Getwd()
	if err != nil {
		return files, fmt.Errorf("error to get now directory: %w", err)
	}
	return FilesInDir(wd)
}

func FilesInDir(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files, fmt.Errorf("error to get entries in directory: %w", err)
	}
	for _, e := range entries {
		if info, err := os.Stat(e.Name()); err == nil {
			if !info.IsDir() {
				files = append(files, e.Name())
			}
		}
	}
	return files, nil
}
