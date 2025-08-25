package config

import (
	"fmt"
	"strings"
)

func AppDatabaseURL(serviceName string) string {
	appPragmas := []string{
		"_pragma=journal_mode(WAL)",            // 동시성 향상
		"_pragma=synchronous(NORMAL)",          // 정합성 향상
		"_pragma=busy_timeout(5000)",           // 쓰기 경합 대기 5초
		"_pragma=foreign_keys(ON)",             // 외래 키 허용
		"_pragma=temp_store(MEMORY)",           // 임시 저장소 메모리
		"_pragma=journal_size_limit(67108864)", // 64MB
	}
	pragmas := strings.Join(appPragmas, "&")
	var url string
	if IsDevEnv() {
		url = fmt.Sprintf(`file:./projects/%s/data/data.db?%s`, serviceName, pragmas)
	} else {
		url = fmt.Sprintf(`file:/srv/%s/data/data.db?%s`, serviceName, pragmas)
	}
	return url
}

func LogDatabaseURL() string {
	logPragmas := []string{
		"_pragma=journal_mode(WAL)",            // 동시성 향상
		"_pragma=synchronous(NORMAL)",          // 정합성 향상
		"_pragma=busy_timeout(5000)",           // 쓰기 경합 대기 5초
		"_pragma=temp_store(MEMORY)",           // 임시 저장소 메모리
		"_pragma=journal_size_limit(67108864)", // 64MB
	}
	pragmas := strings.Join(logPragmas, "&")
	var url string
	if IsDevEnv() {
		url = fmt.Sprintf(`file:./shared/log/auxiliary.db?%s`, pragmas)
	} else {
		url = fmt.Sprintf(`file:/srv/log/auxiliary.db?%s`, pragmas)
	}
	return url
}
