package department

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetAll(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	cache := new(mockCacheProvider)

	cache.On("Get", mock.Anything, "department:all").Return("", redis.Nil)
	repo.On("FindAll", mock.Anything).Return([]Department{
		{ID: 1, Name: "Engineering"},
	}, nil)
	cache.On("Set", mock.Anything, "department:all", mock.Anything, mock.Anything).Return(nil)

	svc := NewService(repo, cache)
	result, err := svc.GetAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Engineering", result[0].Name)
	repo.AssertExpectations(t)
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	cache := new(mockCacheProvider)

	repo.On("FindByID", mock.Anything, uint(1)).Return(&Department{ID: 1, Name: "IT"}, nil)

	svc := NewService(repo, cache)
	result, err := svc.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "IT", result.Name)
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	cache := new(mockCacheProvider)

	req := &CreateDepartmentRequest{Name: "Finance"}
	repo.On("ExistsByName", mock.Anything, "Finance", uint(0)).Return(false, nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*department.Department")).Return(nil)
	cache.On("Del", mock.Anything, "department:all").Return(nil)

	svc := NewService(repo, cache)
	result, err := svc.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "Finance", result.Name)
}

func TestService_Create_Duplicate(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	cache := new(mockCacheProvider)

	req := &CreateDepartmentRequest{Name: "Finance"}
	repo.On("ExistsByName", mock.Anything, "Finance", uint(0)).Return(true, nil)

	svc := NewService(repo, cache)
	_, err := svc.Create(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	cache := new(mockCacheProvider)

	req := &UpdateDepartmentRequest{Name: "Tech"}
	repo.On("FindByID", mock.Anything, uint(1)).Return(&Department{ID: 1, Name: "Old"}, nil)
	repo.On("ExistsByName", mock.Anything, "Tech", uint(1)).Return(false, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*department.Department")).Return(nil)
	cache.On("Del", mock.Anything, "department:all").Return(nil)

	svc := NewService(repo, cache)
	result, err := svc.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "Tech", result.Name)
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	cache := new(mockCacheProvider)

	repo.On("FindByID", mock.Anything, uint(1)).Return(&Department{ID: 1, Name: "IT"}, nil)
	repo.On("Delete", mock.Anything, uint(1)).Return(nil)
	cache.On("Del", mock.Anything, "department:all").Return(nil)

	svc := NewService(repo, cache)
	err := svc.Delete(ctx, 1)

	assert.NoError(t, err)
}
