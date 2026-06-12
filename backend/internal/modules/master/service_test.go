package master

import (
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetAllShifts(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		setupMocks func(*mockRepo, *mockCacheProvider)
		wantErr    bool
	}{
		{
			name: "from db on cache miss",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("", redis.Nil)
				repo.On("FindAllShifts", mock.Anything).Return([]Shift{
					{ID: 1, Name: "Day"},
				}, nil)
				cache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "from cache",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(`[{"id":1,"name":"Day"}]`, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, cache := newTestMasterService()
			tt.setupMocks(repo, cache)

			_, err := svc.GetAllShifts(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetAllLeaveTypes(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		setupMocks func(*mockRepo, *mockCacheProvider)
		wantErr    bool
	}{
		{
			name: "from db on cache miss",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("", redis.Nil)
				repo.On("FindAllLeaveTypes", mock.Anything).Return([]LeaveType{
					{ID: 1, Name: "Annual", DefaultQuota: 12, IsDeducted: true},
				}, nil)
				cache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "from cache",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(`[{"id":1,"name":"Annual","default_quota":12,"is_deducted":true}]`, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, cache := newTestMasterService()
			tt.setupMocks(repo, cache)

			results, err := svc.GetAllLeaveTypes(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, results)
			}
		})
	}
}
