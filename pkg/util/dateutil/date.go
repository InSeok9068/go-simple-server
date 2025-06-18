package dateutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 상수 정의
const (
	// 표준 날짜 포맷
	DateFormatYYYYMMDD       = "20060102"
	DateFormatYYYYMMDDHHMMSS = "20060102150405"
	DateFormatISO            = "2006-01-02"
	DateFormatISOTime        = "2006-01-02 15:04:05"
	DateFormatRFC3339        = time.RFC3339

	// 한국식 날짜 포맷
	DateFormatKor     = "2006년 01월 02일"
	DateFormatKorTime = "2006년 01월 02일 15시 04분 05초"
)

// 한국어 요일 맵
var koreanWeekdays = map[time.Weekday]string{
	time.Sunday:    "일요일",
	time.Monday:    "월요일",
	time.Tuesday:   "화요일",
	time.Wednesday: "수요일",
	time.Thursday:  "목요일",
	time.Friday:    "금요일",
	time.Saturday:  "토요일",
}

// AddDaysToDate 는 YYYYMMDD 형식의 날짜 문자열에 일수를 더하고 같은 형식으로 반환합니다.
//
// 예시:
//
//	// "20230110" 반환
//	result, err := AddDaysToDate("20230105", 5)
func AddDaysToDate(dateStr string, days int) (string, error) {
	// 문자열을 time.Time으로 파싱
	t, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	// 날짜 계산
	newDate := t.AddDate(0, 0, days)

	// 다시 문자열로 포맷
	return newDate.Format(DateFormatYYYYMMDD), nil
}

// MustAddDaysToDate 는 에러를 반환하지 않는 AddDaysToDate의 변형입니다.
// 오류 발생 시 "00000000"을 반환합니다.
//
// 예시:
//
//	// "20230110" 반환
//	result := MustAddDaysToDate("20230105", 5)
func MustAddDaysToDate(dateStr string, days int) string {
	t, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "00000000" // 혹은 panic("날짜 포맷 오류")
	}

	newDate := t.AddDate(0, 0, days)
	return newDate.Format(DateFormatYYYYMMDD)
}

// MustFormatDateKor 는 YYYYMMDD 형식의 날짜 문자열을 한국어 형식으로 변환합니다.
// 에러 발생 시 빈 문자열을 반환합니다.
//
// 예시:
//
//	// "2023년 01월 05일" 반환
//	result := MustFormatDateKor("20230105")
func MustFormatDateKor(date string) string {
	t, err := time.Parse(DateFormatYYYYMMDD, date)
	if err != nil {
		return ""
	}
	return t.Format(DateFormatKor)
}

// MustFormatDateKorWithWeekDay 는 YYYYMMDD 형식의 날짜 문자열을 한국어 형식으로 변환하고 요일을 추가합니다.
// 에러 발생 시 빈 문자열을 반환합니다.
//
// 예시:
//
//	// "2023년 01월 05일 목요일" 반환
//	result := MustFormatDateKorWithWeekDay("20230105")
func MustFormatDateKorWithWeekDay(date string) string {
	t, err := time.Parse(DateFormatYYYYMMDD, date)
	if err != nil {
		return ""
	}

	korWeekday := koreanWeekdays[t.Weekday()]
	return t.Format(DateFormatKor) + " " + korWeekday
}

// AddMonthsToDate 는 YYYYMMDD 형식의 날짜 문자열에 개월 수를 더하고 같은 형식으로 반환합니다.
//
// 예시:
//
//	// "20230405" 반환
//	result, err := AddMonthsToDate("20230105", 3)
func AddMonthsToDate(dateStr string, months int) (string, error) {
	t, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	newDate := t.AddDate(0, months, 0)
	return newDate.Format(DateFormatYYYYMMDD), nil
}

// AddYearsToDate 는 YYYYMMDD 형식의 날짜 문자열에 연도를 더하고 같은 형식으로 반환합니다.
//
// 예시:
//
//	// "20240105" 반환
//	result, err := AddYearsToDate("20230105", 1)
func AddYearsToDate(dateStr string, years int) (string, error) {
	t, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	newDate := t.AddDate(years, 0, 0)
	return newDate.Format(DateFormatYYYYMMDD), nil
}

// FormatDate 는 YYYYMMDD 형식의 날짜 문자열을 지정된 포맷으로 변환합니다.
//
// 예시:
//
//	// "2023-01-05" 반환
//	result, err := FormatDate("20230105", DateFormatISO)
func FormatDate(dateStr string, format string) (string, error) {
	t, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	return t.Format(format), nil
}

// ParseDate 는 지정된 포맷의 날짜 문자열을 time.Time으로 파싱합니다.
//
// 예시:
//
//	// time.Time 객체 반환
//	t, err := ParseDate("2023-01-05", DateFormatISO)
func ParseDate(dateStr string, format string) (time.Time, error) {
	return time.Parse(format, dateStr)
}

