package infrastructure

import (
	"context"
	"hris-backend/pkg/constants"

	"gorm.io/gorm"
)

type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type gormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) TransactionManager {
	return &gormTransactionManager{db}
}

func (t *gormTransactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// check if already in transaction
	if _, ok := ctx.Value(constants.TxContextKey).(*gorm.DB); ok {
		return fn(ctx)
	}

	return t.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, constants.TxContextKey, tx)
		return fn(txCtx)
	})
}
