package reimbursement

import (
	"context"
	"mime/multipart"
)

type StorageProvider interface {
	UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error)
}