// Now 는 현재 시간을 지정된 포맷의 문자열로 반환합니다.
//
// 예시:
//
//	// 현재 날짜를 "20230105" 형식으로 반환
//	today := Now(DateFormatYYYYMMDD)
func Now(format string) string {
	return time.Now().Format(format)
}

// Today 는 현재 날짜를 YYYYMMDD 형식으로 반환합니다.
//
// 예시:
//
//	// 현재 날짜를 "20230105" 형식으로 반환
//	today := Today()
func Today() string {
	return time.Now().Format(DateFormatYYYYMMDD)
}

// DaysBetween 은 두 날짜 사이의、일수 차이를 계산합니다.
// 결과는 end - start 이므로 end가 start보다 이후 날짜면 양수, 이전이면 음수입니다.
//
// 예시:
//
//	// 5 반환
//	days, err := DaysBetween("20230105", "20230110")
func DaysBetween(startDateStr, endDateStr string) (int, error) {
	startDate, err := time.Parse(DateFormatYYYYMMDD, startDateStr)
	if err != nil {
		return 0, fmt.Errorf("시작 날짜 파싱 오류: %w", err)
	}

	endDate, err := time.Parse(DateFormatYYYYMMDD, endDateStr)
	if err != nil {
		return 0, fmt.Errorf("종료 날짜 파싱 오류: %w", err)
	}

	// UTC로 변환하여 시간대 차이 제거
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)

	diff := endDate.Sub(startDate)
	return int(diff.Hours() / 24), nil
}

// IsWeekend 는 지정된 날짜가 주말(토요일 또는 일요일)인지 확인합니다.
//
// 예시:
//
//	// true/false 반환
//	isWeekend, err := IsWeekend("20230107") // 토요일인 경우 true
func IsWeekend(dateStr string) (bool, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return false, fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday, nil
}

// GetWeekday 는 지정된 날짜의 요일을 반환합니다.
//
// 예시:
//
//	// time.Weekday 반환
//	weekday, err := GetWeekday("20230105")
func GetWeekday(dateStr string) (time.Weekday, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return 0, fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	return date.Weekday(), nil
}

// GetKoreanWeekday 는 지정된 날짜의 한국어 요일을 반환합니다.
//
// 예시:
//
//	// "목요일" 반환
//	korWeekday, err := GetKoreanWeekday("20230105")
func GetKoreanWeekday(dateStr string) (string, error) {
	weekday, err := GetWeekday(dateStr)
	if err != nil {
		return "", err
	}

	return koreanWeekdays[weekday], nil
}

// GetFirstDayOfMonth 는 지정된 날짜가 속한 월의 첫째 날을 반환합니다.
//
// 예시:
//
//	// "20230101" 반환
//	firstDay, err := GetFirstDayOfMonth("20230115")
func GetFirstDayOfMonth(dateStr string) (string, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	firstDay := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	return firstDay.Format(DateFormatYYYYMMDD), nil
}

// GetLastDayOfMonth 는 지정된 날짜가 속한 월의 마지막 날을 반환합니다.
//
// 예시:
//
//	// "20230131" 반환
//	lastDay, err := GetLastDayOfMonth("20230115")
func GetLastDayOfMonth(dateStr string) (string, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	// 다음 달의 첫째 날에서 하루를 빼면 이번 달의 마지막 날
	nextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	lastDay := nextMonth.AddDate(0, 0, -1)

	return lastDay.Format(DateFormatYYYYMMDD), nil
}

// IsValidDate 는 문자열이 유효한 날짜 형식인지 확인합니다.
//
// 예시:
//
//	// 유효한 날짜면 true 반환
//	isValid := IsValidDate("20230105", DateFormatYYYYMMDD)
func IsValidDate(dateStr string, format string) bool {
	_, err := time.Parse(format, dateStr)
	return err == nil
}

// FormatTimeAgo 는 주어진 시간과 현재 시간 사이의 경과 시간을 사람이 읽기 쉬운 형식으로 변환합니다.
//
// 예시:
//
//	// "3분 전", "1시간 전", "2일 전" 등 반환
//	timeAgo := FormatTimeAgo(pastTime)
func FormatTimeAgo(past time.Time) string {
	now := time.Now()
	diff := now.Sub(past)

	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := hours / 24
	months := days / 30
	years := days / 365

	switch {
	case seconds < 60:
		return fmt.Sprintf("%d초 전", seconds)
	case minutes < 60:
		return fmt.Sprintf("%d분 전", minutes)
	case hours < 24:
		return fmt.Sprintf("%d시간 전", hours)
	case days < 30:
		return fmt.Sprintf("%d일 전", days)
	case months < 12:
		return fmt.Sprintf("%d개월 전", months)
	default:
		return fmt.Sprintf("%d년 전", years)
	}
}

