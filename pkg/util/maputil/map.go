package maputil

import "strconv"

// GetString은 맵에서 문자열 값을 안전하게 추출합니다.
func GetString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetInt64은 맵에서 int64 값을 안전하게 추출합니다.
func GetInt64(m map[string]interface{}, key string, defaultValue int64) int64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return int64(v)
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				return i
			}
		}
	}
	return defaultValue
}
