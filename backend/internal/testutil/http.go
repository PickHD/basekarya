package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

// APITest wraps echo context and response recorder for handler tests.
type APITest struct {
	Echo    *echo.Echo
	Req     *http.Request
	Rec     *httptest.ResponseRecorder
	Context echo.Context
}

// NewAPITest creates a new API test helper.
// method: HTTP method, target: URL path, body: optional JSON body (can be nil).
func NewAPITest(t *testing.T, method, target string, body interface{}) *APITest {
	t.Helper()

	var bodyReader *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(data)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, target, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()

	e := echo.New()
	e.Validator = utils.NewValidator()

	return &APITest{
		Echo:    e,
		Req:     req,
		Rec:     rec,
		Context: e.NewContext(req, rec),
	}
}

// WithToken sets an Authorization Bearer token on the request.
func (at *APITest) WithToken(token string) *APITest {
	at.Req.Header.Set("Authorization", "Bearer "+token)
	return at
}

// WithAuthContext injects auth claims into the echo context (bypasses JWT middleware).
func (at *APITest) WithAuthContext(claims *infrastructure.MyClaims) *APITest {
	at.Context.Set("user", claims)

	ctx := at.Context.Request().Context()
	ctx = context.WithValue(ctx, constants.CompanyIDContextKey, claims.CompanyID)
	ctx = context.WithValue(ctx, constants.IsPlatformAdminContextKey, claims.IsPlatformAdmin)
	ctx = context.WithValue(ctx, constants.UserIDContextKey, claims.UserID)
	at.Context.SetRequest(at.Context.Request().WithContext(ctx))
	return at
}

// WithPathParams sets path parameters on the echo context.
func (at *APITest) WithPathParams(params map[string]string) *APITest {
	names := make([]string, 0, len(params))
	values := make([]string, 0, len(params))
	for k, v := range params {
		names = append(names, k)
		values = append(values, v)
	}
	at.Context.SetParamNames(names...)
	at.Context.SetParamValues(values...)
	return at
}

// Execute runs the handler and returns the response.
func (at *APITest) Execute(handler echo.HandlerFunc) (*httptest.ResponseRecorder, error) {
	err := handler(at.Context)
	return at.Rec, err
}

// DecodeResponse decodes the response body into the target.
func DecodeResponse(t *testing.T, rec *httptest.ResponseRecorder, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), target); err != nil {
		t.Fatalf("failed to decode response: %v\nbody: %s", err, rec.Body.String())
	}
}

// GenerateTestToken generates a valid JWT token for the given user details.
func GenerateTestToken(t *testing.T, jwt *infrastructure.JwtProvider, userID, companyID uint, isPlatformAdmin bool, role string, permissions []string) string {
	t.Helper()
	employeeID := uint(1)
	token, err := jwt.GenerateToken(userID, companyID, isPlatformAdmin, role, &employeeID, permissions)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	return token
}
