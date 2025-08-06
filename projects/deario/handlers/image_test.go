package handlers

import "testing"

func TestParseFirebaseURL(t *testing.T) {
	raw := "https://firebasestorage.googleapis.com/v0/b/test-bucket/o/diary%2Fuid%2Fimg.png?alt=media&token=abc"
	bucket, object, err := parseFirebaseURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bucket != "test-bucket" {
		t.Fatalf("bucket mismatch: %s", bucket)
	}
	if object != "diary/uid/img.png" {
		t.Fatalf("object mismatch: %s", object)
	}
}

func TestParseFirebaseURL_Invalid(t *testing.T) {
	if _, _, err := parseFirebaseURL("https://example.com"); err == nil {
		t.Fatalf("expected error")
	}
}
