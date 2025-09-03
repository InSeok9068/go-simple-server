package aiclient

import (
	"context"
	"testing"

	"simple-server/internal/config"
)

func TestTranscribeAudio(t *testing.T) {
	config.EnvMap = map[string]string{}
	t.Run("empty data", func(t *testing.T) {
		if _, err := TranscribeAudio(context.Background(), nil, "audio/webm"); err == nil {
			t.Fatalf("error expected")
		}
	})

	t.Run("missing api key", func(t *testing.T) {
		if _, err := TranscribeAudio(context.Background(), []byte{1}, "audio/webm"); err == nil {
			t.Fatalf("error expected")
		}
	})
}
