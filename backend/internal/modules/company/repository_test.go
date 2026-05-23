package company

import (
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCompanyTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(&Company{})
	t.Cleanup(tdb.Close)
	return tdb
}

func seedCompanyTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	require.NoError(t, db.DB.Create(&Company{
		ID:       1,
		Name:     "Test Company",
		Email:    "test@company.com",
		PhoneNumber: "081234567890",
		Address:  "123 Test St",
		Website:  "https://test.com",
		TaxNumber: "NPWP123",
	}).Error)
}

func TestRepoCompany_FindByID(t *testing.T) {
	tdb := setupCompanyTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedCompanyTestData(t, tdb)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: 1, wantErr: false},
		{name: "not found", id: 99, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp, err := repo.FindByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "Test Company", comp.Name)
			}
		})
	}
}

func TestRepoCompany_CreateCompany(t *testing.T) {
	tdb := setupCompanyTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	comp := &Company{
		Name:  "New Company",
		Email: "new@company.com",
	}

	err := repo.CreateCompany(ctx, comp)
	require.NoError(t, err)
	assert.NotZero(t, comp.ID)
}

func TestRepoCompany_Update(t *testing.T) {
	tdb := setupCompanyTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedCompanyTestData(t, tdb)

	comp, _ := repo.FindByID(ctx, 1)
	comp.Name = "Updated Company"
	comp.Address = "456 Updated St"

	err := repo.Update(ctx, comp)
	require.NoError(t, err)

	updated, _ := repo.FindByID(ctx, 1)
	assert.Equal(t, "Updated Company", updated.Name)
	assert.Equal(t, "456 Updated St", updated.Address)
}
