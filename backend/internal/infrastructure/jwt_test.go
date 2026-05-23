package infrastructure

import (
	"basekarya-backend/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJwtProvider_GenerateAndValidateToken(t *testing.T) {
	p := NewJWTProvider(&config.JWTConfig{Secret: "test-secret", ExpiresIn: 1})

	employeeID := uint(10)
	token, err := p.GenerateToken(1, 2, true, "ADMIN", &employeeID, []string{"READ", "WRITE"})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := p.ValidateToken(token)
	require.NoError(t, err)

	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, uint(2), claims.CompanyID)
	assert.True(t, claims.IsPlatformAdmin)
	assert.Equal(t, "ADMIN", claims.Role)
	assert.Equal(t, &employeeID, claims.EmployeeID)
	assert.Equal(t, []string{"READ", "WRITE"}, claims.Permissions)
}

func TestJwtProvider_ValidateToken_Expired(t *testing.T) {
	p := NewJWTProvider(&config.JWTConfig{Secret: "test-secret", ExpiresIn: 0})

	token, err := p.GenerateToken(1, 1, false, "USER", nil, nil)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	_, err = p.ValidateToken(token)
	assert.Error(t, err)
}

func TestJwtProvider_ValidateToken_InvalidString(t *testing.T) {
	p := NewJWTProvider(&config.JWTConfig{Secret: "test-secret", ExpiresIn: 1})

	_, err := p.ValidateToken("not.a.valid.token")
	assert.Error(t, err)
}

func TestJwtProvider_ValidateToken_WrongSecret(t *testing.T) {
	p1 := NewJWTProvider(&config.JWTConfig{Secret: "secret-one", ExpiresIn: 1})
	p2 := NewJWTProvider(&config.JWTConfig{Secret: "secret-two", ExpiresIn: 1})

	token, err := p1.GenerateToken(1, 1, false, "USER", nil, nil)
	require.NoError(t, err)

	_, err = p2.ValidateToken(token)
	assert.Error(t, err)
}

func TestJwtProvider_TokenContainsAllClaims(t *testing.T) {
	p := NewJWTProvider(&config.JWTConfig{Secret: "my-secret", ExpiresIn: 24})

	employeeID := uint(99)
	perms := []string{"CREATE", "READ", "UPDATE", "DELETE"}
	token, err := p.GenerateToken(42, 7, false, "MANAGER", &employeeID, perms)
	require.NoError(t, err)

	claims, err := p.ValidateToken(token)
	require.NoError(t, err)

	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, uint(7), claims.CompanyID)
	assert.False(t, claims.IsPlatformAdmin)
	assert.Equal(t, "MANAGER", claims.Role)
	assert.NotNil(t, claims.EmployeeID)
	assert.Equal(t, uint(99), *claims.EmployeeID)
	assert.Equal(t, perms, claims.Permissions)
	assert.Equal(t, "hris-app", claims.Issuer)
}
