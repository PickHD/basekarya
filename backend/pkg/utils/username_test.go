package utils

import (
	"strings"
	"testing"
)

func TestGenerateUsername_UnderscoreSeparator(t *testing.T) {
	result := GenerateUsername("John Doe")
	parts := strings.SplitN(result, "_", 2)
	if len(parts) != 2 {
		t.Errorf("expected underscore separator, got %s", result)
	}
	if len(parts[1]) != 6 {
		t.Errorf("expected 6-char suffix, got %s (len %d)", parts[1], len(parts[1]))
	}
}

func TestGenerateUsername_UsesFirstName(t *testing.T) {
	result := GenerateUsername("Alice Smith")
	prefix := strings.SplitN(result, "_", 2)[0]
	if prefix != "alice" {
		t.Errorf("expected 'alice', got '%s'", prefix)
	}
}

func TestGenerateUsername_TruncatesLongFirstName(t *testing.T) {
	result := GenerateUsername("Alexander Superlong")
	prefix := strings.SplitN(result, "_", 2)[0]
	if len(prefix) != 8 {
		t.Errorf("expected 8-char prefix, got '%s' (len %d)", prefix, len(prefix))
	}
	if prefix != "alexande" {
		t.Errorf("expected 'alexande', got '%s'", prefix)
	}
}

func TestGenerateUsername_SingleWord(t *testing.T) {
	result := GenerateUsername("Madonna")
	prefix := strings.SplitN(result, "_", 2)[0]
	if prefix != "madonna" {
		t.Errorf("expected 'madonna', got '%s'", prefix)
	}
}
