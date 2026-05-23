package loan

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
	"gorm.io/gorm"
)

func newTestLoanService() (Service, *mockRepo, *mockNotification, *mockUserProvider, *testutil.MockTransactionManager, *mockExcel) {
	repo := new(mockRepo)
	notif := new(mockNotification)
	userProv := new(mockUserProvider)
	tm := testutil.NewMockTransactionManager()
	excel := new(mockExcel)

	svc := NewService(repo, notif, userProv, tm, excel)
	return svc, repo, notif, userProv, tm, excel
}

func TestService_Create(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *LoanRequest
		setupMocks func(*mockRepo, *mockNotification, *mockUserProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req: &LoanRequest{
				UserID:            1,
				EmployeeID:        1,
				TotalAmount:       5000000,
				InstallmentAmount: 500000,
				Reason:            "Emergency",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindActiveLoanByUserID", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*loan.Loan")).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_LOAN)).Return([]uint{10, 11}, nil)
				notif.On("BlastNotification", mock.Anything, []uint{10, 11}, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error user has active loan",
			req: &LoanRequest{
				UserID:            1,
				EmployeeID:        1,
				TotalAmount:       5000000,
				InstallmentAmount: 500000,
				Reason:            "Emergency",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindActiveLoanByUserID", mock.Anything, uint(1)).Return(&Loan{ID: 1, Status: constants.LoanStatusApproved}, nil)
			},
			wantErr: true,
			errMsg:  "users still have loan",
		},
		{
			name: "error user not found",
			req: &LoanRequest{
				UserID:     0,
				EmployeeID: 0,
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
			},
			wantErr: true,
			errMsg:  "user not found",
		},
		{
			name: "error exceeds maximum amount",
			req: &LoanRequest{
				UserID:            1,
				EmployeeID:        1,
				TotalAmount:       constants.LoanMaximumTotalAmount + 1,
				InstallmentAmount: 500000,
				Reason:            "Emergency",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindActiveLoanByUserID", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errMsg:  "cannot exceed maximum loan request",
		},
		{
			name: "error repo create fails",
			req: &LoanRequest{
				UserID:            1,
				EmployeeID:        1,
				TotalAmount:       5000000,
				InstallmentAmount: 500000,
				Reason:            "Emergency",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindActiveLoanByUserID", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*loan.Loan")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name: "error find approval users fails",
			req: &LoanRequest{
				UserID:            1,
				EmployeeID:        1,
				TotalAmount:       5000000,
				InstallmentAmount: 500000,
				Reason:            "Emergency",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindActiveLoanByUserID", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*loan.Loan")).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_LOAN)).Return([]uint(nil), errors.New("user service error"))
			},
			wantErr: true,
			errMsg:  "user service error",
		},
		{
			name: "error find active loan db error",
			req: &LoanRequest{
				UserID:            1,
				EmployeeID:        1,
				TotalAmount:       5000000,
				InstallmentAmount: 500000,
				Reason:            "Emergency",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindActiveLoanByUserID", mock.Anything, uint(1)).Return(nil, errors.New("db connection error"))
			},
			wantErr: true,
			errMsg:  "db connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, notif, userProv, _, _ := newTestLoanService()
			tt.setupMocks(repo, notif, userProv)

			err := svc.Create(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetLoanDetail(t *testing.T) {
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
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Loan{
					ID:                1,
					EmployeeID:        1,
					UserID:            1,
					TotalAmount:       5000000,
					InstallmentAmount: 500000,
					RemainingAmount:   5000000,
					Status:            constants.LoanStatusPending,
					Reason:            "Emergency",
					RejectionReason:   sql.NullString{Valid: false},
					CreatedAt:         time.Now(),
					User:              user.User{ID: 1},
					Employee:          user.Employee{ID: 1, FullName: "John Doe", NIK: "EMP001"},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "success with rejection reason",
			id:   2,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&Loan{
					ID:                2,
					EmployeeID:        1,
					UserID:            1,
					TotalAmount:       5000000,
					InstallmentAmount: 500000,
					RemainingAmount:   5000000,
					Status:            constants.LoanStatusRejected,
					Reason:            "Emergency",
					RejectionReason:   sql.NullString{String: "Not eligible", Valid: true},
					CreatedAt:         time.Now(),
					User:              user.User{ID: 1},
					Employee:          user.Employee{ID: 1, FullName: "John Doe", NIK: "EMP001"},
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
			name: "error data user not found",
			id:   3,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(3)).Return(&Loan{
					ID:       3,
					User:     user.User{ID: 0},
					Employee: user.Employee{ID: 0},
				}, nil)
			},
			wantErr: true,
			errMsg:  "data user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestLoanService()
			tt.setupMocks(repo)

			resp, err := svc.GetLoanDetail(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.id, resp.ID)
			}
		})
	}
}

