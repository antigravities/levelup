package util

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables
func LoadEnv() {
	env := os.Getenv("LU_ENV")

	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")

	if "test" != env {
		godotenv.Load(".env.local")
	}

	godotenv.Load(".env." + env)
	godotenv.Load()
}
