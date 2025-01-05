package internal

import (
	resources "simple-server"

	"github.com/joho/godotenv"
)

var EnvMap map[string]string

func LoadEnv() {
	envData, _ := resources.EmbeddedFiles.ReadFile(".env")
	EnvMap, _ = godotenv.Unmarshal(string(envData))
}
