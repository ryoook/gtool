package rootpath

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
)

var (
	selfPath = ""
	rootPath atomic.Value
)

func init() {
	selfPath, _ = exec.LookPath(os.Args[0])
	if selfPath != "" {
		selfPath, _ = filepath.Abs(selfPath)
	}
	if selfPath == "" {
		selfPath, _ = filepath.Abs(os.Args[0])
	}
}

func RootPath() []string {
	if rootPath.Load() != nil {
		return rootPath.Load().([]string)
	}
	var res []string
	// Working dir
	path, err := os.Getwd()
	if err == nil {
		res = append(res, path)
	}
	// Main package dir
	mainPkgPath := getMainPkgPath()
	if mainPkgPath != "" {
		res = append(res, mainPkgPath)
	}

	// Binary dir
	binPath := filepath.Dir(selfPath)
	if binPath != "" {
		res = append(res, binPath)
	}
	rootPath.Store(res)
	return res
}
