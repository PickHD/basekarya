package testutil

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"time"

	"github.com/stretchr/testify/mock"
)

// === Repository Mocks ===
// These are generic mock helpers. Each module will create its own typed mocks.

// MockStorageProvider mocks infrastructure.StorageProvider.
type MockStorageProvider struct {
	mock.Mock
}

func (m *MockStorageProvider) UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	args := m.Called(ctx, file, objectName)
	return args.String(0), args.Error(1)
}

func (m *MockStorageProvider) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, objectName, reader, size, contentType)
	return args.String(0), args.Error(1)
}

// MockCacheProvider mocks the CacheProvider interface.
type MockCacheProvider struct {
	mock.Mock
}

func (m *MockCacheProvider) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheProvider) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheProvider) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheProvider) FlushDB(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockHasher mocks the Hasher interface.
type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

// MockTokenProvider mocks the TokenProvider interface.
type MockTokenProvider struct {
	mock.Mock
}

func (m *MockTokenProvider) GenerateToken(userID uint, companyID uint, isPlatformAdmin bool, role string, employeeID *uint, permissions []string) (string, error) {
	args := m.Called(userID, companyID, isPlatformAdmin, role, employeeID, permissions)
	return args.String(0), args.Error(1)
}

// MockEmailProvider mocks the EmailProvider interface.
type MockEmailProvider struct {
	mock.Mock
}

func (m *MockEmailProvider) Send(to, subject, htmlBody string) error {
	args := m.Called(to, subject, htmlBody)
	return args.Error(0)
}

func (m *MockEmailProvider) SendWithAttachment(to, subject, htmlBody, fileName string, attachmentBytes []byte) error {
	args := m.Called(to, subject, htmlBody, fileName, attachmentBytes)
	return args.Error(0)
}

// MockNotificationProvider mocks notification.NotificationProvider.
type MockNotificationProvider struct {
	mock.Mock
}

func (m *MockNotificationProvider) SendNotification(ctx context.Context, userID uint, notifType string, title string, message string, relatedID uint) error {
	args := m.Called(ctx, userID, notifType, title, message, relatedID)
	return args.Error(0)
}

func (m *MockNotificationProvider) BlastNotification(ctx context.Context, userIDs []uint, notifType string, title string, message string, relatedID uint) error {
	args := m.Called(ctx, userIDs, notifType, title, message, relatedID)
	return args.Error(0)
}

// MockExcelProvider mocks infrastructure.ExcelProvider.
type MockExcelProvider struct {
	mock.Mock
}

func (m *MockExcelProvider) GenerateSimpleExcel(sheetName string, headers []string, rows [][]interface{}) ([]byte, error) {
	args := m.Called(sheetName, headers, rows)
	return args.Get(0).([]byte), args.Error(1)
}

// MockTransactionManager executes the function directly without a real DB transaction.
type MockTransactionManager struct{}

func NewMockTransactionManager() *MockTransactionManager {
	return &MockTransactionManager{}
}

func (m *MockTransactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// Helper to create a multipart file header for testing.
func CreateMultipartFileHeader(filename, content string) (*multipart.FileHeader, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	part.Write([]byte(content))
	writer.Close()

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(1024)
	if err != nil {
		return nil, err
	}

	return form.File["file"][0], nil
}
