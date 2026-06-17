package bpjs

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) FindActiveByType(ctx context.Context, bpjsType string, effectiveDate time.Time) ([]BPJSRateConfig, error) {
	args := m.Called(ctx, bpjsType, effectiveDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]BPJSRateConfig), args.Error(1)
}

func (m *mockRepo) FindAllActive(ctx context.Context, effectiveDate time.Time) ([]BPJSRateConfig, error) {
	args := m.Called(ctx, effectiveDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]BPJSRateConfig), args.Error(1)
}

func (m *mockRepo) Create(ctx context.Context, config *BPJSRateConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*BPJSRateConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BPJSRateConfig), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, config *BPJSRateConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *mockRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) List(ctx context.Context, filter BPJSRateConfigFilter) ([]BPJSRateConfig, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]BPJSRateConfig), args.Get(1).(int64), args.Error(2)
}
