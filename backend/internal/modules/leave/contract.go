package leave

import (
	"context"
	"io"
)

type StorageProvider interface {
	UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
}
