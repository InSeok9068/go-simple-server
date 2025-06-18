package stringutil

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// IsEmpty 는 문자열이 비어있는지 확인합니다.
// 공백만 있는 문자열도 비어있다고 간주합니다.
//
// 예시:
//
//	if util.IsEmpty(userInput) {
//	    return errors.New("입력값이 비어 있습니다")
//	}
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 는 문자열이 비어있지 않은지 확인합니다.
// IsEmpty의 반대 동작을 수행합니다.
//
// 예시:
//
//	if util.IsNotEmpty(userInput) {
//	    processInput(userInput)
//	}
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// Truncate 는 문자열을 지정된 최대 길이로 자릅니다.
// 문자열이 최대 길이보다 길 경우 '...'을 추가합니다.
//
// 예시:
//
//	// "안녕하세요, 반갑..."을 반환합니다
//	shortText := util.Truncate("안녕하세요, 반갑습니다", 10)
func Truncate(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}

	runes := []rune(s)
	return string(runes[:maxLen]) + "..."
}

// TruncateWithSuffix 는 문자열을 지정된 최대 길이로 자르고 사용자 지정 접미사를 추가합니다.
//
// 예시:
//
//	// "안녕하세요, 반갑[더보기]"을 반환합니다
//	shortText := util.TruncateWithSuffix("안녕하세요, 반갑습니다", 10, "[더보기]")
func TruncateWithSuffix(s string, maxLen int, suffix string) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}

	runes := []rune(s)
	return string(runes[:maxLen]) + suffix
}

// Mask 는 문자열의 일부를 마스킹(*)합니다.
// start, end는 마스킹할 범위를 지정합니다. 음수 인덱스는 뒤에서부터 카운트합니다.
//
// 예시:
//
//	// "ab**fg"를 반환합니다
//	masked := util.Mask("abcdefg", 2, 4)
//
//	// "1234-56**-****-3456"을 반환합니다 (카드번호 마스킹)
//	maskedCardNumber := util.Mask("1234-5678-9012-3456", 7, -5)
func Mask(s string, start, end int) string {
	runes := []rune(s)
	length := len(runes)

	// 음수 인덱스 처리
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}

	// 유효한 범위 확인
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start > end || start >= length || end <= 0 {
		return s
	}

	result := make([]rune, length)
	copy(result, runes)

	for i := start; i < end; i++ {
		result[i] = '*'
	}

	return string(result)
}

// MaskEmail 은 이메일 주소를 마스킹합니다.
// ID의 일부와 도메인은 표시하고 나머지는 마스킹 처리합니다.
//
// 예시:
//
//	// "a***@example.com"을 반환합니다
//	masked := util.MaskEmail("admin@example.com")
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // 유효한 이메일이 아님
	}

	username := parts[0]
	domain := parts[1]

	// 사용자명 마스킹 처리
	if len(username) <= 1 {
		// 사용자명이 너무 짧으면 마스킹하지 않음
		return email
	}

	visible := int(math.Min(float64(len(username)/3+1), 3))
	maskedUsername := username[:visible] + strings.Repeat("*", len(username)-visible)

	return maskedUsername + "@" + domain
}

// MaskPhoneNumber 는 전화번호를 마스킹합니다.
// 전화번호의 중간 부분을 마스킹 처리합니다.
//
// 예시:
//
//	// "010-****-5678"을 반환합니다
//	masked := util.MaskPhoneNumber("010-1234-5678")
func MaskPhoneNumber(phoneNumber string) string {
	// 하이픈 제거
	digits := strings.ReplaceAll(phoneNumber, "-", "")
	length := len(digits)

	if length <= 4 {
		return phoneNumber // 너무 짧은 번호는 마스킹하지 않음
	}

	// 뒤 4자리를 제외한 나머지를 마스킹
	visible := 4
	if length < 7 { // 짧은 번호인 경우 뒤 3자리만 표시
		visible = 3
	}

	prefix := ""
	if length > 8 { // 휴대폰 번호인 경우 앞 3자리 표시
		prefix = digits[:3]
		digits = digits[3:]
		length -= 3
	}

	masked := strings.Repeat("*", length-visible) + digits[length-visible:]

	// 하이픈 처리
	if prefix != "" {
		return prefix + "-" + masked[:length-visible] + "-" + masked[length-visible:]
	}

	if length > 6 {
		return masked[:length-visible-4] + "-" + masked[length-visible-4:length-visible] + "-" + masked[length-visible:]
	}

	return masked
}

