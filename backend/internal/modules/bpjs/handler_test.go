package bpjs

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) CalculateAll(ctx context.Context, grossMonthlyIncome float64) ([]BPJSComponent, error) {
	args := m.Called(ctx, grossMonthlyIncome)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]BPJSComponent), args.Error(1)
}

func (m *mockService) Create(ctx context.Context, req *BPJSRateConfigRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetByID(ctx context.Context, id uint) (*BPJSRateConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BPJSRateConfig), args.Error(1)
}

func (m *mockService) Update(ctx context.Context, id uint, req *BPJSRateConfigRequest) error {
	return m.Called(ctx, id, req).Error(0)
}

func (m *mockService) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) List(ctx context.Context, filter BPJSRateConfigFilter) ([]BPJSRateConfig, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]BPJSRateConfig), args.Get(1).(int64), args.Error(2)
}

func TestHandler_List(t *testing.T) {
	svc := new(mockService)
	handler := NewHandler(svc)

	configs := []BPJSRateConfig{
		{ID: 1, Type: "JKK", EmployeeRate: 0.0, EmployerRate: 0.0024},
		{ID: 2, Type: "JHT", EmployeeRate: 0.02, EmployerRate: 0.037},
	}
	svc.On("List", mock.Anything, mock.Anything).Return(configs, int64(2), nil)

	at := testutil.NewAPITest(t, http.MethodGet, "/api/v1/admin/bpjs/configs?page=1&limit=10", nil)
	at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

	rec, err := at.Execute(handler.List)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Message string `json:"message"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp.Message)
}

func TestHandler_Create(t *testing.T) {
	svc := new(mockService)
	handler := NewHandler(svc)

	svc.On("Create", mock.Anything, mock.AnythingOfType("*bpjs.BPJSRateConfigRequest")).Return(nil)

	at := testutil.NewAPITest(t, http.MethodPost, "/api/v1/admin/bpjs/configs", map[string]interface{}{
		"type":           "JHT",
		"employee_rate":  0.02,
		"employer_rate":  0.037,
		"is_active":      true,
		"effective_from": "2026-01-01",
	})
	at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

	rec, err := at.Execute(handler.Create)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}
