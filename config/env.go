package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment variables from the .env-[ENV] or .env file
// and returns a kafka.ConfigMap with the client configuration
func LoadEnv() {
	// Load environment variables
	env := os.Getenv("ENV")

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	configFilePath := ""
	if env != "" {
		configFilePath = fmt.Sprintf("%v/config-%v.properties", currentDir, env)
	} else {
		configFilePath = fmt.Sprintf("%v/config.properties", currentDir)
	}
	fmt.Printf("Loading ENV variables from: %v\n", configFilePath)
	// Loading Topic mainly
	godotenv.Load(configFilePath)
}
