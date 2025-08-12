package handlers

import (
	"strings"
	"testing"
)

func TestSplitRunes(t *testing.T) {
	text := strings.Repeat("a", 25)
	lines := splitRunes(text, 10)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != strings.Repeat("a", 10) {
		t.Fatalf("first line mismatch: %q", lines[0])
	}
	if lines[1] != strings.Repeat("a", 10) {
		t.Fatalf("second line mismatch: %q", lines[1])
	}
	if lines[2] != strings.Repeat("a", 5) {
		t.Fatalf("third line mismatch: %q", lines[2])
	}
}
