package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Load environment variables
	env := os.Getenv("ENV")
	if env != "" {
		fmt.Printf("Loading ENV variables from: .env.%v\n", env)
	}

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	// First load the environment specific file, it takes precedence over the original .env file
	godotenv.Load(fmt.Sprintf("%v/.env.%v", currentDir, env))
	// Load the original .env file
	godotenv.Load(fmt.Sprintf("%v/.env", currentDir))
}
