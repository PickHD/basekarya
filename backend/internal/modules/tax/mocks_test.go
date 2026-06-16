package tax

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) FindTERBrackets(ctx context.Context, category string, effectiveDate time.Time) ([]TERBracket, error) {
	args := m.Called(ctx, category, effectiveDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]TERBracket), args.Error(1)
}

func (m *mockRepo) CreateTERBracket(ctx context.Context, bracket *TERBracket) error {
	return m.Called(ctx, bracket).Error(0)
}

func (m *mockRepo) FindTERBracketByID(ctx context.Context, id uint) (*TERBracket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TERBracket), args.Error(1)
}

func (m *mockRepo) UpdateTERBracket(ctx context.Context, bracket *TERBracket) error {
	return m.Called(ctx, bracket).Error(0)
}

func (m *mockRepo) DeleteTERBracket(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) ListTERBrackets(ctx context.Context, filter TERBracketFilter) ([]TERBracket, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]TERBracket), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) FindPTKPByYear(ctx context.Context, year int) ([]PTKPConfig, error) {
	args := m.Called(ctx, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]PTKPConfig), args.Error(1)
}

func (m *mockRepo) CreatePTKPConfig(ctx context.Context, ptkp *PTKPConfig) error {
	return m.Called(ctx, ptkp).Error(0)
}

func (m *mockRepo) FindPTKPConfigByID(ctx context.Context, id uint) (*PTKPConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PTKPConfig), args.Error(1)
}

func (m *mockRepo) UpdatePTKPConfig(ctx context.Context, ptkp *PTKPConfig) error {
	return m.Called(ctx, ptkp).Error(0)
}

func (m *mockRepo) DeletePTKPConfig(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) ListPTKPConfigs(ctx context.Context, year int) ([]PTKPConfig, int64, error) {
	args := m.Called(ctx, year)
	return args.Get(0).([]PTKPConfig), args.Get(1).(int64), args.Error(2)
}
