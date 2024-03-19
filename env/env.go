package env

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

const (
	Develop = "develop"
)

var currentEnv string

func init() {
	currentEnv = os.Getenv("ENV")
	if currentEnv == "" {
		currentEnv = Develop
	}
}

// GetEnv Get current running environment
func GetEnv() string {
	return currentEnv
}

// SetEnv ...
func SetEnv(newEnv string) {
	currentEnv = newEnv
}
