package finance

import (
	"database/sql"
	"testing"
	"time"

	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupFinanceTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	tdb := &testutil.TestDB{DB: db}
	t.Cleanup(tdb.Close)

	err = db.AutoMigrate(
		&rbac.Role{},
		&master.Department{},
		&master.Shift{},
		&user.User{},
		&user.Employee{},
	)
	require.NoError(t, err)

	db.Exec(`CREATE TABLE IF NOT EXISTS finance_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(20) NOT NULL,
		description TEXT,
		company_id INTEGER NOT NULL
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS finance_transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		finance_category_id INTEGER NOT NULL,
		created_by INTEGER NOT NULL,
		company_id INTEGER NOT NULL,
		approved_by INTEGER,
		type VARCHAR(20) NOT NULL,
		amount REAL NOT NULL,
		description TEXT,
		transaction_date DATE NOT NULL,
		reference_number VARCHAR(100),
		status VARCHAR(20) DEFAULT 'PENDING',
		rejection_reason TEXT
	)`)

	return tdb
}

func seedFinanceTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	companyID := uint(1)

	role := &rbac.Role{ID: 1, Name: "Admin", CompanyID: companyID}
	require.NoError(t, db.DB.Create(role).Error)

	dept := &master.Department{ID: 1, Name: "Engineering", CompanyID: companyID}
	require.NoError(t, db.DB.Create(dept).Error)

	shift := &master.Shift{ID: 1, Name: "Day", StartTime: "09:00", EndTime: "17:00", CompanyID: companyID}
	require.NoError(t, db.DB.Create(shift).Error)

	usr := &user.User{ID: 1, Username: "john", RoleID: 1, CompanyID: companyID}
	require.NoError(t, db.DB.Create(usr).Error)

	emp := &user.Employee{
		ID: 1, UserID: 1, CompanyID: companyID, DepartmentID: 1, ShiftID: 1,
		NIK: "EMP001", FullName: "John Doe",
	}
	require.NoError(t, db.DB.Create(emp).Error)

	cat := &FinanceCategory{
		ID: 1, Name: "Salary", Type: constants.FinanceTypeIncome, CompanyID: companyID,
	}
	require.NoError(t, db.DB.Table("finance_categories").Create(cat).Error)
}

func TestRepo_CreateCategory(t *testing.T) {
	tdb := setupFinanceTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name    string
		cat     *FinanceCategory
		wantErr bool
	}{
		{
			name: "success",
			cat: &FinanceCategory{
				Name: "Bonus", Type: constants.FinanceTypeIncome, CompanyID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateCategory(ctx, tt.cat)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.cat.ID)
			}
		})
	}
}

func TestRepo_FindCategoryByID(t *testing.T) {
	tdb := setupFinanceTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedFinanceTestData(t, tdb)

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
			cat, err := repo.FindCategoryByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, cat.ID)
				assert.Equal(t, "Salary", cat.Name)
			}
		})
	}
}

func TestRepo_FindAllCategories(t *testing.T) {
	tdb := setupFinanceTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedFinanceTestData(t, tdb)

	require.NoError(t, repo.CreateCategory(ctx, &FinanceCategory{
		Name: "Rent", Type: constants.FinanceTypeExpense, CompanyID: 1,
	}))

	tests := []struct {
		name      string
		catType   string
		wantCount int
		wantErr   bool
	}{
		{name: "all categories", catType: "", wantCount: 2, wantErr: false},
		{name: "income only", catType: "INCOME", wantCount: 1, wantErr: false},
		{name: "expense only", catType: "EXPENSE", wantCount: 1, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cats, err := repo.FindAllCategories(ctx, tt.catType)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, cats, tt.wantCount)
			}
		})
	}
}

func TestRepo_UpdateCategory(t *testing.T) {
	tdb := setupFinanceTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedFinanceTestData(t, tdb)

	cat, err := repo.FindCategoryByID(ctx, 1)
	require.NoError(t, err)

	cat.Name = "Monthly Salary"
	cat.Description = sql.NullString{String: "Regular salary", Valid: true}

	err = repo.UpdateCategory(ctx, cat)
	require.NoError(t, err)

	updated, err := repo.FindCategoryByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "Monthly Salary", updated.Name)
}