// GetAge 는 생년월일로부터 현재 나이를 계산합니다.
//
// 예시:
//
//	// 나이 반환
//	age, err := GetAge("19900105")
func GetAge(birthDateStr string) (int, error) {
	birthDate, err := time.Parse(DateFormatYYYYMMDD, birthDateStr)
	if err != nil {
		return 0, fmt.Errorf("생년월일 파싱 오류: %w", err)
	}

	now := time.Now()

	// 나이 계산
	age := now.Year() - birthDate.Year()

	// 아직 생일이 지나지 않았으면 1을 뺍니다
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age, nil
}

// GetQuarter 는 지정된 날짜가 속한 분기를 반환합니다.
//
// 예시:
//
//	// 1~4 사이의 분기 번호 반환
//	quarter, err := GetQuarter("20230215") // 1 반환 (1분기)
func GetQuarter(dateStr string) (int, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return 0, fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	month := int(date.Month())
	quarter := (month-1)/3 + 1

	return quarter, nil
}

// GetFirstDayOfQuarter 는 지정된 날짜가 속한 분기의 첫째 날을 반환합니다.
//
// 예시:
//
//	// "20230101" 반환 (1분기의 첫날)
//	firstDay, err := GetFirstDayOfQuarter("20230215")
func GetFirstDayOfQuarter(dateStr string) (string, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	quarter := (int(date.Month())-1)/3 + 1
	firstMonth := (quarter-1)*3 + 1

	firstDay := time.Date(date.Year(), time.Month(firstMonth), 1, 0, 0, 0, 0, date.Location())
	return firstDay.Format(DateFormatYYYYMMDD), nil
}

// GetLastDayOfQuarter 는 지정된 날짜가 속한 분기의 마지막 날을 반환합니다.
//
// 예시:
//
//	// "20230331" 반환 (1분기의 마지막날)
//	lastDay, err := GetLastDayOfQuarter("20230215")
func GetLastDayOfQuarter(dateStr string) (string, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	quarter := (int(date.Month())-1)/3 + 1
	lastMonth := quarter * 3

	// 다음 달의 첫째 날에서 하루를 빼면 이번 달의 마지막 날
	nextMonth := time.Date(date.Year(), time.Month(lastMonth)+1, 1, 0, 0, 0, 0, date.Location())
	lastDay := nextMonth.AddDate(0, 0, -1)

	return lastDay.Format(DateFormatYYYYMMDD), nil
}

// IsBetween 은 대상 날짜가 시작일과 종료일 사이에 있는지 확인합니다.
//
// 예시:
//
//	// 날짜가 범위 내에 있으면 true 반환
//	isBetween, err := IsBetween("20230110", "20230105", "20230115")
func IsBetween(dateStr, startDateStr, endDateStr string) (bool, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return false, fmt.Errorf("대상 날짜 파싱 오류: %w", err)
	}

	startDate, err := time.Parse(DateFormatYYYYMMDD, startDateStr)
	if err != nil {
		return false, fmt.Errorf("시작 날짜 파싱 오류: %w", err)
	}

	endDate, err := time.Parse(DateFormatYYYYMMDD, endDateStr)
	if err != nil {
		return false, fmt.Errorf("종료 날짜 파싱 오류: %w", err)
	}

	return !date.Before(startDate) && !date.After(endDate), nil
}

// GetDayOfYear 는 지정된 날짜가 해당 연도의 몇 번째 날인지 반환합니다.
//
// 예시:
//
//	// 5 반환 (1월 5일은 연중 5번째 날)
//	dayOfYear, err := GetDayOfYear("20230105")
func GetDayOfYear(dateStr string) (int, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return 0, fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	startOfYear := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, date.Location())
	diff := date.Sub(startOfYear)

	return int(diff.Hours()/24) + 1, nil
}

// GetWeekNumber 는 지정된 날짜의 ISO 주차를 반환합니다.
//
// 예시:
//
//	// 1 반환 (2023년 1월 5일은 1주차)
//	week, err := GetWeekNumber("20230105")
func GetWeekNumber(dateStr string) (int, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return 0, fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	_, week := date.ISOWeek()
	return week, nil
}

// GetFirstDayOfWeek 는 지정된 날짜가 속한 주의 첫째 날(월요일)을 반환합니다.
//
// 예시:
//
//	// "20230102" 반환 (2023년 1월 5일이 속한 주의 월요일)
//	monday, err := GetFirstDayOfWeek("20230105")
func GetFirstDayOfWeek(dateStr string) (string, error) {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}

	// 지정된 날짜의 요일
	weekday := int(date.Weekday())

	// 일요일(0)을 주의 7번째 날로 처리
	if weekday == 0 {
		weekday = 7
	}

	// 월요일(1)까지의 차이를 계산하여 빼기
	daysToSubtract := weekday - 1
	monday := date.AddDate(0, 0, -daysToSubtract)

	return monday.Format(DateFormatYYYYMMDD), nil
}

