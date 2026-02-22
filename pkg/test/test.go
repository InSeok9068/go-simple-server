package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
)

// Request ✅ 공통 HTTP 요청 실행 함수
func Request(
	t *testing.T,
	method string,
	serverURL string,
	path string,
	queryParams map[string]string,
	body []byte,
) (int, string) {
	t.Helper()

	reqURL, err := buildRequestURL(serverURL, path, queryParams)
	if err != nil {
		t.Fatalf("요청 URL 생성 실패: %v", err)
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

	if err := validateRequestTarget(serverURL, reqURL); err != nil {
		t.Fatalf("요청 대상 검증 실패: %v", err)
	}

	// ✅ 요청 실행
	client := &http.Client{}
	// #nosec G704 -- 테스트 서버 URL(host/scheme 검증 완료)로만 요청합니다.
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

func buildRequestURL(serverURL, path string, queryParams map[string]string) (string, error) {
	reqURL, err := url.Parse(fmt.Sprintf("%s%s", serverURL, path))
	if err != nil {
		return "", err
	}

	query := reqURL.Query()
	for key, value := range queryParams {
		query.Set(key, value)
	}
	reqURL.RawQuery = query.Encode()
	return reqURL.String(), nil
}

func validateRequestTarget(serverURL, reqURL string) error {
	baseURL, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("서버 URL 파싱 실패: %w", err)
	}
	targetURL, err := url.Parse(reqURL)
	if err != nil {
		return fmt.Errorf("요청 URL 파싱 실패: %w", err)
	}
	if targetURL.Scheme != baseURL.Scheme || targetURL.Host != baseURL.Host {
		return fmt.Errorf("허용되지 않은 요청 대상입니다: %s", reqURL)
	}
	return nil
}
