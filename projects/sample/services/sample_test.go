package services

import "testing"

func TestPlus(t *testing.T) {
	result := plus(1, 2)
	if result != 3 {
		t.Errorf("1 + 2 = %d, want 3", result)
	}
}
