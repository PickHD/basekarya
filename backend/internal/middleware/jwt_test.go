package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestAuthMiddleware(t *testing.T) (*AuthMiddleware, *infrastructure.JwtProvider) {
	t.Helper()
	jwtProvider := testutil.NewTestJWT()
	authMW := NewAuthMiddleware(jwtProvider)
	return authMW, jwtProvider
}

func okHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "ok")
}

func TestAuthMiddleware_VerifyToken_ValidBearerToken(t *testing.T) {
	authMW, jwtProvider := newTestAuthMiddleware(t)
	token := testutil.GenerateTestToken(t, jwtProvider, 1, 1, false, "ADMIN", []string{"read"})

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithToken(token)

	handler := authMW.VerifyToken(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_VerifyToken_MissingAuthorizationHeader(t *testing.T) {
	authMW, _ := newTestAuthMiddleware(t)

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)

	handler := authMW.VerifyToken(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_VerifyToken_InvalidToken(t *testing.T) {
	authMW, _ := newTestAuthMiddleware(t)

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithToken("invalid-token")

	handler := authMW.VerifyToken(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_VerifyToken_QueryParamToken(t *testing.T) {
	authMW, jwtProvider := newTestAuthMiddleware(t)
	token := testutil.GenerateTestToken(t, jwtProvider, 1, 1, false, "ADMIN", []string{"read"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ws?token="+token, nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/v1/ws")

	handler := authMW.VerifyToken(okHandler)
	err := handler(ctx)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_VerifyToken_SetsClaimsOnContext(t *testing.T) {
	authMW, jwtProvider := newTestAuthMiddleware(t)
	token := testutil.GenerateTestToken(t, jwtProvider, 1, 1, false, "ADMIN", []string{"read"})

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithToken(token)

	var gotClaims *infrastructure.MyClaims
	handler := authMW.VerifyToken(func(ctx echo.Context) error {
		val := ctx.Get("user")
		gotClaims = val.(*infrastructure.MyClaims)
		return ctx.String(http.StatusOK, "ok")
	})

	_, err := at.Execute(handler)
	require.NoError(t, err)

	require.NotNil(t, gotClaims)
	assert.Equal(t, uint(1), gotClaims.UserID)
	assert.Equal(t, uint(1), gotClaims.CompanyID)
	assert.Equal(t, "ADMIN", gotClaims.Role)
	assert.Equal(t, []string{"read"}, gotClaims.Permissions)
}

func TestAuthMiddleware_VerifyToken_SetsContextValues(t *testing.T) {
	authMW, jwtProvider := newTestAuthMiddleware(t)
	token := testutil.GenerateTestToken(t, jwtProvider, 5, 10, true, "ADMIN", nil)

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithToken(token)

	var reqCtx context.Context
	handler := authMW.VerifyToken(func(ctx echo.Context) error {
		reqCtx = ctx.Request().Context()
		return ctx.String(http.StatusOK, "ok")
	})

	_, err := at.Execute(handler)
	require.NoError(t, err)

	assert.Equal(t, uint(5), reqCtx.Value(constants.UserIDContextKey))
	assert.Equal(t, uint(10), reqCtx.Value(constants.CompanyIDContextKey))
	assert.Equal(t, true, reqCtx.Value(constants.IsPlatformAdminContextKey))
}

func TestAuthMiddleware_GrantPermission_UserHasPermission(t *testing.T) {
	authMW, _ := newTestAuthMiddleware(t)
	claims := &infrastructure.MyClaims{
		UserID:          1,
		CompanyID:       1,
		IsPlatformAdmin: false,
		Permissions:     []string{"users.read", "users.write"},
	}

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithAuthContext(claims)

	handler := authMW.GrantPermission("users.read")(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_GrantPermission_UserNoPermission(t *testing.T) {
	authMW, _ := newTestAuthMiddleware(t)
	claims := &infrastructure.MyClaims{
		UserID:          1,
		CompanyID:       1,
		IsPlatformAdmin: false,
		Permissions:     []string{"users.read"},
	}

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithAuthContext(claims)

	handler := authMW.GrantPermission("users.delete")(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthMiddleware_GrantPermission_PlatformAdminBypass(t *testing.T) {
	authMW, _ := newTestAuthMiddleware(t)
	claims := &infrastructure.MyClaims{
		UserID:          1,
		CompanyID:       1,
		IsPlatformAdmin: true,
		Permissions:     []string{},
	}

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithAuthContext(claims)

	handler := authMW.GrantPermission("users.delete")(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_GrantAnyPermission_UserHasOne(t *testing.T) {
	authMW, _ := newTestAuthMiddleware(t)
	claims := &infrastructure.MyClaims{
		UserID:          1,
		CompanyID:       1,
		IsPlatformAdmin: false,
		Permissions:     []string{"users.read"},
	}

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithAuthContext(claims)

	handler := authMW.GrantAnyPermission("users.delete", "users.read")(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_GrantAnyPermission_UserHasNone(t *testing.T) {
	authMW, _ := newTestAuthMiddleware(t)
	claims := &infrastructure.MyClaims{
		UserID:          1,
		CompanyID:       1,
		IsPlatformAdmin: false,
		Permissions:     []string{"settings.read"},
	}

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithAuthContext(claims)

	handler := authMW.GrantAnyPermission("users.delete", "users.write")(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequirePlatformAdmin_IsAdmin(t *testing.T) {
	authMW, _ := newTestAuthMiddleware(t)
	claims := &infrastructure.MyClaims{
		UserID:          1,
		CompanyID:       1,
		IsPlatformAdmin: true,
	}

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithAuthContext(claims)

	handler := RequirePlatformAdmin(authMW)(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequirePlatformAdmin_NotAdmin(t *testing.T) {
	authMW, _ := newTestAuthMiddleware(t)
	claims := &infrastructure.MyClaims{
		UserID:          1,
		CompanyID:       1,
		IsPlatformAdmin: false,
	}

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithAuthContext(claims)

	handler := RequirePlatformAdmin(authMW)(okHandler)
	rec, err := at.Execute(handler)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
