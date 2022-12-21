package fileutil

import (
	"io/fs"
	"os"
	"path/filepath"
)

var (
	mkdirAllFunc  = os.MkdirAll
	removeAllFunc = os.RemoveAll
	removeFunc    = os.Remove
)

// MkDirIfNotExist creates given dir if it's not existed.
func MkDirIfNotExist(path string) error {
	if !Exist(path) {
		if e := mkdirAllFunc(path, os.ModePerm); e != nil {
			return e
		}
	}
	return nil
}

// RemoveDir deletes dir include children if exist.
func RemoveDir(path string) error {
	if Exist(path) {
		if e := removeAllFunc(path); e != nil {
			return e
		}
	}
	return nil
}

// RemoveFile removes the file if it's existed.
func RemoveFile(file string) error {
	if Exist(file) {
		if e := removeFunc(file); e != nil {
			return e
		}
	}
	return nil
}

// MkDir creates dir.
func MkDir(path string) error {
	if e := mkdirAllFunc(path, os.ModePerm); e != nil {
		return e
	}
	return nil
}

// ListDir reads the directory named by dirname and returns a list of filename.
func ListDir(path string) ([]string, error) {
	var result []string
	if err := readDir(path, func(f fs.DirEntry) {
		result = append(result, f.Name())
	}); err != nil {
		return nil, err
	}
	return result, nil
}

// GetDirectoryList reads the directory named by dirname and returns a list of directory.
func GetDirectoryList(path string) ([]string, error) {
	var result []string
	if err := readDir(path, func(f fs.DirEntry) {
		if f.IsDir() {
			result = append(result, f.Name())
		}
	}); err != nil {
		return nil, err
	}
	return result, nil
}

// Exist check if file or dir is existed.
func Exist(file string) bool {
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// GetExistPath gets exist path based on the given path.
func GetExistPath(path string) string {
	if Exist(path) {
		return path
	}
	dir, _ := filepath.Split(path)
	length := len(dir)
	if length == 0 {
		return dir
	}
	if length > 0 && os.IsPathSeparator(dir[length-1]) {
		dir = dir[:length-1]
	}
	return GetExistPath(dir)
}

// readDir lists all files/directories.
func readDir(path string, fn func(f fs.DirEntry)) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		fn(file)
	}
	return nil
}
