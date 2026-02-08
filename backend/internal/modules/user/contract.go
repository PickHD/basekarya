package user

import (
	"context"
	"mime/multipart"

	"gorm.io/gorm"
)

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type StorageProvider interface {
	UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error)
}

type LeaveBalanceGenerator interface {
	GenerateInitialBalance(ctx context.Context, tx interface{}, employeeID uint) (*gorm.DB, error)
}
