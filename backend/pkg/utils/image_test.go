package utils

import (
	"encoding/base64"
	"testing"
)

func TestDecodeBase64Image_DataURIPrefix(t *testing.T) {
	original := []byte("hello image")
	encoded := base64.StdEncoding.EncodeToString(original)
	dataURI := "data:image/png;base64," + encoded

	decoded, err := DecodeBase64Image(dataURI)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(decoded) != string(original) {
		t.Errorf("expected %s, got %s", original, decoded)
	}
}

func TestDecodeBase64Image_PlainBase64(t *testing.T) {
	original := []byte("raw bytes")
	encoded := base64.StdEncoding.EncodeToString(original)

	decoded, err := DecodeBase64Image(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(decoded) != string(original) {
		t.Errorf("expected %s, got %s", original, decoded)
	}
}

func TestDecodeBase64Image_Invalid(t *testing.T) {
	_, err := DecodeBase64Image("!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
}
