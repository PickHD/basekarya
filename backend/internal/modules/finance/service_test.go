package finance

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestFinanceService() (Service, *mockRepo, *mockNotificationProvider, *mockUserProvider, *testutil.MockTransactionManager, *mockExcel) {
	repo := new(mockRepo)
	notif := new(mockNotificationProvider)
	userProv := new(mockUserProvider)
	tm := testutil.NewMockTransactionManager()
	excel := new(mockExcel)

	svc := NewService(repo, notif, userProv, tm, excel)
	return svc, repo, notif, userProv, tm, excel
}

func TestService_CreateTransaction(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CreateTransactionRequest
		setupMocks func(*mockRepo, *mockUserProvider, *mockNotificationProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req: &CreateTransactionRequest{
				CreatedBy:         1,
				FinanceCategoryID: 1,
				Type:              "INCOME",
				Amount:            5000000,
				TransactionDate:   "2026-01-15",
			},
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, notif *mockNotificationProvider) {
				repo.On("FindCategoryByID", mock.Anything, uint(1)).Return(&FinanceCategory{ID: 1, Name: "Salary"}, nil)
				repo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*finance.FinanceTransaction")).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_FINANCE)).Return([]uint{10, 11}, nil)
				notif.On("BlastNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error user not found",
			req: &CreateTransactionRequest{
				CreatedBy:         0,
				FinanceCategoryID: 1,
				Type:              "INCOME",
				Amount:            5000000,
				TransactionDate:   "2026-01-15",
			},
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, notif *mockNotificationProvider) {},
			wantErr:    true,
			errMsg:     "user not found",
		},
		{
			name: "error invalid date format",
			req: &CreateTransactionRequest{
				CreatedBy:         1,
				FinanceCategoryID: 1,
				Type:              "INCOME",
				Amount:            5000000,
				TransactionDate:   "not-a-date",
			},
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, notif *mockNotificationProvider) {},
			wantErr:    true,
			errMsg:     "invalid transaction_date format, use YYYY-MM-DD",
		},
		{
			name: "error category not found",
			req: &CreateTransactionRequest{
				CreatedBy:         1,
				FinanceCategoryID: 99,
				Type:              "INCOME",
				Amount:            5000000,
				TransactionDate:   "2026-01-15",
			},
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, notif *mockNotificationProvider) {
				repo.On("FindCategoryByID", mock.Anything, uint(99)).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errMsg:  "finance category not found",
		},
		{
			name: "error create transaction fails",
			req: &CreateTransactionRequest{
				CreatedBy:         1,
				FinanceCategoryID: 1,
				Type:              "INCOME",
				Amount:            5000000,
				TransactionDate:   "2026-01-15",
			},
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, notif *mockNotificationProvider) {
				repo.On("FindCategoryByID", mock.Anything, uint(1)).Return(&FinanceCategory{ID: 1}, nil)
				repo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*finance.FinanceTransaction")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, notif, userProv, _, _ := newTestFinanceService()
			tt.setupMocks(repo, userProv, notif)

			err := svc.CreateTransaction(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetTransactionDetail(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		id         uint
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(repo *mockRepo) {
				approverID := uint(2)
				repo.On("FindTransactionByID", mock.Anything, uint(1)).Return(&FinanceTransaction{
					ID:        1,
					Type:      constants.FinanceTypeIncome,
					Amount:    5000000,
					Status:    constants.FinanceStatusApproved,
					ApprovedBy: &approverID,
					Creator:   user.User{ID: 1, Employee: &user.Employee{FullName: "John"}},
					Approver:  &user.User{ID: 2, Employee: &user.Employee{FullName: "Admin"}},
					FinanceCategory: FinanceCategory{Name: "Salary", Type: constants.FinanceTypeIncome},
					Description:     sql.NullString{String: "Monthly", Valid: true},
					ReferenceNumber: sql.NullString{String: "REF-001", Valid: true},
					RejectionReason: sql.NullString{Valid: false},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindTransactionByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestFinanceService()
			tt.setupMocks(repo)

			resp, err := svc.GetTransactionDetail(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.id, resp.ID)
				assert.Equal(t, "John", resp.CreatorName)
				assert.Equal(t, "Admin", resp.ApproverName)
			}
		})
	}
}

func TestService_GetTransactions(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     TransactionFilter
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:   "success with data",
			filter: TransactionFilter{Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllTransactions", mock.Anything, mock.AnythingOfType("finance.TransactionFilter")).Return([]FinanceTransaction{
					{
						ID: 1, Type: constants.FinanceTypeIncome, Amount: 5000000, Status: constants.FinanceStatusPending,
						Creator:         user.User{ID: 1, Employee: &user.Employee{FullName: "John"}},
						FinanceCategory: FinanceCategory{Name: "Salary"},
						ReferenceNumber: sql.NullString{Valid: false},
					},
				}, (*response.Cursor)(nil), nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "success empty",
			filter: TransactionFilter{Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllTransactions", mock.Anything, mock.AnythingOfType("finance.TransactionFilter")).Return([]FinanceTransaction{}, (*response.Cursor)(nil), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "repo error returns empty",
			filter: TransactionFilter{Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllTransactions", mock.Anything, mock.AnythingOfType("finance.TransactionFilter")).Return([]FinanceTransaction(nil), (*response.Cursor)(nil), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestFinanceService()
			tt.setupMocks(repo)

			list, meta, err := svc.GetTransactions(ctx, tt.filter)

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
		setupMocks func(*mockRepo, *mockNotificationProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "approve success",
			req: &ActionRequest{
				ID:           1,
				SuperAdminID: 2,
				Action:       string(constants.FinanceActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider) {
				repo.On("FindTransactionByID", mock.Anything, uint(1)).Return(&FinanceTransaction{
					ID: 1, CreatedBy: 1, Status: constants.FinanceStatusPending, Type: constants.FinanceTypeIncome,
					RejectionReason: sql.NullString{Valid: false},
				}, nil)
				repo.On("UpdateTransaction", mock.Anything, mock.AnythingOfType("*finance.FinanceTransaction")).Return(nil)
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "reject success",
			req: &ActionRequest{
				ID:              2,
				SuperAdminID:    2,
				Action:          string(constants.FinanceActionReject),
				RejectionReason: "Invalid",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider) {
				repo.On("FindTransactionByID", mock.Anything, uint(2)).Return(&FinanceTransaction{
					ID: 2, CreatedBy: 1, Status: constants.FinanceStatusPending, Type: constants.FinanceTypeIncome,
					RejectionReason: sql.NullString{Valid: false},
				}, nil)
				repo.On("UpdateTransaction", mock.Anything, mock.AnythingOfType("*finance.FinanceTransaction")).Return(nil)
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not pending",
			req: &ActionRequest{
				ID:           3,
				SuperAdminID: 2,
				Action:       string(constants.FinanceActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider) {
				repo.On("FindTransactionByID", mock.Anything, uint(3)).Return(&FinanceTransaction{
					ID: 3, Status: constants.FinanceStatusApproved, Type: constants.FinanceTypeIncome,
					RejectionReason: sql.NullString{Valid: false},
				}, nil)
			},
			wantErr: true,
			errMsg:  "cannot process transaction with status APPROVED",
		},
		{
			name: "error rejection reason required",
			req: &ActionRequest{
				ID:           4,
				SuperAdminID: 2,
				Action:       string(constants.FinanceActionReject),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider) {
				repo.On("FindTransactionByID", mock.Anything, uint(4)).Return(&FinanceTransaction{
					ID: 4, Status: constants.FinanceStatusPending, Type: constants.FinanceTypeIncome,
					RejectionReason: sql.NullString{Valid: false},
				}, nil)
			},
			wantErr: true,
			errMsg:  "rejection reason is required",
		},
		{
			name: "error invalid action",
			req: &ActionRequest{
				ID:           5,
				SuperAdminID: 2,
				Action:       "INVALID",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider) {
				repo.On("FindTransactionByID", mock.Anything, uint(5)).Return(&FinanceTransaction{
					ID: 5, Status: constants.FinanceStatusPending, Type: constants.FinanceTypeIncome,
					RejectionReason: sql.NullString{Valid: false},
				}, nil)
			},
			wantErr: true,
			errMsg:  "invalid action: INVALID",
		},
		{
			name: "error transaction not found",
			req: &ActionRequest{
				ID:           99,
				SuperAdminID: 2,
				Action:       string(constants.FinanceActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider) {
				repo.On("FindTransactionByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, notif, _, _, _ := newTestFinanceService()
			tt.setupMocks(repo, notif)

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

func TestService_ExportTransactions(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     TransactionFilter
		setupMocks func(*mockRepo, *mockExcel)
		wantErr    bool
	}{
		{
			name:   "success",
			filter: TransactionFilter{Limit: 0},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAllTransactions", mock.Anything, mock.AnythingOfType("finance.TransactionFilter")).Return([]FinanceTransaction{
					{
						ID: 1, Type: constants.FinanceTypeIncome, Amount: 5000000, Status: constants.FinanceStatusPending,
						Creator:         user.User{ID: 1, Employee: &user.Employee{FullName: "John"}},
						FinanceCategory: FinanceCategory{Name: "Salary"},
						TransactionDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
						ReferenceNumber: sql.NullString{Valid: false},
						CreatedAt:       time.Now(),
					},
				}, (*response.Cursor)(nil), nil)
				excel.On("GenerateSimpleExcel", "Finance Transactions", mock.Anything, mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "error repo fails",
			filter: TransactionFilter{Limit: 0},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAllTransactions", mock.Anything, mock.AnythingOfType("finance.TransactionFilter")).Return([]FinanceTransaction(nil), (*response.Cursor)(nil), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, excel := newTestFinanceService()
			tt.setupMocks(repo, excel)

			data, err := svc.ExportTransactions(ctx, tt.filter)

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

func TestService_CreateCategory(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CategoryRequest
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success",
			req: &CategoryRequest{
				Name: "Bonus",
				Type: "INCOME",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateCategory", mock.Anything, mock.AnythingOfType("*finance.FinanceCategory")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success with description",
			req: &CategoryRequest{
				Name:        "Rent",
				Type:        "EXPENSE",
				Description: "Office rent",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateCategory", mock.Anything, mock.AnythingOfType("*finance.FinanceCategory")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestFinanceService()
			tt.setupMocks(repo)

			err := svc.CreateCategory(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetCategories(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		catType    string
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:    "success",
			catType: "",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllCategories", mock.Anything, "").Return([]FinanceCategory{
					{ID: 1, Name: "Salary", Type: constants.FinanceTypeIncome, Description: sql.NullString{Valid: false}},
				}, nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "error",
			catType: "",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllCategories", mock.Anything, "").Return([]FinanceCategory(nil), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestFinanceService()
			tt.setupMocks(repo)

			list, err := svc.GetCategories(ctx, tt.catType)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, list, tt.wantLen)
			}
		})
	}
}

func TestService_UpdateCategory(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		id         uint
		req        *CategoryRequest
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success",
			id:   1,
			req: &CategoryRequest{
				Name: "Updated Salary",
				Type: "INCOME",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindCategoryByID", mock.Anything, uint(1)).Return(&FinanceCategory{
					ID: 1, Name: "Salary", Type: constants.FinanceTypeIncome,
					Description: sql.NullString{Valid: false},
				}, nil)
				repo.On("UpdateCategory", mock.Anything, mock.AnythingOfType("*finance.FinanceCategory")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   99,
			req: &CategoryRequest{
				Name: "Test",
				Type: "INCOME",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindCategoryByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestFinanceService()
			tt.setupMocks(repo)

			err := svc.UpdateCategory(ctx, tt.id, tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_DeleteCategory(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		id         uint
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(repo *mockRepo) {
				repo.On("DeleteCategory", mock.Anything, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("DeleteCategory", mock.Anything, uint(99)).Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestFinanceService()
			tt.setupMocks(repo)

			err := svc.DeleteCategory(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetDashboard(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		startDate  string
		endDate    string
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name:      "success",
			startDate: "2026-01-01",
			endDate:   "2026-12-31",
			setupMocks: func(repo *mockRepo) {
				repo.On("GetDashboardSummary", mock.Anything, "2026-01-01", "2026-12-31").Return(&DashboardResponse{
					TotalIncome:  10000000,
					TotalExpense: 5000000,
					NetBalance:   5000000,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:      "error",
			startDate: "",
			endDate:   "",
			setupMocks: func(repo *mockRepo) {
				repo.On("GetDashboardSummary", mock.Anything, "", "").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestFinanceService()
			tt.setupMocks(repo)

			resp, err := svc.GetDashboard(ctx, tt.startDate, tt.endDate)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}
