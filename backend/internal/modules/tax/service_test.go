package tax

import (
	"basekarya-backend/pkg/constants"
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newMockRepo() *mockRepo {
	return &mockRepo{}
}

func stubCategoryABrackets() []TERBracket {
	return []TERBracket{
		{Category: "A", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Category: "A", BracketNumber: 2, MinMonthlySalary: 5400000, Rate: 0.0025, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Category: "A", BracketNumber: 3, MinMonthlySalary: 10000000, Rate: 0.01, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Category: "A", BracketNumber: 4, MinMonthlySalary: 15000000, Rate: 0.02, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Category: "A", BracketNumber: 5, MinMonthlySalary: 100000000, Rate: 0.15, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
}

func stubCategoryBBrackets() []TERBracket {
	return []TERBracket{
		{Category: "B", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Category: "B", BracketNumber: 2, MinMonthlySalary: 6200000, Rate: 0.0075, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
}

func stubCategoryCBrackets() []TERBracket {
	return []TERBracket{
		{Category: "C", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Category: "C", BracketNumber: 2, MinMonthlySalary: 5000000, Rate: 0.02, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
}

func TestCalculateTER_CategoryA_Single(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	brackets := stubCategoryABrackets()
	m.On("FindTERBrackets", ctx, "A", mock.AnythingOfType("time.Time")).Return(brackets, nil)

	result, err := svc.CalculateTER(ctx, 10000000, constants.MaritalStatusSingle, 0)
	assert.NoError(t, err)
	assert.Equal(t, "A", result.TERCategory)
	assert.Equal(t, "TK/0", result.PTKPCode)
	assert.Equal(t, float64(10000000), result.GrossMonthly)
	assert.Equal(t, float64(0.01), result.TERRate)
	assert.Equal(t, float64(100000), result.MonthlyPPh21)
	m.AssertExpectations(t)
}

func TestCalculateTER_HighestBracket(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	brackets := stubCategoryABrackets()
	m.On("FindTERBrackets", ctx, "A", mock.AnythingOfType("time.Time")).Return(brackets, nil)

	result, err := svc.CalculateTER(ctx, 200000000, constants.MaritalStatusSingle, 0)
	assert.NoError(t, err)
	assert.Equal(t, float64(0.15), result.TERRate)
	m.AssertExpectations(t)
}

func TestCalculateTER_ExactBoundary(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	brackets := stubCategoryABrackets()
	m.On("FindTERBrackets", ctx, "A", mock.AnythingOfType("time.Time")).Return(brackets, nil)

	result, err := svc.CalculateTER(ctx, 5400000, constants.MaritalStatusSingle, 0)
	assert.NoError(t, err)
	assert.Equal(t, float64(0.0025), result.TERRate)
	m.AssertExpectations(t)
}

func TestCalculateTER_CategoryB(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	brackets := stubCategoryBBrackets()
	m.On("FindTERBrackets", ctx, "B", mock.AnythingOfType("time.Time")).Return(brackets, nil)

	result, err := svc.CalculateTER(ctx, 6200000, constants.MaritalStatusMarried, 1)
	assert.NoError(t, err)
	assert.Equal(t, "B", result.TERCategory)
	assert.Equal(t, "K/1", result.PTKPCode)
	assert.Equal(t, float64(0.0075), result.TERRate)
	m.AssertExpectations(t)
}

func TestCalculateTER_CategoryC(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	brackets := stubCategoryCBrackets()
	m.On("FindTERBrackets", ctx, "C", mock.AnythingOfType("time.Time")).Return(brackets, nil)

	result, err := svc.CalculateTER(ctx, 5000000, constants.MaritalStatusMarried, 3)
	assert.NoError(t, err)
	assert.Equal(t, "C", result.TERCategory)
	assert.Equal(t, "K/3", result.PTKPCode)
	m.AssertExpectations(t)
}

func TestCalculateTER_NoBrackets(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	m.On("FindTERBrackets", ctx, "A", mock.AnythingOfType("time.Time")).Return([]TERBracket{}, nil)

	_, err := svc.CalculateTER(ctx, 10000000, constants.MaritalStatusSingle, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no TER brackets found")
	m.AssertExpectations(t)
}

func TestCalculateTER_RepoError(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	m.On("FindTERBrackets", ctx, "A", mock.AnythingOfType("time.Time")).Return(nil, errors.New("db error"))

	_, err := svc.CalculateTER(ctx, 10000000, constants.MaritalStatusSingle, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find TER brackets")
	m.AssertExpectations(t)
}

func TestDerivePTKPCode(t *testing.T) {
	tests := []struct {
		status     constants.MaritalStatus
		dependents int
		expected   string
	}{
		{constants.MaritalStatusSingle, 0, "TK/0"},
		{constants.MaritalStatusSingle, 1, "TK/1"},
		{constants.MaritalStatusSingle, 3, "TK/3"},
		{constants.MaritalStatusSingle, 5, "TK/3"},
		{constants.MaritalStatusMarried, 0, "K/0"},
		{constants.MaritalStatusMarried, 1, "K/1"},
		{constants.MaritalStatusMarried, 2, "K/2"},
		{constants.MaritalStatusMarried, 3, "K/3"},
		{constants.MaritalStatusMarried, 5, "K/3"},
	}
	for _, tt := range tests {
		result := derivePTKPCode(tt.status, tt.dependents)
		assert.Equal(t, tt.expected, result, "status=%s deps=%d", tt.status, tt.dependents)
	}
}

func TestDeriveTERCategory(t *testing.T) {
	tests := []struct {
		ptkpCode string
		expected string
	}{
		{"TK/0", "A"},
		{"TK/1", "A"},
		{"K/0", "A"},
		{"TK/2", "B"},
		{"TK/3", "B"},
		{"K/1", "B"},
		{"K/2", "B"},
		{"K/3", "C"},
		{"K/4", "A"},
		{"UNKNOWN", "A"},
	}
	for _, tt := range tests {
		result := deriveTERCategory(tt.ptkpCode)
		assert.Equal(t, tt.expected, result, "ptkpCode=%s", tt.ptkpCode)
	}
}

func TestFindTERRate(t *testing.T) {
	brackets := stubCategoryABrackets()

	tests := []struct {
		gross    float64
		expected float64
	}{
		{0, 0.0},
		{3000000, 0.0},
		{5400000, 0.0025},
		{7000000, 0.0025},
		{10000000, 0.01},
		{12000000, 0.01},
		{15000000, 0.02},
		{50000000, 0.02},
		{100000000, 0.15},
		{200000000, 0.15},
	}
	for _, tt := range tests {
		rate := findTERRate(brackets, tt.gross)
		assert.Equal(t, tt.expected, rate, "gross=%.0f", tt.gross)
	}
}

func TestCalculateBiayaJabatanMonthly(t *testing.T) {
	tests := []struct {
		gross    float64
		expected float64
	}{
		{0, 0},
		{5000000, 250000},
		{10000000, 500000},
		{20000000, 500000},
		{50000000, 500000},
	}
	for _, tt := range tests {
		result := calculateBiayaJabatanMonthly(tt.gross)
		assert.Equal(t, tt.expected, result, "gross=%.0f", tt.gross)
	}
}

func TestCalculateProgressiveTax(t *testing.T) {
	tests := []struct {
		pkp      float64
		expected float64
	}{
		{0, 0},
		{50000000, 2500000},
		{60000000, 3000000},
		{100000000, 9000000},
		{250000000, 31500000},
		{500000000, 94000000},
		{1000000000, 244000000},
	}
	for _, tt := range tests {
		result := calculateProgressiveTax(tt.pkp)
		assert.Equal(t, tt.expected, result, "pkp=%.0f", tt.pkp)
	}
}

func TestReconcileAnnual(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	ptkps := []PTKPConfig{
		{Code: "TK/0", AnnualAmount: 54000000, EffectiveYear: 2026},
		{Code: "K/0", AnnualAmount: 58500000, EffectiveYear: 2026},
		{Code: "K/1", AnnualAmount: 63000000, EffectiveYear: 2026},
	}
	m.On("FindPTKPByYear", ctx, 2026).Return(ptkps, nil)

	grossAnnual := 200000000.0
	monthlyDetails := []MonthlyTaxEntry{
		{Month: 1, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 2, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 3, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 4, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 5, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 6, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 7, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 8, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 9, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 10, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 11, GrossIncome: 16000000, PPh21Paid: 80000},
		{Month: 12, GrossIncome: 16000000, PPh21Paid: 80000},
	}

	result, err := svc.ReconcileAnnual(ctx, 2026, "K/1", grossAnnual, monthlyDetails)
	assert.NoError(t, err)

	assert.Equal(t, grossAnnual, result.GrossAnnual)

	expectedBiayaJabatan := math.Min(grossAnnual*0.05, 6000000)
	assert.Equal(t, expectedBiayaJabatan, result.BiayaJabatan)

	assert.Equal(t, float64(63000000), result.PTKP)

	expectedPKP := math.Floor((grossAnnual - expectedBiayaJabatan - 63000000) / 1000) * 1000
	assert.Equal(t, expectedPKP, result.PKP)

	expectedTaxPayable := calculateProgressiveTax(expectedPKP)
	assert.Equal(t, expectedTaxPayable, result.TaxPayable)

	expectedTERPaidYTD := 12.0 * 80000.0
	assert.Equal(t, expectedTERPaidYTD, result.TERPaidYTD)

	assert.Equal(t, expectedTaxPayable-expectedTERPaidYTD, result.Delta)
	m.AssertExpectations(t)
}

func TestReconcileAnnual_PTKPNotFound(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	m.On("FindPTKPByYear", ctx, 2026).Return([]PTKPConfig{}, nil)

	grossAnnual := 200000000.0
	monthlyDetails := []MonthlyTaxEntry{}
	result, err := svc.ReconcileAnnual(ctx, 2026, "K/1", grossAnnual, monthlyDetails)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), result.PTKP)
	m.AssertExpectations(t)
}

func TestReconcileAnnual_RepoError(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	m.On("FindPTKPByYear", ctx, 2026).Return(nil, errors.New("db error"))

	monthlyDetails := []MonthlyTaxEntry{}
	_, err := svc.ReconcileAnnual(ctx, 2026, "K/1", 200000000, monthlyDetails)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find PTKP configs")
	m.AssertExpectations(t)
}

func TestService_CRUD_TERBracket(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	req := &TERBracketRequest{
		Category: "A", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0,
		EffectiveFrom: "2026-01-01",
	}
	m.On("CreateTERBracket", ctx, mock.AnythingOfType("*tax.TERBracket")).Return(nil)
	err := svc.CreateTERBracket(ctx, req)
	assert.NoError(t, err)

	m.On("FindTERBracketByID", ctx, uint(1)).Return(&TERBracket{ID: 1, Category: "A"}, nil)
	found, err := svc.GetTERBracketByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "A", found.Category)

	m.On("FindTERBracketByID", ctx, uint(1)).Return(&TERBracket{ID: 1, Category: "A"}, nil)
	m.On("UpdateTERBracket", ctx, mock.AnythingOfType("*tax.TERBracket")).Return(nil)
	err = svc.UpdateTERBracket(ctx, 1, req)
	assert.NoError(t, err)

	m.On("DeleteTERBracket", ctx, uint(1)).Return(nil)
	err = svc.DeleteTERBracket(ctx, 1)
	assert.NoError(t, err)

	m.AssertExpectations(t)
}

func TestService_CRUD_TERBracket_InvalidDate(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	req := &TERBracketRequest{
		Category: "A", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0,
		EffectiveFrom: "not-a-date",
	}
	err := svc.CreateTERBracket(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid effective_from")
}

func TestService_CRUD_PTKPConfig(t *testing.T) {
	m := newMockRepo()
	svc := NewService(m)
	ctx := context.Background()

	req := &PTKPConfigRequest{Code: "TK/0", AnnualAmount: 54000000, EffectiveYear: 2026}
	m.On("CreatePTKPConfig", ctx, mock.AnythingOfType("*tax.PTKPConfig")).Return(nil)
	err := svc.CreatePTKPConfig(ctx, req)
	assert.NoError(t, err)

	m.On("FindPTKPConfigByID", ctx, uint(1)).Return(&PTKPConfig{ID: 1, Code: "TK/0"}, nil)
	found, err := svc.GetPTKPConfigByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "TK/0", found.Code)

	m.On("FindPTKPConfigByID", ctx, uint(1)).Return(&PTKPConfig{ID: 1, Code: "TK/0"}, nil)
	m.On("UpdatePTKPConfig", ctx, mock.AnythingOfType("*tax.PTKPConfig")).Return(nil)
	err = svc.UpdatePTKPConfig(ctx, 1, req)
	assert.NoError(t, err)

	m.On("DeletePTKPConfig", ctx, uint(1)).Return(nil)
	err = svc.DeletePTKPConfig(ctx, 1)
	assert.NoError(t, err)

	m.AssertExpectations(t)
}
