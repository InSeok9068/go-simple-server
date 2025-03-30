package util

import (
	"fmt"
	"time"
)

func AddDaysToDate(dateStr string, days int) (string, error) {
	// 문자열을 time.Time으로 파싱
	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	// 날짜 계산
	newDate := t.AddDate(0, 0, days)

	// 다시 문자열로 포맷
	return newDate.Format("20060102"), nil
}

func MustAddDaysToDate(dateStr string, days int) string {
	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		// 실패 시 기본값 또는 panic
		return "00000000" // 혹은 panic("날짜 포맷 오류")
	}

	newDate := t.AddDate(0, 0, days)
	return newDate.Format("20060102")
}
