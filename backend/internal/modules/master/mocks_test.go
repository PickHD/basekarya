package master

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindAllShifts(ctx context.Context) ([]Shift, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Shift), args.Error(1)
}

func (m *mockRepo) FindAllLeaveTypes(ctx context.Context) ([]LeaveType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LeaveType), args.Error(1)
}

func (m *mockRepo) FindShiftByName(ctx context.Context, name string) (*Shift, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Shift), args.Error(1)
}

func (m *mockRepo) SeedDefaults(ctx context.Context, companyID uint) error {
	return m.Called(ctx, companyID).Error(0)
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

func (m *mockService) GetAllShifts(ctx context.Context) ([]LookupResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LookupResponse), args.Error(1)
}

func (m *mockService) GetAllLeaveTypes(ctx context.Context) ([]LookupLeaveTypeResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]LookupLeaveTypeResponse), args.Error(1)
}

func newTestMasterService() (Service, *mockRepo, *mockCacheProvider) {
	repo := new(mockRepo)
	cache := new(mockCacheProvider)
	return NewService(repo, cache), repo, cache
}
