package reimbursement

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestReimbursementService() (Service, *mockRepo, *mockStorage, *mockNotification, *mockUserProvider, *testutil.MockTransactionManager, *mockExcel) {
	repo := new(mockRepo)
	storage := new(mockStorage)
	notif := new(mockNotification)
	userProv := new(mockUserProvider)
	tm := testutil.NewMockTransactionManager()
	excel := new(mockExcel)

	svc := NewService(repo, storage, notif, userProv, tm, excel)
	return svc, repo, storage, notif, userProv, tm, excel
}

func TestService_Create(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)
	fileHeader, _ := testutil.CreateMultipartFileHeader("receipt.jpg", "fake content")

	tests := []struct {
		name       string
		req        *ReimbursementRequest
		setupMocks func(*mockRepo, *mockStorage, *mockNotification, *mockUserProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req: &ReimbursementRequest{
				UserID:      1,
				Title:       "Office Supplies",
				Description: "Purchased supplies",
				Amount:      50000,
				Date:        "2026-01-15",
				File:        fileHeader,
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {
				storage.On("UploadFileMultipart", mock.Anything, mock.Anything, mock.Anything).Return("https://storage.example.com/file.jpg", nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*reimbursement.Reimbursement")).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_REIMBURSEMENT)).Return([]uint{10, 11}, nil)
				notif.On("BlastNotification", mock.Anything, []uint{10, 11}, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error invalid user id",
			req: &ReimbursementRequest{
				UserID: 0,
				Title:  "Test",
				Amount: 50000,
				Date:   "2026-01-15",
				File:   fileHeader,
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {},
			wantErr:    true,
			errMsg:     "user id is invalid",
		},
		{
			name: "error upload fails",
			req: &ReimbursementRequest{
				UserID: 1,
				Title:  "Test",
				Amount: 50000,
				Date:   "2026-01-15",
				File:   fileHeader,
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {
				storage.On("UploadFileMultipart", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("upload failed"))
			},
			wantErr: true,
		},
		{
			name: "error invalid date format",
			req: &ReimbursementRequest{
				UserID: 1,
				Title:  "Test",
				Amount: 50000,
				Date:   "not-a-date",
				File:   fileHeader,
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {
				storage.On("UploadFileMultipart", mock.Anything, mock.Anything, mock.Anything).Return("https://storage.example.com/file.jpg", nil)
			},
			wantErr: true,
			errMsg:  "invalid date format",
		},
		{
			name: "error repo create fails",
			req: &ReimbursementRequest{
				UserID: 1,
				Title:  "Test",
				Amount: 50000,
				Date:   "2026-01-15",
				File:   fileHeader,
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {
				storage.On("UploadFileMultipart", mock.Anything, mock.Anything, mock.Anything).Return("https://storage.example.com/file.jpg", nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*reimbursement.Reimbursement")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name: "error find approval users fails",
			req: &ReimbursementRequest{
				UserID: 1,
				Title:  "Test",
				Amount: 50000,
				Date:   "2026-01-15",
				File:   fileHeader,
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {
				storage.On("UploadFileMultipart", mock.Anything, mock.Anything, mock.Anything).Return("https://storage.example.com/file.jpg", nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*reimbursement.Reimbursement")).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_REIMBURSEMENT)).Return([]uint(nil), errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, storage, notif, userProv, _, _ := newTestReimbursementService()
			tt.setupMocks(repo, storage, notif, userProv)

			err := svc.Create(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetReimburseDetail(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		id         uint
		setupMocks func(*mockRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Reimbursement{
					ID:              1,
					Title:           "Office Supplies",
					Description:     "Purchased supplies",
					Amount:          50000,
					DateOfExpense:   time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
					ProofFileURL:    "https://storage.example.com/file.jpg",
					Status:          constants.ReimbursementStatusPending,
					RejectionReason: sql.NullString{Valid: false},
					User:            user.User{ID: 1, Username: "john"},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name: "error user not found",
			id:   2,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&Reimbursement{
					ID:              2,
					Title:           "Test",
					Status:          constants.ReimbursementStatusPending,
					RejectionReason: sql.NullString{Valid: false},
					User:            user.User{ID: 0},
				}, nil)
			},
			wantErr: true,
			errMsg:  "data user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _ := newTestReimbursementService()
			tt.setupMocks(repo)

			resp, err := svc.GetReimburseDetail(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.id, resp.ID)
				assert.Equal(t, "john", resp.RequesterName)
				assert.Equal(t, "PENDING", resp.Status)
			}
		})
	}
}

func TestService_GetReimbursements(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     ReimbursementFilter
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:   "success with data",
			filter: ReimbursementFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("reimbursement.ReimbursementFilter")).
					Return([]Reimbursement{
						{
							ID:            1,
							Title:         "Test",
							Amount:        50000,
							DateOfExpense: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
							ProofFileURL:  "https://storage.example.com/file.jpg",
							Status:        constants.ReimbursementStatusPending,
						},
					}, int64(1), nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "success empty list",
			filter: ReimbursementFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("reimbursement.ReimbursementFilter")).
					Return([]Reimbursement{}, int64(0), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "repo returns error returns empty",
			filter: ReimbursementFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("reimbursement.ReimbursementFilter")).
					Return([]Reimbursement(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _ := newTestReimbursementService()
			tt.setupMocks(repo)

			list, meta, err := svc.GetReimbursements(ctx, tt.filter)

			require.NoError(t, err)
			assert.Len(t, list, tt.wantLen)
			if tt.wantLen > 0 {
				assert.NotNil(t, meta)
			}
		})
	}
}

func TestService_ProcessAction(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *ActionRequest
		setupMocks func(*mockRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "approve success",
			req: &ActionRequest{
				ID:           1,
				SuperAdminID: 10,
				Action:       string(constants.ReimbursementActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Reimbursement{
					ID:     1,
					Status: constants.ReimbursementStatusPending,
					User:   user.User{ID: 1},
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*reimbursement.Reimbursement")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "reject success",
			req: &ActionRequest{
				ID:              2,
				SuperAdminID:    10,
				Action:          string(constants.ReimbursementActionReject),
				RejectionReason: "Not eligible",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&Reimbursement{
					ID:     2,
					Status: constants.ReimbursementStatusPending,
					User:   user.User{ID: 1},
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*reimbursement.Reimbursement")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not pending",
			req: &ActionRequest{
				ID:           3,
				SuperAdminID: 10,
				Action:       string(constants.ReimbursementActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(3)).Return(&Reimbursement{
					ID:     3,
					Status: constants.ReimbursementStatusApproved,
					User:   user.User{ID: 1},
				}, nil)
			},
			wantErr: true,
			errMsg:  "cannot process reimburstment with status APPROVED",
		},
		{
			name: "error rejection reason required",
			req: &ActionRequest{
				ID:              4,
				SuperAdminID:    10,
				Action:          string(constants.ReimbursementActionReject),
				RejectionReason: "",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(4)).Return(&Reimbursement{
					ID:     4,
					Status: constants.ReimbursementStatusPending,
					User:   user.User{ID: 1},
				}, nil)
			},
			wantErr: true,
			errMsg:  "rejection reason is required",
		},
		{
			name: "error not found",
			req: &ActionRequest{
				ID:           99,
				SuperAdminID: 10,
				Action:       string(constants.ReimbursementActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name: "error invalid action",
			req: &ActionRequest{
				ID:           5,
				SuperAdminID: 10,
				Action:       "INVALID",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(5)).Return(&Reimbursement{
					ID:     5,
					Status: constants.ReimbursementStatusPending,
					User:   user.User{ID: 1},
				}, nil)
			},
			wantErr: true,
			errMsg:  "invalid action: INVALID",
		},
		{
			name: "error update fails",
			req: &ActionRequest{
				ID:           6,
				SuperAdminID: 10,
				Action:       string(constants.ReimbursementActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(6)).Return(&Reimbursement{
					ID:     6,
					Status: constants.ReimbursementStatusPending,
					User:   user.User{ID: 1},
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*reimbursement.Reimbursement")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, notif, _, _, _ := newTestReimbursementService()
			tt.setupMocks(repo)

			if !tt.wantErr {
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}

			err := svc.ProcessAction(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Export(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     ReimbursementFilter
		setupMocks func(*mockRepo, *mockExcel)
		wantErr    bool
	}{
		{
			name:   "success",
			filter: ReimbursementFilter{Status: "PENDING"},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("reimbursement.ReimbursementFilter")).
					Return([]Reimbursement{
						{
							ID:            1,
							UserID:        1,
							User:          user.User{ID: 1, Username: "john"},
							Title:         "Office Supplies",
							Amount:        50000,
							DateOfExpense: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
							Status:        constants.ReimbursementStatusPending,
							Description:   "Purchased supplies",
							CreatedAt:     time.Now(),
						},
					}, int64(1), nil)
				excel.On("GenerateSimpleExcel", "Reimbursements", mock.Anything, mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "error repo fails",
			filter: ReimbursementFilter{},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("reimbursement.ReimbursementFilter")).
					Return([]Reimbursement(nil), int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, excel := newTestReimbursementService()
			tt.setupMocks(repo, excel)

			data, err := svc.Export(ctx, tt.filter)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, data)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, data)
			}
		})
	}
}
