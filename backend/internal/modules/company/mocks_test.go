package company

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Company), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, company *Company) error {
	return m.Called(ctx, company).Error(0)
}

func (m *mockRepo) CreateCompany(ctx context.Context, c *Company) error {
	return m.Called(ctx, c).Error(0)
}

func (m *mockRepo) FindPlanIDBySlug(ctx context.Context, slug string) (uint, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(uint), args.Error(1)
}

func (m *mockRepo) FindPlanByCompanyID(ctx context.Context, companyID uint) (string, int, string, error) {
	args := m.Called(ctx, companyID)
	return args.String(0), args.Get(1).(int), args.String(2), args.Error(3)
}

func (m *mockRepo) FindModulesByCompanyID(ctx context.Context, companyID uint) ([]string, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

type mockCacheProvider struct{ mock.Mock }

func (m *mockCacheProvider) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockCacheProvider) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.Called(ctx, key, value, expiration).Error(0)
}

func (m *mockCacheProvider) Del(ctx context.Context, key string) error {
	return m.Called(ctx, key).Error(0)
}

type mockStorageProvider struct{ mock.Mock }

func (m *mockStorageProvider) UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	args := m.Called(ctx, file, objectName)
	return args.String(0), args.Error(1)
}

type mockService struct{ mock.Mock }

func (m *mockService) GetProfile(ctx context.Context) (*CompanyProfileResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CompanyProfileResponse), args.Error(1)
}

func (m *mockService) UpdateProfile(ctx context.Context, req *UpdateCompanyProfileRequest, file *multipart.FileHeader) error {
	return m.Called(ctx, req, file).Error(0)
}

func newTestCompanyService() (Service, *mockRepo, *mockCacheProvider, *mockStorageProvider) {
	repo := new(mockRepo)
	cache := new(mockCacheProvider)
	storage := new(mockStorageProvider)
	return NewService(repo, cache, storage), repo, cache, storage
}
