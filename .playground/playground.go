package main

import (
	"fmt"
	"strings"
)

func main() {
	appPragmas := []string{
		"_pragma=journal_mode(WAL)",            // 동시성 향상
		"_pragma=synchronous(NORMAL)",          // 정합성 향상
		"_pragma=busy_timeout(5000)",           // 쓰기 경합 대기 5초
		"_pragma=foreign_keys(ON)",             // 외래 키 허용
		"_pragma=temp_store(MEMORY)",           // 임시 저장소 메모리
		"_pragma=journal_size_limit(67108864)", // 64MB
	}
	pragmas := strings.Join(appPragmas, "&")
	fmt.Println(pragmas)
}
