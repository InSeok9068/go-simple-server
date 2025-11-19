package views

import "time"

func DateView(date string) string {
	parsed, _ := time.Parse("20060102", date)
	return parsed.Format("1월 2일")
}