// GetLastDayOfWeek 는 지정된 날짜가 속한 주의 마지막 날(일요일)을 반환합니다.
//
// 예시:
//
//	// "20230108" 반환 (2023년 1월 5일이 속한 주의 일요일)
//	sunday, err := GetLastDayOfWeek("20230105")
func GetLastDayOfWeek(dateStr string) (string, error) {
	firstDay, err := GetFirstDayOfWeek(dateStr)
	if err != nil {
		return "", err
	}

	firstDate, err := time.Parse(DateFormatYYYYMMDD, firstDay)
	if err != nil {
		return "", fmt.Errorf("날짜 파싱 오류: %w", err)
	}
	sunday := firstDate.AddDate(0, 0, 6) // 월요일부터 6일 후는 일요일

	return sunday.Format(DateFormatYYYYMMDD), nil
}

// ParseDurationString 은 "1h30m", "3d", "2w" 등의 문자열을 time.Duration으로 파싱합니다.
// 지원하는 단위: s(초), m(분), h(시간), d(일), w(주)
//
// 예시:
//
//	// 1.5시간에 해당하는 time.Duration 반환
//	duration, err := ParseDurationString("1h30m")
//
//	// 3일에 해당하는 time.Duration 반환
//	duration, err := ParseDurationString("3d")
func ParseDurationString(durationStr string) (time.Duration, error) {
	// 일(d)과 주(w) 처리
	if strings.Contains(durationStr, "d") || strings.Contains(durationStr, "w") {
		var totalDuration time.Duration

		// 문자열을 숫자와 단위로 분리
		re := regexp.MustCompile(`(\d+)([a-z])`)
		matches := re.FindAllStringSubmatch(durationStr, -1)

		for _, match := range matches {
			if len(match) != 3 {
				continue
			}

			value, err := strconv.Atoi(match[1])
			if err != nil {
				return 0, fmt.Errorf("유효하지 않은 지속 시간: %s", durationStr)
			}

			unit := match[2]

			switch unit {
			case "d":
				totalDuration += time.Duration(value) * 24 * time.Hour
			case "w":
				totalDuration += time.Duration(value) * 7 * 24 * time.Hour
			default:
				// 표준 단위는 그대로 유지
				standardDuration, err := time.ParseDuration(match[1] + unit)
				if err != nil {
					return 0, fmt.Errorf("유효하지 않은 지속 시간 단위: %s", unit)
				}
				totalDuration += standardDuration
			}
		}

		return totalDuration, nil
	}

	// 표준 단위(s, m, h)인 경우 time.ParseDuration 사용
	return time.ParseDuration(durationStr)
}

// ConvertTimezone 은 시간을 다른 시간대로 변환합니다.
//
// 예시:
//
//	// UTC 시간을 아시아/서울 시간대로 변환
//	seoulTime, err := ConvertTimezone(utcTime, "UTC", "Asia/Seoul")
func ConvertTimezone(t time.Time, fromTZ, toTZ string) (time.Time, error) {
	// 원본 시간대 로드
	fromLoc, err := time.LoadLocation(fromTZ)
	if err != nil {
		return time.Time{}, fmt.Errorf("원본 시간대 로드 오류: %w", err)
	}

	// 대상 시간대 로드
	toLoc, err := time.LoadLocation(toTZ)
	if err != nil {
		return time.Time{}, fmt.Errorf("대상 시간대 로드 오류: %w", err)
	}

	// 원본 시간대로 설정
	inFromTZ := t.In(fromLoc)

	// 대상 시간대로 변환
	inToTZ := inFromTZ.In(toLoc)

	return inToTZ, nil
}

// IsLeapYear 는 지정된 연도가 윤년인지 확인합니다.
//
// 예시:
//
//	// true 반환 (2024년은 윤년)
//	isLeap := IsLeapYear(2024)
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// GetDaysInMonth 는 지정된 연도와 월의 일수를 반환합니다.
//
// 예시:
//
//	// 31 반환 (1월은 31일)
//	days := GetDaysInMonth(2023, 1)
func GetDaysInMonth(year, month int) int {
	// 각 월별 일수
	daysInMonth := []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// 2월이고 윤년이면 29일
	if month == 2 && IsLeapYear(year) {
		return 29
	}

	return daysInMonth[month]
}

// GetDaysInYear 는 지정된 연도의 일수를 반환합니다.
//
// 예시:
//
//	// 365 또는 366 반환
//	days := GetDaysInYear(2023) // 365 반환
func GetDaysInYear(year int) int {
	if IsLeapYear(year) {
		return 366
	}
	return 365
}
