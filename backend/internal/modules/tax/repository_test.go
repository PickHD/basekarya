package tax

import (
	"basekarya-backend/internal/testutil"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTaxTestDB(t *testing.T) *testutil.TestDB {
	db := testutil.NewTestDB(&TERBracket{}, &PTKPConfig{})
	return db
}

func TestRepository_FindTERBrackets(t *testing.T) {
	db := setupTaxTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	b1 := &TERBracket{Category: "A", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	b2 := &TERBracket{Category: "A", BracketNumber: 2, MinMonthlySalary: 5400000, Rate: 0.0025, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	db.DB.Create(b1)
	db.DB.Create(b2)

	brackets, err := repo.FindTERBrackets(ctx, "A", time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	assert.Len(t, brackets, 2)
}

func TestRepository_FindTERBrackets_Expired(t *testing.T) {
	db := setupTaxTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	until := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	b1 := &TERBracket{Category: "A", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), EffectiveUntil: &until}
	db.DB.Create(b1)

	brackets, err := repo.FindTERBrackets(ctx, "A", time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	assert.Len(t, brackets, 0)
}

func TestRepository_FindPTKPByYear(t *testing.T) {
	db := setupTaxTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	ptkp := &PTKPConfig{Code: "TK/0", AnnualAmount: 54000000, EffectiveYear: 2026}
	db.DB.Create(ptkp)

	configs, err := repo.FindPTKPByYear(ctx, 2026)
	assert.NoError(t, err)
	assert.Len(t, configs, 1)
	assert.Equal(t, "TK/0", configs[0].Code)
}

func TestRepository_FindPTKPByYear_NoResults(t *testing.T) {
	db := setupTaxTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	configs, err := repo.FindPTKPByYear(ctx, 2025)
	assert.NoError(t, err)
	assert.Len(t, configs, 0)
}

func TestRepository_CRUD_TERBracket(t *testing.T) {
	db := setupTaxTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	b := &TERBracket{Category: "A", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	err := repo.CreateTERBracket(ctx, b)
	assert.NoError(t, err)
	assert.NotZero(t, b.ID)

	found, err := repo.FindTERBracketByID(ctx, b.ID)
	assert.NoError(t, err)
	assert.Equal(t, "A", found.Category)

	found.Rate = 0.0050
	err = repo.UpdateTERBracket(ctx, found)
	assert.NoError(t, err)

	updated, _ := repo.FindTERBracketByID(ctx, b.ID)
	assert.Equal(t, 0.0050, updated.Rate)

	err = repo.DeleteTERBracket(ctx, b.ID)
	assert.NoError(t, err)

	_, err = repo.FindTERBracketByID(ctx, b.ID)
	assert.Error(t, err)
}

func TestRepository_CRUD_PTKPConfig(t *testing.T) {
	db := setupTaxTestDB(t)
	defer db.Close()
	repo := NewRepository(db.DB)
	ctx := context.Background()

	ptkp := &PTKPConfig{Code: "TK/0", AnnualAmount: 54000000, EffectiveYear: 2026}
	err := repo.CreatePTKPConfig(ctx, ptkp)
	assert.NoError(t, err)
	assert.NotZero(t, ptkp.ID)

	found, err := repo.FindPTKPConfigByID(ctx, ptkp.ID)
	assert.NoError(t, err)
	assert.Equal(t, "TK/0", found.Code)

	found.AnnualAmount = 55000000
	err = repo.UpdatePTKPConfig(ctx, found)
	assert.NoError(t, err)

	updated, _ := repo.FindPTKPConfigByID(ctx, ptkp.ID)
	assert.Equal(t, float64(55000000), updated.AnnualAmount)

	err = repo.DeletePTKPConfig(ctx, ptkp.ID)
	assert.NoError(t, err)

	_, err = repo.FindPTKPConfigByID(ctx, ptkp.ID)
	assert.Error(t, err)
}
