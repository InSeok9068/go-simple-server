package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
)

// ✅ 공통 HTTP 요청 실행 함수
func TestRequest(
	t *testing.T,
	method string,
	serverURL string,
	path string,
	queryParams map[string]string,
	body []byte,
) (int, string) {
	t.Helper()

	// ✅ Query Parameters 처리
	reqURL := fmt.Sprintf("%s%s", serverURL, path)
	if len(queryParams) > 0 {
		query := url.Values{}
		for key, value := range queryParams {
			query.Add(key, value)
		}
		reqURL = fmt.Sprintf("%s?%s", reqURL, query.Encode())
	}

	// ✅ 요청 생성
	req, err := http.NewRequestWithContext(
		t.Context(),
		method,
		reqURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("요청 생성 실패: %v", err)
	}

	// ✅ Content-Type 설정 (JSON)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// ✅ 요청 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// ✅ 응답 바디 읽기
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("응답 바디 읽기 실패: %v", err)
	}
	bodyString := string(bodyBytes)

	return resp.StatusCode, bodyString
}
