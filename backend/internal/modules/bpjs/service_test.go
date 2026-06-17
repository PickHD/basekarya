package bpjs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateAll(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)
	ctx := context.Background()

	cap12M := float64(12000000)
	configs := []BPJSRateConfig{
		{Type: "KESEHATAN", EmployeeRate: 0.01, EmployerRate: 0.04, MaxSalaryCap: &cap12M, IsActive: true},
		{Type: "JHT", EmployeeRate: 0.02, EmployerRate: 0.037, IsActive: true},
		{Type: "JKK", EmployeeRate: 0, EmployerRate: 0.0054, IsActive: true},
		{Type: "JKM", EmployeeRate: 0, EmployerRate: 0.003, IsActive: true},
		{Type: "JP", EmployeeRate: 0.01, EmployerRate: 0.02, IsActive: true},
	}

	repo.On("FindAllActive", ctx, mock.Anything).Return(configs, nil)

	components, err := svc.CalculateAll(ctx, 10000000)
	assert.NoError(t, err)
	assert.Len(t, components, 8)

	assert.Equal(t, "KESEHATAN", components[0].Type)
	assert.Equal(t, "BPJS_KESEHATAN_E", components[0].Code)
	assert.Equal(t, float64(100000), components[0].EmployeeAmount)
	assert.Equal(t, float64(0), components[0].EmployerAmount)
	assert.False(t, components[0].IsEmployerBorne)

	assert.Equal(t, "KESEHATAN", components[1].Type)
	assert.Equal(t, "BPJS_KESEHATAN_R", components[1].Code)
	assert.Equal(t, float64(0), components[1].EmployeeAmount)
	assert.Equal(t, float64(400000), components[1].EmployerAmount)
	assert.True(t, components[1].IsEmployerBorne)

	repo.AssertExpectations(t)
}

func TestCalculateAll_WithCap(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)
	ctx := context.Background()

	cap12M := float64(12000000)
	configs := []BPJSRateConfig{
		{Type: "KESEHATAN", EmployeeRate: 0.01, EmployerRate: 0.04, MaxSalaryCap: &cap12M, IsActive: true},
	}

	t.Run("salary above cap", func(t *testing.T) {
		repo.On("FindAllActive", ctx, mock.Anything).Return(configs, nil).Once()

		components, err := svc.CalculateAll(ctx, 20000000)
		assert.NoError(t, err)
		assert.Len(t, components, 2)
		assert.Equal(t, float64(120000), components[0].EmployeeAmount)
	})

	t.Run("salary below cap", func(t *testing.T) {
		repo.On("FindAllActive", ctx, mock.Anything).Return(configs, nil).Once()

		components, err := svc.CalculateAll(ctx, 5000000)
		assert.NoError(t, err)
		assert.Len(t, components, 2)
		assert.Equal(t, float64(50000), components[0].EmployeeAmount)
	})

	repo.AssertExpectations(t)
}

func TestCalculateAll_NoConfig(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)
	ctx := context.Background()

	repo.On("FindAllActive", ctx, mock.Anything).Return([]BPJSRateConfig{}, nil)

	components, err := svc.CalculateAll(ctx, 10000000)
	assert.NoError(t, err)
	assert.Len(t, components, 0)

	repo.AssertExpectations(t)
}

func TestCalculateAll_InactiveConfig(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)
	ctx := context.Background()

	configs := []BPJSRateConfig{
		{Type: "KESEHATAN", EmployeeRate: 0.01, EmployerRate: 0.04, IsActive: false},
	}

	repo.On("FindAllActive", ctx, mock.Anything).Return(configs, nil)

	components, err := svc.CalculateAll(ctx, 10000000)
	assert.NoError(t, err)
	assert.Len(t, components, 0)

	repo.AssertExpectations(t)
}

func TestService_CRUD(t *testing.T) {
	repo := new(mockRepo)
	svc := NewService(repo)
	ctx := context.Background()

	req := &BPJSRateConfigRequest{
		Type:          "JHT",
		EmployeeRate:  0.02,
		EmployerRate:  0.037,
		IsActive:      true,
		EffectiveFrom: "2026-01-01",
	}

	repo.On("Create", ctx, mock.AnythingOfType("*bpjs.BPJSRateConfig")).Return(nil)

	err := svc.Create(ctx, req)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}
