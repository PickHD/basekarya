package reimbursement

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("Create", mock.Anything, mock.AnythingOfType("*reimbursement.ReimbursementRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.CREATE_REIMBURSEMENT},
				})
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("Create", mock.Anything, mock.AnythingOfType("*reimbursement.ReimbursementRequest")).Return(errors.New("upload failed"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.CREATE_REIMBURSEMENT},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			writer.WriteField("title", "Office Supplies")
			writer.WriteField("description", "Purchased supplies")
			writer.WriteField("amount", "50000")
			writer.WriteField("date", "2026-01-15")
			part, _ := writer.CreateFormFile("file", "receipt.jpg")
			part.Write([]byte("fake content"))
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/api/reimbursements", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e := echo.New()
			e.Validator = utils.NewValidator()
			ctx := e.NewContext(req, rec)

			at := &testutil.APITest{Echo: e, Req: req, Rec: rec, Context: ctx}
			tt.setupContext(at)

			err := handler.Create(at.Context)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tt.wantStatus < 400 {
				assert.Nil(t, resp["error"])
			}
		})
	}
}

func TestHandler_GetAll(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:        "success",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetReimbursements", mock.Anything, mock.AnythingOfType("reimbursement.ReimbursementFilter")).
					Return([]ReimbursementListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_REIMBURSEMENT},
				})
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/reimbursements"+tt.queryParams, nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetAll)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetDetail(t *testing.T) {
	tests := []struct {
		name       string
		pathParams map[string]string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name:       "success",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("GetReimburseDetail", mock.Anything, uint(1)).Return(&ReimbursementDetailResponse{ID: 1}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetReimburseDetail", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/reimbursements/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_ProcessAction(t *testing.T) {
	tests := []struct {
		name         string
		pathParams   map[string]string
		body         interface{}
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:       "approve success",
			pathParams: map[string]string{"id": "1"},
			body: ActionRequest{
				Action:          "APPROVE",
				RejectionReason: "",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*reimbursement.ActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.APPROVAL_REIMBURSEMENT},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "2"},
			body: ActionRequest{
				Action:          "REJECT",
				RejectionReason: "Not eligible",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*reimbursement.ActionRequest")).Return(errors.New("cannot process"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.APPROVAL_REIMBURSEMENT},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/reimbursements/:id/action", tt.body)
			at.WithPathParams(tt.pathParams)
			if tt.setupContext != nil {
				tt.setupContext(at)
			}

			rec, err := at.Execute(handler.ProcessAction)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_Export(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:        "success",
			queryParams: "?status=PENDING",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("reimbursement.ReimbursementFilter")).
					Return([]byte("fake-excel-content"), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.EXPORT_REIMBURSEMENT},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("reimbursement.ReimbursementFilter")).
					Return(nil, errors.New("export failed"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.EXPORT_REIMBURSEMENT},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/reimbursements/export"+tt.queryParams, nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.Export)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
