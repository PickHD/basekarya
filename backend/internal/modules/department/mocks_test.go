package department

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindAll(ctx context.Context) ([]Department, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Department), args.Error(1)
}
func (m *mockRepo) FindByID(ctx context.Context, id uint) (*Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Department), args.Error(1)
}
func (m *mockRepo) FindByName(ctx context.Context, name string) (*Department, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Department), args.Error(1)
}
func (m *mockRepo) Create(ctx context.Context, dept *Department) error {
	return m.Called(ctx, dept).Error(0)
}
func (m *mockRepo) Update(ctx context.Context, dept *Department) error {
	return m.Called(ctx, dept).Error(0)
}
func (m *mockRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockRepo) CountEmployees(ctx context.Context, departmentID uint) (int64, error) {
	args := m.Called(ctx, departmentID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockRepo) ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error) {
	args := m.Called(ctx, name, excludeID)
	return args.Bool(0), args.Error(1)
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

type mockService struct{ mock.Mock }

func (m *mockService) GetAll(ctx context.Context) ([]LookupResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LookupResponse), args.Error(1)
}
func (m *mockService) GetByID(ctx context.Context, id uint) (*LookupResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LookupResponse), args.Error(1)
}
func (m *mockService) Create(ctx context.Context, req *CreateDepartmentRequest) (*LookupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LookupResponse), args.Error(1)
}
func (m *mockService) Update(ctx context.Context, id uint, req *UpdateDepartmentRequest) (*LookupResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LookupResponse), args.Error(1)
}
func (m *mockService) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
