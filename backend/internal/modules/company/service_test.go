package company

import (
	"errors"
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetProfile(t *testing.T) {
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
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Company{
					ID:    1,
					Name:  "Test Co",
					Email: "test@co.com",
				}, nil)
				repo.On("FindPlanByCompanyID", mock.Anything, uint(1)).Return("Free", 5, `{"modules":[]}`, nil)
				cache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "from cache",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(`{"id":1,"name":"Test Co","email":"test@co.com"}`, nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("", errors.New("cache miss"))
				repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, cache, _ := newTestCompanyService()
			tt.setupMocks(repo, cache)

			resp, err := svc.GetProfile(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestService_UpdateProfile(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *UpdateCompanyProfileRequest
		setupMocks func(*mockRepo, *mockCacheProvider)
		wantErr    bool
	}{
		{
			name: "success",
			req: &UpdateCompanyProfileRequest{
				Name:    "Updated Co",
				Address: "456 New St",
				Email:   "new@co.com",
			},
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Company{
					ID:    1,
					Name:  "Test Co",
					Email: "test@co.com",
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*company.Company")).Return(nil)
				cache.On("Del", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "company not found",
			req: &UpdateCompanyProfileRequest{
				Name: "Updated Co",
			},
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, cache, _ := newTestCompanyService()
			tt.setupMocks(repo, cache)

			err := svc.UpdateProfile(ctx, tt.req, nil)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