func TestService_GetLoans(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     LoanFilter
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:   "success with data",
			filter: LoanFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]Loan{
						{
							ID:                1,
							EmployeeID:        1,
							TotalAmount:       5000000,
							InstallmentAmount: 500000,
							RemainingAmount:   5000000,
							Status:            constants.LoanStatusPending,
							Employee:          user.Employee{FullName: "John Doe", NIK: "EMP001"},
							CreatedAt:         time.Now(),
						},
					}, int64(1), nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "success empty list",
			filter: LoanFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]Loan{}, int64(0), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "repo returns error returns empty",
			filter: LoanFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]Loan(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestLoanService()
			tt.setupMocks(repo)

			list, meta, err := svc.GetLoans(ctx, tt.filter)

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
		setupMocks func(*mockRepo, *mockNotification)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "approve success",
			req: &ActionRequest{
				ID:           1,
				SuperAdminID: 10,
				Action:       string(constants.LoanActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Loan{
					ID:                1,
					UserID:            1,
					EmployeeID:        1,
					Status:            constants.LoanStatusPending,
					TotalAmount:       5000000,
					InstallmentAmount: 500000,
					RemainingAmount:   5000000,
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*loan.Loan")).Return(nil)
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "reject success",
			req: &ActionRequest{
				ID:              2,
				SuperAdminID:    10,
				Action:          string(constants.LoanActionReject),
				RejectionReason: "Not eligible",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&Loan{
					ID:                2,
					UserID:            1,
					EmployeeID:        1,
					Status:            constants.LoanStatusPending,
					TotalAmount:       5000000,
					InstallmentAmount: 500000,
					RemainingAmount:   5000000,
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*loan.Loan")).Return(nil)
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not pending status",
			req: &ActionRequest{
				ID:           3,
				SuperAdminID: 10,
				Action:       string(constants.LoanActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindByID", mock.Anything, uint(3)).Return(&Loan{
					ID:     3,
					Status: constants.LoanStatusApproved,
				}, nil)
			},
			wantErr: true,
			errMsg:  "cannot process loan with status APPROVED",
		},
		{
			name: "error rejection reason required",
			req: &ActionRequest{
				ID:              4,
				SuperAdminID:    10,
				Action:          string(constants.LoanActionReject),
				RejectionReason: "",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindByID", mock.Anything, uint(4)).Return(&Loan{
					ID:     4,
					UserID: 1,
					Status: constants.LoanStatusPending,
				}, nil)
			},
			wantErr: true,
			errMsg:  "rejection reason is required",
		},
		{
			name: "error loan not found",
			req: &ActionRequest{
				ID:           99,
				SuperAdminID: 10,
				Action:       string(constants.LoanActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
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
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindByID", mock.Anything, uint(5)).Return(&Loan{
					ID:     5,
					UserID: 1,
					Status: constants.LoanStatusPending,
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
				Action:       string(constants.LoanActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindByID", mock.Anything, uint(6)).Return(&Loan{
					ID:                6,
					UserID:            1,
					EmployeeID:        1,
					Status:            constants.LoanStatusPending,
					TotalAmount:       5000000,
					InstallmentAmount: 500000,
					RemainingAmount:   5000000,
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*loan.Loan")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, notif, _, _, _ := newTestLoanService()
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

func TestService_Export(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     LoanFilter
		setupMocks func(*mockRepo, *mockExcel)
		wantErr    bool
	}{
		{
			name:   "success",
			filter: LoanFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]Loan{
						{
							ID:                1,
							EmployeeID:        1,
							TotalAmount:       5000000,
							InstallmentAmount: 500000,
							RemainingAmount:   5000000,
							Status:            constants.LoanStatusPending,
							Employee:          user.Employee{FullName: "John Doe"},
							CreatedAt:         time.Now(),
						},
					}, int64(1), nil)
				excel.On("GenerateSimpleExcel", "Loans", mock.Anything, mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "success with empty employee name",
			filter: LoanFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]Loan{
						{
							ID:                2,
							EmployeeID:        1,
							TotalAmount:       3000000,
							InstallmentAmount: 300000,
							RemainingAmount:   3000000,
							Status:            constants.LoanStatusApproved,
							Employee:          user.Employee{FullName: ""},
							CreatedAt:         time.Now(),
						},
					}, int64(1), nil)
				excel.On("GenerateSimpleExcel", "Loans", mock.Anything, mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "error repo fails",
			filter: LoanFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]Loan(nil), int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "error excel generation fails",
			filter: LoanFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]Loan{
						{
							ID:                1,
							TotalAmount:       5000000,
							InstallmentAmount: 500000,
							RemainingAmount:   5000000,
							Status:            constants.LoanStatusPending,
							Employee:          user.Employee{FullName: "John Doe"},
							CreatedAt:         time.Now(),
						},
					}, int64(1), nil)
				excel.On("GenerateSimpleExcel", "Loans", mock.Anything, mock.Anything).Return(nil, errors.New("excel error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, excel := newTestLoanService()
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
