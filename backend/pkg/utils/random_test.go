package utils

import (
	"testing"
	"unicode"
)

func TestGenerateRandomNumber_Length(t *testing.T) {
	result := GenerateRandomNumber(6)
	if len(result) != 6 {
		t.Errorf("expected length 6, got %d", len(result))
	}
}

func TestGenerateRandomNumber_VariousLengths(t *testing.T) {
	for _, n := range []int{1, 4, 10, 20} {
		result := GenerateRandomNumber(n)
		if len(result) != n {
			t.Errorf("expected length %d, got %d", n, len(result))
		}
	}
}

func TestGenerateRandomNumber_OnlyDigits(t *testing.T) {
	result := GenerateRandomNumber(100)
	for _, r := range result {
		if !unicode.IsDigit(r) {
			t.Errorf("expected only digits, got %c", r)
		}
	}
}
