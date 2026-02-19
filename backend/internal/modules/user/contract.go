package user

import (
	"context"
	"mime/multipart"
)

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type StorageProvider interface {
	UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error)
}

type LeaveBalanceGenerator interface {
	GenerateInitialBalance(ctx context.Context, employeeID uint) error
}
