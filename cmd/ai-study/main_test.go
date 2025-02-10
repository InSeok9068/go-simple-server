package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"simple-server/internal"
	"simple-server/pkg/test"
	"testing"
)

func runTestServer() *httptest.Server {
	/* 환경 설정 */
	internal.LoadEnv()
	os.Setenv("SERVICE_NAME", "ai-study")
	os.Setenv("APP_TITLE", "🕵️‍♀️ AI 공부 길잡이")
	/* 환경 설정 */

	return httptest.NewServer(setUpServer())
}

func TestAIStudy(t *testing.T) {
	server := runTestServer()
	defer server.Close()

	t.Run("AI Study API 호출", func(t *testing.T) {
		input := "Go 언어"

		code, body := test.TestRequest(t, http.MethodPost, server.URL, "/ai-study", map[string]string{"input": input}, nil)

		// ✅ 응답 출력
		if code != http.StatusOK {
			t.Fatalf("예상한 응답 코드가 아닙니다. 예상: %d, 결과: %d", http.StatusOK, code)
		}

		t.Log("응답 바디:", body)
	})

	t.Run("AI Study API 호출 랜덤", func(t *testing.T) {
		code, body := test.TestRequest(t, http.MethodPost, server.URL, "/ai-study-random", nil, nil)

		// ✅ 응답 출력
		if code != http.StatusOK {
			t.Fatalf("예상한 응답 코드가 아닙니다. 예상: %d, 결과: %d", http.StatusOK, code)
		}

		t.Log("응답 바디:", body)
	})
}
