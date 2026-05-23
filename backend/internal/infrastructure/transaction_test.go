package infrastructure

import (
	"context"
	"testing"
	"time"

	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type testRole struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	CompanyID uint      `gorm:"not null"`
	CreatedAt time.Time
}

func (testRole) TableName() string {
	return "test_roles"
}

func setupTransactionTestDB(t *testing.T) (*gorm.DB, TransactionManager) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&testRole{}))
	tm := NewGormTransactionManager(db)
	return db, tm
}

func TestRunInTransaction_CommitsOnSuccess(t *testing.T) {
	db, tm := setupTransactionTestDB(t)

	role := &testRole{Name: "ADMIN", CompanyID: 1}

	err := tm.RunInTransaction(context.Background(), func(ctx context.Context) error {
		tx := ctx.Value(constants.TxContextKey).(*gorm.DB)
		return tx.Create(role).Error
	})

	require.NoError(t, err)

	var found testRole
	require.NoError(t, db.First(&found, "name = ?", "ADMIN").Error)
	assert.Equal(t, "ADMIN", found.Name)
}

func TestRunInTransaction_RollsBackOnError(t *testing.T) {
	db, tm := setupTransactionTestDB(t)

	err := tm.RunInTransaction(context.Background(), func(ctx context.Context) error {
		tx := ctx.Value(constants.TxContextKey).(*gorm.DB)
		role := &testRole{Name: "TO_DELETE", CompanyID: 1}
		if e := tx.Create(role).Error; e != nil {
			return e
		}
		return assert.AnError
	})

	assert.Error(t, err)

	var count int64
	db.Model(&testRole{}).Where("name = ?", "TO_DELETE").Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestRunInTransaction_ReusesExistingTx(t *testing.T) {
	db, tm := setupTransactionTestDB(t)

	err := tm.RunInTransaction(context.Background(), func(ctx context.Context) error {
		role := &testRole{Name: "INNER", CompanyID: 1}

		return tm.RunInTransaction(ctx, func(ctx context.Context) error {
			tx := ctx.Value(constants.TxContextKey).(*gorm.DB)
			return tx.Create(role).Error
		})
	})

	require.NoError(t, err)

	var found testRole
	require.NoError(t, db.First(&found, "name = ?", "INNER").Error)
	assert.Equal(t, "INNER", found.Name)
}
