package department

import (
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEmployee struct {
	ID           uint `gorm:"primaryKey"`
	DepartmentID *uint
	CompanyID    uint
}

func (testEmployee) TableName() string {
	return "employees"
}

func setupDeptTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(&Department{}, &testEmployee{})
	t.Cleanup(tdb.Close)
	return tdb
}

func seedDeptTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	companyID := uint(1)
	require.NoError(t, db.DB.Create(&Department{ID: 1, Name: "Engineering", CompanyID: companyID}).Error)
}

func TestRepo_FindAll(t *testing.T) {
	tdb := setupDeptTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedDeptTestData(t, tdb)

	deps, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.Len(t, deps, 1)
	assert.Equal(t, "Engineering", deps[0].Name)
}

func TestRepo_FindByID(t *testing.T) {
	tdb := setupDeptTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedDeptTestData(t, tdb)

	dept, err := repo.FindByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "Engineering", dept.Name)
}

func TestRepo_FindByID_NotFound(t *testing.T) {
	tdb := setupDeptTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	_, err := repo.FindByID(ctx, 99)
	require.Error(t, err)
}

func TestRepo_FindByName(t *testing.T) {
	tdb := setupDeptTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedDeptTestData(t, tdb)

	dept, err := repo.FindByName(ctx, "Engineering")
	require.NoError(t, err)
	assert.Equal(t, "Engineering", dept.Name)
}

func TestRepo_Create(t *testing.T) {
	tdb := setupDeptTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	err := repo.Create(ctx, &Department{Name: "HR", CompanyID: 1})
	require.NoError(t, err)

	deps, _ := repo.FindAll(ctx)
	assert.Len(t, deps, 1)
	assert.Equal(t, "HR", deps[0].Name)
}

func TestRepo_Update(t *testing.T) {
	tdb := setupDeptTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedDeptTestData(t, tdb)

	dept, _ := repo.FindByID(ctx, 1)
	dept.Name = "Updated"
	err := repo.Update(ctx, dept)
	require.NoError(t, err)

	updated, _ := repo.FindByID(ctx, 1)
	assert.Equal(t, "Updated", updated.Name)
}

func TestRepo_Delete(t *testing.T) {
	tdb := setupDeptTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedDeptTestData(t, tdb)

	err := repo.Delete(ctx, 1)
	require.NoError(t, err)

	_, err = repo.FindByID(ctx, 1)
	require.Error(t, err)
}

func TestRepo_ExistsByName(t *testing.T) {
	tdb := setupDeptTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedDeptTestData(t, tdb)

	exists, err := repo.ExistsByName(ctx, "Engineering", 0)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByName(ctx, "Unknown", 0)
	require.NoError(t, err)
	assert.False(t, exists)
}
