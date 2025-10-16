package config

import (
	"log/slog"
	"os" 

	"github.com/joho/godotenv"
)

var EnvMap map[string]string

func LoadEnv() {
	// .env 파일 읽기
	envMap, err := godotenv.Read(".env")
	if err != nil {
		slog.Error("환경 변수 파싱 실패", "error", err)
		return
	}

	// 환경 변수로 설정
	for key, value := range envMap {
		if os.Getenv(key) == "" { // 이미 설정된 환경 변수가 없을 때만 설정
			if err := os.Setenv(key, value); err != nil {
				slog.Error("환경 변수 설정 실패", "key", key, "error", err)
				continue
			}
		}
	}

	// 전역 변수에도 저장 (필요한 경우)
	EnvMap = envMap
}

func IsDevEnv() bool {
	return EnvMap["ENV"] == "dev"
}

func IsProdEnv() bool {
	return EnvMap["ENV"] == "prod"
}
