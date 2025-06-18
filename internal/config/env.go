package config

import (
	"log/slog"
	resources "simple-server"

	"github.com/joho/godotenv"
)

var EnvMap map[string]string

func LoadEnv() {
	envData, err := resources.EmbeddedFiles.ReadFile(".env")
	if err != nil {
		slog.Error("환경 파일 읽기 실패", "error", err)
		return
	}

	envMap, err := godotenv.Unmarshal(string(envData))
	if err != nil {
		slog.Error("환경 변수 파싱 실패", "error", err)
		return
	}

	EnvMap = envMap
}

func IsDevEnv() bool {
	return EnvMap["ENV"] == "dev"
}

func IsProdEnv() bool {
	return EnvMap["ENV"] == "prod"
}
