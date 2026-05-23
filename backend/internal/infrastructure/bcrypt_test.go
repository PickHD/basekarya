package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBcryptHasher_HashPassword_ReturnsNonEmpty(t *testing.T) {
	h := NewBcryptHasher(4)

	hash, err := h.HashPassword("password123")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestBcryptHasher_CheckPasswordHash_Correct(t *testing.T) {
	h := NewBcryptHasher(4)

	hash, err := h.HashPassword("mypassword")
	require.NoError(t, err)

	assert.True(t, h.CheckPasswordHash("mypassword", hash))
}

func TestBcryptHasher_CheckPasswordHash_Wrong(t *testing.T) {
	h := NewBcryptHasher(4)

	hash, err := h.HashPassword("mypassword")
	require.NoError(t, err)

	assert.False(t, h.CheckPasswordHash("wrongpassword", hash))
}

func TestBcryptHasher_DifferentPasswordsDifferentHashes(t *testing.T) {
	h := NewBcryptHasher(4)

	hash1, err := h.HashPassword("password1")
	require.NoError(t, err)

	hash2, err := h.HashPassword("password2")
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2)
}