// RandomString 은 지정된 길이의 랜덤 문자열을 생성합니다.
// 영문 대소문자와 숫자로 구성됩니다.
//
// 예시:
//
//	// 16자리 랜덤 문자열 생성
//	token := util.RandomString(16)
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b)
}

// RandomHex 는 지정된 바이트 수만큼의 랜덤 16진수 문자열을 생성합니다.
//
// 예시:
//
//	// 16바이트(32자) 16진수 문자열 생성
//	hexString := util.RandomHex(16)
func RandomHex(bytes int) string {
	b := make([]byte, bytes)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

// RandomBase64 는 지정된 바이트 수만큼의 랜덤 Base64 문자열을 생성합니다.
//
// 예시:
//
//	// 24바이트 Base64 문자열 생성
//	b64String := util.RandomBase64(24)
func RandomBase64(bytes int) string {
	b := make([]byte, bytes)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(b)
}

// SanitizeFileName 은 파일명에서 사용할 수 없는 문자를 제거합니다.
//
// 예시:
//
//	// "my_document_2023"을 반환합니다
//	safeFileName := util.SanitizeFileName("my/document:2023")
func SanitizeFileName(filename string) string {
	// 윈도우와 유닉스 모두에서 사용할 수 없는 문자 제거
	sanitized := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`).ReplaceAllString(filename, "_")
	sanitized = strings.TrimSpace(sanitized)

	if sanitized == "" {
		return "untitled"
	}

	return sanitized
}

// ToSnakeCase 는 카멜 케이스 문자열을 스네이크 케이스로 변환합니다.
//
// 예시:
//
//	// "user_name"을 반환합니다
//	snake := util.ToSnakeCase("UserName")
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToCamelCase 는 스네이크 케이스 또는 케밥 케이스 문자열을 카멜 케이스로 변환합니다.
//
// 예시:
//
//	// "userName"을 반환합니다
//	camel := util.ToCamelCase("user_name")
//
//	// "productInfo"를 반환합니다
//	camel2 := util.ToCamelCase("product-info")
func ToCamelCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	result := strings.Builder{}
	result.WriteString(strings.ToLower(parts[0]))

	for _, part := range parts[1:] {
		if part == "" {
			continue
		}
		result.WriteString(strings.ToUpper(part[:1]) + strings.ToLower(part[1:]))
	}

	return result.String()
}

// ToPascalCase 는 스네이크 케이스 또는 케밥 케이스 문자열을 파스칼 케이스로 변환합니다.
//
// 예시:
//
//	// "UserName"을 반환합니다
//	pascal := util.ToPascalCase("user_name")
//
//	// "ProductInfo"를 반환합니다
//	pascal2 := util.ToPascalCase("product-info")
func ToPascalCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	result := strings.Builder{}

	for _, part := range parts {
		if part == "" {
			continue
		}
		result.WriteString(strings.ToUpper(part[:1]) + strings.ToLower(part[1:]))
	}

	return result.String()
}

// ToKebabCase 는 카멜 케이스 또는 파스칼 케이스 문자열을 케밥 케이스로 변환합니다.
//
// 예시:
//
//	// "user-name"을 반환합니다
//	kebab := util.ToKebabCase("UserName")
func ToKebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ReverseString 은 문자열을 뒤집습니다.
//
// 예시:
//
//	// "dlrow olleh"를 반환합니다
//	reversed := util.ReverseString("hello world")
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Slugify 는 문자열을 URL에 적합한 슬러그로 변환합니다.
// 공백은 하이픈으로 바꾸고, 영숫자와 하이픈 외의 모든 문자를 제거합니다.
//
// 예시:
//
//	// "hello-world"를 반환합니다
//	slug := util.Slugify("Hello, World!")
func Slugify(s string) string {
	// 모든 문자를 소문자로 변환
	s = strings.ToLower(s)

	// 문자 외의 문자를 공백으로 변환
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	s = reg.ReplaceAllString(s, "")

	// 연속된 공백을 하나의 공백으로 변환
	reg = regexp.MustCompile(`\s+`)
	s = reg.ReplaceAllString(s, " ")

	// 공백을 하이픈으로 변환
	s = strings.ReplaceAll(s, " ", "-")

	// 연속된 하이픈을 하나의 하이픈으로 변환
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	// 앞뒤 하이픈 제거
	s = strings.Trim(s, "-")

	return s
}

// OnlyDigits 는 문자열에서 숫자만 추출합니다.
//
// 예시:
//
//	// "1234"를 반환합니다
//	digits := util.OnlyDigits("abc123-4xyz")
func OnlyDigits(s string) string {
	reg := regexp.MustCompile(`[0-9]`)
	return strings.Join(reg.FindAllString(s, -1), "")
}

// OnlyLetters 는 문자열에서 알파벳만 추출합니다.
//
// 예시:
//
//	// "abcxyz"를 반환합니다
//	letters := util.OnlyLetters("abc123-4xyz")
func OnlyLetters(s string) string {
	reg := regexp.MustCompile(`[a-zA-Z]`)
	return strings.Join(reg.FindAllString(s, -1), "")
}

// OnlyAlphanumeric 는 문자열에서 영숫자만 추출합니다.
//
// 예시:
//
//	// "abc1234xyz"를 반환합니다
//	alphanum := util.OnlyAlphanumeric("abc123-4xyz")
func OnlyAlphanumeric(s string) string {
	reg := regexp.MustCompile(`[a-zA-Z0-9]`)
	return strings.Join(reg.FindAllString(s, -1), "")
}

// HasUppercase 는 문자열에 대문자가 포함되어 있는지 확인합니다.
//
// 예시:
//
//	// true를 반환합니다
//	hasUpper := util.HasUppercase("Hello")
func HasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// HasLowercase 는 문자열에 소문자가 포함되어 있는지 확인합니다.
//
// 예시:
//
//	// true를 반환합니다
//	hasLower := util.HasLowercase("Hello")
func HasLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// HasDigit 는 문자열에 숫자가 포함되어 있는지 확인합니다.
//
// 예시:
//
//	// true를 반환합니다
//	hasDigit := util.HasDigit("Hello123")
func HasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// HasSpecialChar 는 문자열에 특수 문자가 포함되어 있는지 확인합니다.
//
// 예시:
//
//	// true를 반환합니다
//	hasSpecial := util.HasSpecialChar("Hello!")
func HasSpecialChar(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

// IsValidEmail 은 문자열이 유효한 이메일 형식인지 확인합니다.
//
// 예시:
//
//	// true를 반환합니다
//	isValid := util.IsValidEmail("example@example.com")
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}
	return match
}

// IsValidKoreanPhoneNumber 는 문자열이 유효한 한국 전화번호 형식인지 확인합니다.
//
// 예시:
//
//	// true를 반환합니다
//	isValid := util.IsValidKoreanPhoneNumber("010-1234-5678")
func IsValidKoreanPhoneNumber(phone string) bool {
	// 하이픈 제거
	phone = strings.ReplaceAll(phone, "-", "")

	// 한국 휴대폰 번호 패턴 검사 (01X-XXXX-XXXX)
	pattern := `^01[0-9][0-9]{7,8}$`
	match, err := regexp.MatchString(pattern, phone)
	if err != nil {
		return false
	}
	return match
}

// FormatKoreanPhoneNumber 는 한국 전화번호에 하이픈을 추가하여 형식화합니다.
//
// 예시:
//
//	// "010-1234-5678"를 반환합니다
//	formatted := util.FormatKoreanPhoneNumber("01012345678")
func FormatKoreanPhoneNumber(phone string) string {
	// 기존 하이픈 제거
	phone = strings.ReplaceAll(phone, "-", "")

	if len(phone) < 9 || len(phone) > 11 {
		return phone // 유효하지 않은 길이
	}

	switch len(phone) {
	case 11:
		// 01X-XXXX-XXXX 형식
		return phone[:3] + "-" + phone[3:7] + "-" + phone[7:]
	case 10:
		// 01X-XXX-XXXX 형식
		return phone[:3] + "-" + phone[3:6] + "-" + phone[6:]
	default:
		// 02-XXXX-XXXX 형식 등
		return phone[:2] + "-" + phone[2:6] + "-" + phone[6:]
	}
}

// BytesToHumanReadable 은 바이트 크기를 사람이 읽기 쉬운 형식으로 변환합니다.
//
// 예시:
//
//	// "1.5 MB"를 반환합니다
//	size := util.BytesToHumanReadable(1500000)
func BytesToHumanReadable(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ToJSON 은 인터페이스를 JSON 문자열로 변환합니다.
//
// 예시:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	// {"name":"John","age":30}를 반환합니다
//	jsonStr := util.ToJSON(Person{Name: "John", Age: 30})
func ToJSON(v interface{}) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// ToPrettyJSON 은 인터페이스를 들여쓰기가 적용된 JSON 문자열로 변환합니다.
//
// 예시:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	// 들여쓰기가 적용된 JSON 문자열을 반환합니다
//	jsonStr := util.ToPrettyJSON(Person{Name: "John", Age: 30})
func ToPrettyJSON(v interface{}) string {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// FromJSON 은 JSON 문자열을 인터페이스로 변환합니다.
//
// 예시:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	var person Person
//	jsonStr := `{"name":"John","age":30}`
//	if util.FromJSON(jsonStr, &person) {
//	    // person 변수에 값이 설정됨
//	}
func FromJSON(jsonStr string, v interface{}) bool {
	err := json.Unmarshal([]byte(jsonStr), v)
	return err == nil
}

// ExtractNumbers 는 문자열에서 모든 숫자를 추출하여 int 슬라이스로 반환합니다.
//
// 예시:
//
//	// [123, 456]을 반환합니다
//	numbers := util.ExtractNumbers("abc123def456")
func ExtractNumbers(s string) []int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(s, -1)

	var numbers []int
	for _, match := range matches {
		num, err := strconv.Atoi(match)
		if err == nil {
			numbers = append(numbers, num)
		}
	}

	return numbers
}

// Capitalize 는 문자열의 첫 글자를 대문자로 변환합니다.
//
// 예시:
//
//	// "Hello world"를 반환합니다
//	capitalized := util.Capitalize("hello world")
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// CenterString 은 지정된 길이의, 지정된 문자로 패딩된 문자열 중앙에 텍스트를 배치합니다.
//
// 예시:
//
//	// "---Hello---"를 반환합니다
//	centered := util.CenterString("Hello", 13, '-')
func CenterString(s string, width int, padChar rune) string {
	if width <= len(s) {
		return s
	}

	padding := width - len(s)
	leftPad := padding / 2
	rightPad := padding - leftPad

	return strings.Repeat(string(padChar), leftPad) + s + strings.Repeat(string(padChar), rightPad)
}

// JoinNonEmpty 는 비어있지 않은 문자열들을 구분자로 결합합니다.
//
// 예시:
//
//	// "foo-bar-baz"를 반환합니다
//	joined := util.JoinNonEmpty([]string{"foo", "", "bar", "baz"}, "-")
func JoinNonEmpty(elements []string, separator string) string {
	var nonEmpty []string
	for _, element := range elements {
		if element != "" {
			nonEmpty = append(nonEmpty, element)
		}
	}
	return strings.Join(nonEmpty, separator)
}

// SubstringBefore 는 지정된 구분자 앞의 문자열 부분을 반환합니다.
//
// 예시:
//
//	// "Hello"를 반환합니다
//	before := util.SubstringBefore("Hello-World", "-")
func SubstringBefore(s, delimiter string) string {
	if delimiter == "" {
		return s
	}
	index := strings.Index(s, delimiter)
	if index == -1 {
		return s
	}
	return s[:index]
}

// SubstringAfter 는 지정된 구분자 뒤의 문자열 부분을 반환합니다.
//
// 예시:
//
//	// "World"를 반환합니다
//	after := util.SubstringAfter("Hello-World", "-")
func SubstringAfter(s, delimiter string) string {
	if delimiter == "" {
		return s
	}
	index := strings.Index(s, delimiter)
	if index == -1 {
		return ""
	}
	return s[index+len(delimiter):]
}

// SubstringBeforeLast 는 마지막 구분자 앞의 문자열 부분을 반환합니다.
//
// 예시:
//
//	// "Hello-beautiful"를 반환합니다
//	beforeLast := util.SubstringBeforeLast("Hello-beautiful-World", "-")
func SubstringBeforeLast(s, delimiter string) string {
	if delimiter == "" {
		return s
	}
	index := strings.LastIndex(s, delimiter)
	if index == -1 {
		return s
	}
	return s[:index]
}

// SubstringAfterLast 는 마지막 구분자 뒤의 문자열 부분을 반환합니다.
//
// 예시:
//
//	// "World"를 반환합니다
//	afterLast := util.SubstringAfterLast("Hello-beautiful-World", "-")
func SubstringAfterLast(s, delimiter string) string {
	if delimiter == "" {
		return s
	}
	index := strings.LastIndex(s, delimiter)
	if index == -1 {
		return ""
	}
	return s[index+len(delimiter):]
}

// CountWords 는 문자열에 포함된 단어의 수를 계산합니다.
//
// 예시:
//
//	// 3을 반환합니다
//	wordCount := util.CountWords("Hello beautiful world")
func CountWords(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	return len(strings.Fields(s))
}
