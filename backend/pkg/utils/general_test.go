package utils

import (
	"testing"
)

func TestFormatNumber_Large(t *testing.T) {
	result := FormatNumber(1500000)
	if result != "1.500.000" {
		t.Errorf("expected '1.500.000', got '%s'", result)
	}
}

func TestFormatNumber_Small(t *testing.T) {
	result := FormatNumber(42)
	if result != "42" {
		t.Errorf("expected '42', got '%s'", result)
	}
}

func TestFormatNumber_Zero(t *testing.T) {
	result := FormatNumber(0)
	if result != "0" {
		t.Errorf("expected '0', got '%s'", result)
	}
}

func TestFormatNumber_Thousands(t *testing.T) {
	result := FormatNumber(1234)
	if result != "1.234" {
		t.Errorf("expected '1.234', got '%s'", result)
	}
}
