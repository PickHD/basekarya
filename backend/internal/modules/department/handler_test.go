package department

import (
	"errors"
	"net/http"
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mockService)
		svc.On("GetAll", mock.Anything).Return([]LookupResponse{
			{ID: 1, Name: "IT"},
		}, nil)

		handler := NewHandler(svc)
		at := testutil.NewAPITest(t, http.MethodGet, "/api/departments", nil)
		rec, err := at.Execute(handler.GetAll)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error", func(t *testing.T) {
		svc := new(mockService)
		svc.On("GetAll", mock.Anything).Return(nil, errors.New("db error"))

		handler := NewHandler(svc)
		at := testutil.NewAPITest(t, http.MethodGet, "/api/departments", nil)
		rec, err := at.Execute(handler.GetAll)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mockService)
		svc.On("GetByID", mock.Anything, uint(1)).Return(&LookupResponse{ID: 1, Name: "IT"}, nil)

		handler := NewHandler(svc)
		at := testutil.NewAPITest(t, http.MethodGet, "/api/departments/1", nil)
		at.WithPathParams(map[string]string{"id": "1"})
		rec, err := at.Execute(handler.GetByID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(mockService)
		svc.On("GetByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))

		handler := NewHandler(svc)
		at := testutil.NewAPITest(t, http.MethodGet, "/api/departments/99", nil)
		at.WithPathParams(map[string]string{"id": "99"})
		rec, err := at.Execute(handler.GetByID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mockService)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*department.CreateDepartmentRequest")).Return(&LookupResponse{ID: 2, Name: "Finance"}, nil)

		handler := NewHandler(svc)
		at := testutil.NewAPITest(t, http.MethodPost, "/api/departments", CreateDepartmentRequest{Name: "Finance"})
		rec, err := at.Execute(handler.Create)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("duplicate", func(t *testing.T) {
		svc := new(mockService)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*department.CreateDepartmentRequest")).Return(nil, errors.New("department name already exists"))

		handler := NewHandler(svc)
		at := testutil.NewAPITest(t, http.MethodPost, "/api/departments", CreateDepartmentRequest{Name: "IT"})
		rec, err := at.Execute(handler.Create)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mockService)
		svc.On("Update", mock.Anything, uint(1), mock.AnythingOfType("*department.UpdateDepartmentRequest")).Return(&LookupResponse{ID: 1, Name: "Tech"}, nil)

		handler := NewHandler(svc)
		at := testutil.NewAPITest(t, http.MethodPut, "/api/departments/1", UpdateDepartmentRequest{Name: "Tech"})
		at.WithPathParams(map[string]string{"id": "1"})
		rec, err := at.Execute(handler.Update)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mockService)
		svc.On("Delete", mock.Anything, uint(1)).Return(nil)

		handler := NewHandler(svc)
		at := testutil.NewAPITest(t, http.MethodDelete, "/api/departments/1", nil)
		at.WithPathParams(map[string]string{"id": "1"})
		rec, err := at.Execute(handler.Delete)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(mockService)
		svc.On("Delete", mock.Anything, uint(99)).Return(errors.New("department not found"))

		handler := NewHandler(svc)
		at := testutil.NewAPITest(t, http.MethodDelete, "/api/departments/99", nil)
		at.WithPathParams(map[string]string{"id": "99"})
		rec, err := at.Execute(handler.Delete)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
