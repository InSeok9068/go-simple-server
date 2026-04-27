package deariodate

import (
	"net/http"
	"strings"

	"simple-server/pkg/util/dateutil"

	"github.com/labstack/echo/v4"
)

// NormalizeWithDefault는 빈 날짜를 오늘 날짜로 보정한 뒤 YYYYMMDD 형식으로 검증한다.
func NormalizeWithDefault(value string) (string, error) {
	if value == "" {
		return dateutil.Today(), nil
	}
	return NormalizeRequired(value)
}

// NormalizeRequired는 날짜를 YYYYMMDD 형식으로 정규화하고 검증한다.
func NormalizeRequired(value string) (string, error) {
	normalized := strings.ReplaceAll(value, "-", "")
	if normalized == "" {
		return "", echo.NewHTTPError(http.StatusBadRequest, "날짜가 필요합니다.")
	}
	if !dateutil.IsValidDate(normalized, dateutil.DateFormatYYYYMMDD) {
		return "", echo.NewHTTPError(http.StatusBadRequest, "날짜 형식이 올바르지 않습니다.")
	}
	return normalized, nil
}
