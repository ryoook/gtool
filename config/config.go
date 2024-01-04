package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/ryoook/gtool/env"
	"github.com/ryoook/gtool/rootpath"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
)

// Get Read your custom config into your structure
func Get(path string, retConfig interface{}) error {
	rawConfig, err := get(path)
	if err != nil {
		return err
	}

	return mapstructure.Decode(rawConfig, retConfig)
}

var (
	container *sync.Map
)

func init() {
	container = &sync.Map{}
}

func get(path string) (interface{}, error) {
	var (
		value interface{}
		err   error
		ok    bool
	)
	value, ok = container.Load(path)
	if !ok {
		// init
		value, err = loadYaml(path)
		if err != nil {
			return nil, err
		}
	}
	container.Store(path, value)
	return value, err
}

func loadYaml(path string) (interface{}, error) {
	content, err := read(path)
	if err != nil {
		return nil, err
	}

	return parseYamlContent(content)
}

func read(path string) ([]byte, error) {
	rootPath := rootpath.RootPath()
	if len(rootPath) == 0 {
		return nil, fmt.Errorf("root path is empty")
	}

	var (
		filePath string
		err      error
	)
	for _, root := range rootPath {
		filePath = filepath.Join(root, fmt.Sprintf("%s.yml", path))
		if _, err = os.Stat(filePath); err != nil {
			continue
		} else {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("file not found, path: %s", path)
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func parseYamlContent(content []byte) (interface{}, error) {
	var config map[string]interface{}
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	if c, exist := config[env.GetEnv()]; exist {
		return c, nil
	}
	return nil, fmt.Errorf("env: %s is not configed", env.GetEnv())
}
