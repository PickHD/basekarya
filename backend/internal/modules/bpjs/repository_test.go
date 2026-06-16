package bpjs

import (
	"basekarya-backend/internal/testutil"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupBPJSTestDB(t *testing.T) *testutil.TestDB {
	db := testutil.NewTestDB(&BPJSRateConfig{})
	return db
}

func TestRepository_FindActiveByType(t *testing.T) {
	db := setupBPJSTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	cfg := &BPJSRateConfig{
		Type:          "KESEHATAN",
		EmployeeRate:  0.01,
		EmployerRate:  0.04,
		MaxSalaryCap:  float64Ptr(12000000),
		IsActive:      true,
		EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	db.DB.Create(cfg)

	configs, err := repo.FindActiveByType(ctx, "KESEHATAN", time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	assert.Len(t, configs, 1)
	assert.Equal(t, "KESEHATAN", configs[0].Type)
	assert.Equal(t, 0.01, configs[0].EmployeeRate)
	assert.Equal(t, 0.04, configs[0].EmployerRate)
	assert.Equal(t, float64(12000000), *configs[0].MaxSalaryCap)
}

func TestRepository_FindActiveByType_Inactive(t *testing.T) {
	db := setupBPJSTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	cfg := &BPJSRateConfig{
		Type:          "KESEHATAN",
		EmployeeRate:  0.01,
		EmployerRate:  0.04,
		IsActive:      true,
		EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	db.DB.Create(cfg)
	db.DB.Model(cfg).Update("is_active", false)

	configs, err := repo.FindActiveByType(ctx, "KESEHATAN", time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	assert.Len(t, configs, 0)
}

func TestRepository_FindAllActive(t *testing.T) {
	db := setupBPJSTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	types := []string{"KESEHATAN", "JHT", "JKK"}
	for _, typ := range types {
		cfg := &BPJSRateConfig{
			Type:          typ,
			EmployeeRate:  0.01,
			EmployerRate:  0.04,
			IsActive:      true,
			EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		db.DB.Create(cfg)
	}

	configs, err := repo.FindAllActive(ctx, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	assert.Len(t, configs, 3)
}

func TestRepository_CRUD(t *testing.T) {
	db := setupBPJSTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	cfg := &BPJSRateConfig{
		Type:          "JHT",
		EmployeeRate:  0.02,
		EmployerRate:  0.037,
		IsActive:      true,
		EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	err := repo.Create(ctx, cfg)
	assert.NoError(t, err)
	assert.NotZero(t, cfg.ID)

	found, err := repo.FindByID(ctx, cfg.ID)
	assert.NoError(t, err)
	assert.Equal(t, "JHT", found.Type)
	assert.Equal(t, 0.02, found.EmployeeRate)
	assert.Equal(t, 0.037, found.EmployerRate)

	found.EmployerRate = 0.04
	err = repo.Update(ctx, found)
	assert.NoError(t, err)

	updated, err := repo.FindByID(ctx, cfg.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0.04, updated.EmployerRate)

	err = repo.Delete(ctx, cfg.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, cfg.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BPJS rate config not found")
}

func float64Ptr(v float64) *float64 {
	return &v
}
