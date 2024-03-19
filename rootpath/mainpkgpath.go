package rootpath

import (
	"github.com/ryoook/gtool/tool"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync/atomic"
)

var (
	mainPkgDir   atomic.Value
	mainPkgRegex atomic.Value
)

func getMainPkgPath() string {
	if mainPkgDir.Load() != nil {
		return mainPkgDir.Load().(string)
	}
	goRootForFilter := runtime.GOROOT()
	var lastFile string
	for i := 1; i < 10000; i++ {
		if pc, fPath, _, ok := runtime.Caller(i); ok {
			if goRootForFilter != "" && len(fPath) >= len(goRootForFilter) && fPath[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			if filepath.Ext(fPath) != ".go" {
				continue
			}
			lastFile = fPath
			if fn := runtime.FuncForPC(pc); fn != nil {
				array := strings.Split(fn.Name(), ".")
				if array[0] != "main" {
					continue
				}
			}
			if getMainPkgRegex().Match(getFileContent(fPath)) {
				mainPkgDir.Store(filepath.Dir(fPath))
				return filepath.Dir(fPath)
			}
		} else {
			break
		}
	}
	if lastFile != "" {
		for path := filepath.Dir(lastFile); len(path) > 1 && path[len(path)-1] != os.PathSeparator; {
			for _, v := range tool.ScanDir(path, "*.go") {
				if getMainPkgRegex().Match(getFileContent(v)) {
					mainPkgDir.Store(filepath.Dir(v))
					return filepath.Dir(v)
				}
			}
			path = filepath.Dir(path)
		}
	}
	return ""
}

func getFileContent(path string) []byte {
	data, _ := os.ReadFile(path)
	return data
}

func getMainPkgRegex() *regexp.Regexp {
	regexIns, ok := mainPkgRegex.Load().(*regexp.Regexp)
	if !ok {
		regexIns = regexp.MustCompile(`package\s+main\s+`)
		mainPkgRegex.Store(regexIns)
	}
	return regexIns
}
