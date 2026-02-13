package aiclient

import (
	"context"
	"testing"
)

func TestTranscribeAudio(t *testing.T) {
	t.Run("empty data", func(t *testing.T) {
		if _, err := TranscribeAudio(context.Background(), nil, "audio/webm"); err == nil {
			t.Fatalf("error expected")
		}
	})

	t.Run("missing api key", func(t *testing.T) {
		t.Setenv("GEMINI_AI_KEY", "")
		if _, err := TranscribeAudio(context.Background(), []byte{1}, "audio/webm"); err == nil {
			t.Fatalf("error expected")
		}
	})
}
