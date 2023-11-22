package gtool

import (
	"github.com/ryoook/gtool/gfile"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	selfPath = ""
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
	var res []string
	// Working dir
	path, err := os.Getwd()
	if err == nil {
		res = append(res, path)
	}
	// Main package dir
	mainPkgPath := gfile.MainPkgPath()
	if mainPkgPath != "" {
		res = append(res, mainPkgPath)
	}

	// Binary dir
	binPath := filepath.Dir(selfPath)
	if binPath != "" {
		res = append(res, binPath)
	}
	return res
}
