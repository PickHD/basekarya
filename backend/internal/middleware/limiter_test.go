package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiterMiddleware_Init_ReturnsNonNil(t *testing.T) {
	m := NewRateLimiterMiddleware()
	mw := m.Init()
	assert.NotNil(t, mw)
}
