package gomutil

import (
	"fmt"
	"testing"
)

func TestCSRFInput(t *testing.T) {
	result := fmt.Sprint(CSRFInput("token"))
	expected := `<input type="hidden" name="csrf" value="token">`
	if result != expected {
		t.Errorf("got %s, want %s", result, expected)
	}
}
