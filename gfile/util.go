package gfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ScanDir(path string, pattern string, recursive ...bool) []string {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	list, err := doScanDir(0, path, pattern, isRecursive, nil)
	if err != nil {
		return nil
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list
}

func doScanDir(depth int, path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	if depth >= 100000 {
		return nil, nil
	}
	var (
		list      []string
		file, err = Open(path)
	)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	var (
		filePath string
		patterns = strings.Split(pattern, ",")
	)
	for _, name := range names {
		filePath = path + string(filepath.Separator) + name
		if IsDir(filePath) && recursive {
			array, _ := doScanDir(depth+1, filePath, pattern, true, handler)
			if len(array) > 0 {
				list = append(list, array...)
			}
		}
		// Handler filtering.
		if handler != nil {
			filePath = handler(filePath)
			if filePath == "" {
				continue
			}
		}
		// If it meets pattern, then add it to the result list.
		for _, p := range patterns {
			if match, _ := filepath.Match(p, name); match {
				if filePath = Abs(filePath); filePath != "" {
					list = append(list, filePath)
				}
			}
		}
	}
	return list, nil
}

func Open(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf(`os.Open failed for name "%s"`, path)
	}
	return file, err
}

func Abs(path string) string {
	p, _ := filepath.Abs(path)
	return p
}
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