func TestRepo_DeleteCategory(t *testing.T) {
	tdb := setupFinanceTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedFinanceTestData(t, tdb)

	err := repo.DeleteCategory(ctx, 1)
	require.NoError(t, err)

	_, err = repo.FindCategoryByID(ctx, 1)
	require.Error(t, err)
}

func TestRepo_CreateTransaction(t *testing.T) {
	tdb := setupFinanceTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedFinanceTestData(t, tdb)

	tests := []struct {
		name    string
		tx      *FinanceTransaction
		wantErr bool
	}{
		{
			name: "success",
			tx: &FinanceTransaction{
				CompanyID:          1,
				FinanceCategoryID: 1,
				CreatedBy:         1,
				Type:              constants.FinanceTypeIncome,
				Amount:            5000000,
				TransactionDate:   time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				Status:            constants.FinanceStatusPending,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateTransaction(ctx, tt.tx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.tx.ID)
			}
		})
	}
}

func TestRepo_FindTransactionByID(t *testing.T) {
	tdb := setupFinanceTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedFinanceTestData(t, tdb)

	tx := &FinanceTransaction{
		CompanyID: 1, FinanceCategoryID: 1, CreatedBy: 1,
		Type: constants.FinanceTypeIncome, Amount: 5000000,
		TransactionDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		Status: constants.FinanceStatusPending,
	}
	require.NoError(t, repo.CreateTransaction(ctx, tx))

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: tx.ID, wantErr: false},
		{name: "not found", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindTransactionByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, found.ID)
				assert.NotNil(t, found.Creator)
				assert.NotNil(t, found.FinanceCategory)
			}
		})
	}
}

func TestRepo_FindAllTransactions(t *testing.T) {
	tdb := setupFinanceTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedFinanceTestData(t, tdb)

	for i := 0; i < 3; i++ {
		require.NoError(t, repo.CreateTransaction(ctx, &FinanceTransaction{
			CompanyID: 1, FinanceCategoryID: 1, CreatedBy: 1,
			Type: constants.FinanceTypeIncome, Amount: 1000000,
			TransactionDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			Status: constants.FinanceStatusPending,
		}))
	}

	tests := []struct {
		name       string
		filter     TransactionFilter
		wantCount  int
		wantCursor bool
		wantErr    bool
	}{
		{
			name:       "all transactions",
			filter:     TransactionFilter{Limit: 10},
			wantCount:  3,
			wantCursor: false,
			wantErr:    false,
		},
		{
			name:       "filter by type",
			filter:     TransactionFilter{Type: "INCOME", Limit: 10},
			wantCount:  3,
			wantCursor: false,
			wantErr:    false,
		},
		{
			name:       "filter by status pending",
			filter:     TransactionFilter{Status: "PENDING", Limit: 10},
			wantCount:  3,
			wantCursor: false,
			wantErr:    false,
		},
		{
			name:       "paginated with next cursor",
			filter:     TransactionFilter{Limit: 2},
			wantCount:  2,
			wantCursor: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txs, nextCursor, err := repo.FindAllTransactions(ctx, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, txs, tt.wantCount)
				if tt.wantCursor {
					assert.NotNil(t, nextCursor)
				} else {
					assert.Nil(t, nextCursor)
				}
			}
		})
	}
}

func TestRepo_UpdateTransaction(t *testing.T) {
	tdb := setupFinanceTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedFinanceTestData(t, tdb)

	tx := &FinanceTransaction{
		CompanyID: 1, FinanceCategoryID: 1, CreatedBy: 1,
		Type: constants.FinanceTypeIncome, Amount: 5000000,
		TransactionDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		Status: constants.FinanceStatusPending,
	}
	require.NoError(t, repo.CreateTransaction(ctx, tx))

	tx.Status = constants.FinanceStatusApproved
	approvedBy := uint(2)
	tx.ApprovedBy = &approvedBy

	err := repo.UpdateTransaction(ctx, tx)
	require.NoError(t, err)

	found, err := repo.FindTransactionByID(ctx, tx.ID)
	require.NoError(t, err)
	assert.Equal(t, constants.FinanceStatusApproved, found.Status)
}
