package tax

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockServiceHandler struct {
	mock.Mock
}

func (m *mockServiceHandler) CalculateTER(ctx context.Context, gross float64, ms constants.MaritalStatus, dc int) (*PPh21Result, error) {
	args := m.Called(ctx, gross, ms, dc)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PPh21Result), args.Error(1)
}

func (m *mockServiceHandler) ReconcileAnnual(ctx context.Context, year int, code string, gross float64, md []MonthlyTaxEntry) (*AnnualSettlement, error) {
	args := m.Called(ctx, year, code, gross, md)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AnnualSettlement), args.Error(1)
}

func (m *mockServiceHandler) CreateTERBracket(ctx context.Context, req *TERBracketRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockServiceHandler) GetTERBracketByID(ctx context.Context, id uint) (*TERBracket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TERBracket), args.Error(1)
}

func (m *mockServiceHandler) UpdateTERBracket(ctx context.Context, id uint, req *TERBracketRequest) error {
	return m.Called(ctx, id, req).Error(0)
}

func (m *mockServiceHandler) DeleteTERBracket(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockServiceHandler) ListTERBrackets(ctx context.Context, filter TERBracketFilter) ([]TERBracket, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]TERBracket), args.Get(1).(int64), args.Error(2)
}

func (m *mockServiceHandler) CreatePTKPConfig(ctx context.Context, req *PTKPConfigRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockServiceHandler) GetPTKPConfigByID(ctx context.Context, id uint) (*PTKPConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PTKPConfig), args.Error(1)
}

func (m *mockServiceHandler) UpdatePTKPConfig(ctx context.Context, id uint, req *PTKPConfigRequest) error {
	return m.Called(ctx, id, req).Error(0)
}

func (m *mockServiceHandler) DeletePTKPConfig(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockServiceHandler) ListPTKPConfigs(ctx context.Context, year int) ([]PTKPConfig, int64, error) {
	args := m.Called(ctx, year)
	return args.Get(0).([]PTKPConfig), args.Get(1).(int64), args.Error(2)
}

func TestHandler_ListTERBrackets(t *testing.T) {
	svc := new(mockServiceHandler)
	handler := NewHandler(svc)

	brackets := []TERBracket{
		{ID: 1, Category: "A", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0},
		{ID: 2, Category: "A", BracketNumber: 2, MinMonthlySalary: 5400000, Rate: 0.0025},
	}
	svc.On("ListTERBrackets", mock.Anything, mock.Anything).Return(brackets, int64(2), nil)

	at := testutil.NewAPITest(t, http.MethodGet, "/api/v1/admin/tax/ter-brackets?category=A&page=1&limit=10", nil)
	at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

	rec, err := at.Execute(handler.ListTERBrackets)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Message string `json:"message"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp.Message)
}

func TestHandler_CreateTERBracket(t *testing.T) {
	svc := new(mockServiceHandler)
	handler := NewHandler(svc)

	svc.On("CreateTERBracket", mock.Anything, mock.AnythingOfType("*tax.TERBracketRequest")).Return(nil)

	at := testutil.NewAPITest(t, http.MethodPost, "/api/v1/admin/tax/ter-brackets", map[string]interface{}{
		"category":           "A",
		"bracket_number":     1,
		"min_monthly_salary": 5400000,
		"rate":               0.0025,
		"effective_from":     "2026-01-01",
	})
	at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

	rec, err := at.Execute(handler.CreateTERBracket)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestHandler_ListPTKPConfigs(t *testing.T) {
	svc := new(mockServiceHandler)
	handler := NewHandler(svc)

	ptkpConfigs := []PTKPConfig{
		{ID: 1, Code: "TK/0", AnnualAmount: 54000000, EffectiveYear: 2026},
	}
	svc.On("ListPTKPConfigs", mock.Anything, 2026).Return(ptkpConfigs, int64(1), nil)

	at := testutil.NewAPITest(t, http.MethodGet, "/api/v1/admin/tax/ptkp-configs?year=2026", nil)
	at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

	rec, err := at.Execute(handler.ListPTKPConfigs)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
