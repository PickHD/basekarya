package payroll

import (
	"errors"
	"testing"
	"time"

	"basekarya-backend/internal/modules/loan"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GenerateAll(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *GenerateRequest
		setupMocks func(*mockRepo, *mockUserProvider, *mockReimbursementProvider, *mockAttendanceProvider, *mockLoanProvider, *mockOvertimeProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success with one employee",
			req:  &GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(repo *mockRepo, userP *mockUserProvider, reimburse *mockReimbursementProvider, attend *mockAttendanceProvider, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider) {
				userP.On("FindAllEmployeeActive", mock.Anything).Return([]user.Employee{
					{ID: 1, UserID: 10, BaseSalary: 5000000},
				}, nil)
				repo.On("GetExistingEmployeeID", mock.Anything, 6, 2025).Return(map[uint]bool{}, nil)
				attend.On("GetBulkLateDuration", mock.Anything, 6, 2025).Return(map[uint]int{1: 30}, nil)
				reimburse.On("GetBulkApprovedAmount", mock.Anything, 6, 2025).Return(map[uint]float64{10: 200000}, nil)
				loanP.On("GetBulkActiveLoansByEmployeeIds", mock.Anything, mock.Anything).Return(map[uint]loan.Loan{}, nil)
				overtimeP.On("GetBulkActiveOvertimesByEmployeeIds", mock.Anything, 6, 2025, mock.Anything).Return(map[uint]int{}, nil)
				repo.On("CreateBulk", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success with overtime",
			req:  &GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(repo *mockRepo, userP *mockUserProvider, reimburse *mockReimbursementProvider, attend *mockAttendanceProvider, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider) {
				userP.On("FindAllEmployeeActive", mock.Anything).Return([]user.Employee{
					{ID: 2, UserID: 20, BaseSalary: 5000000},
				}, nil)
				repo.On("GetExistingEmployeeID", mock.Anything, 6, 2025).Return(map[uint]bool{}, nil)
				attend.On("GetBulkLateDuration", mock.Anything, 6, 2025).Return(map[uint]int{}, nil)
				reimburse.On("GetBulkApprovedAmount", mock.Anything, 6, 2025).Return(map[uint]float64{}, nil)
				loanP.On("GetBulkActiveLoansByEmployeeIds", mock.Anything, mock.Anything).Return(map[uint]loan.Loan{}, nil)
				overtimeP.On("GetBulkActiveOvertimesByEmployeeIds", mock.Anything, 6, 2025, mock.Anything).Return(map[uint]int{2: 120}, nil)
				repo.On("CreateBulk", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "skip existing employee",
			req:  &GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(repo *mockRepo, userP *mockUserProvider, reimburse *mockReimbursementProvider, attend *mockAttendanceProvider, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider) {
				userP.On("FindAllEmployeeActive", mock.Anything).Return([]user.Employee{
					{ID: 1, UserID: 10, BaseSalary: 5000000},
				}, nil)
				repo.On("GetExistingEmployeeID", mock.Anything, 6, 2025).Return(map[uint]bool{1: true}, nil)
				attend.On("GetBulkLateDuration", mock.Anything, 6, 2025).Return(map[uint]int{}, nil)
				reimburse.On("GetBulkApprovedAmount", mock.Anything, 6, 2025).Return(map[uint]float64{}, nil)
				loanP.On("GetBulkActiveLoansByEmployeeIds", mock.Anything, mock.Anything).Return(map[uint]loan.Loan{}, nil)
				overtimeP.On("GetBulkActiveOvertimesByEmployeeIds", mock.Anything, 6, 2025, mock.Anything).Return(map[uint]int{}, nil)
			},
			wantErr: false,
		},
		{
			name: "error fetch employees",
			req:  &GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(repo *mockRepo, userP *mockUserProvider, reimburse *mockReimbursementProvider, attend *mockAttendanceProvider, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider) {
				userP.On("FindAllEmployeeActive", mock.Anything).Return([]user.Employee(nil), errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "failed to fetch all employee active: db error",
		},
		{
			name: "error fetch existing",
			req:  &GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(repo *mockRepo, userP *mockUserProvider, reimburse *mockReimbursementProvider, attend *mockAttendanceProvider, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider) {
				userP.On("FindAllEmployeeActive", mock.Anything).Return([]user.Employee{{ID: 1, BaseSalary: 5000000}}, nil)
				repo.On("GetExistingEmployeeID", mock.Anything, 6, 2025).Return(map[uint]bool(nil), errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "failed to fetch existing employee id: db error",
		},
		{
			name: "error fetch attendance",
			req:  &GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(repo *mockRepo, userP *mockUserProvider, reimburse *mockReimbursementProvider, attend *mockAttendanceProvider, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider) {
				userP.On("FindAllEmployeeActive", mock.Anything).Return([]user.Employee{{ID: 1, BaseSalary: 5000000}}, nil)
				repo.On("GetExistingEmployeeID", mock.Anything, 6, 2025).Return(map[uint]bool{}, nil)
				attend.On("GetBulkLateDuration", mock.Anything, 6, 2025).Return(map[uint]int(nil), errors.New("attend error"))
			},
			wantErr: true,
			errMsg:  "failed to fetch bulk late duration: attend error",
		},
		{
			name: "error fetch reimbursement",
			req:  &GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(repo *mockRepo, userP *mockUserProvider, reimburse *mockReimbursementProvider, attend *mockAttendanceProvider, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider) {
				userP.On("FindAllEmployeeActive", mock.Anything).Return([]user.Employee{{ID: 1, BaseSalary: 5000000}}, nil)
				repo.On("GetExistingEmployeeID", mock.Anything, 6, 2025).Return(map[uint]bool{}, nil)
				attend.On("GetBulkLateDuration", mock.Anything, 6, 2025).Return(map[uint]int{}, nil)
				reimburse.On("GetBulkApprovedAmount", mock.Anything, 6, 2025).Return(map[uint]float64(nil), errors.New("reimburse error"))
			},
			wantErr: true,
			errMsg:  "failed to fetch bulk approved amount: reimburse error",
		},
		{
			name: "error create bulk",
			req:  &GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(repo *mockRepo, userP *mockUserProvider, reimburse *mockReimbursementProvider, attend *mockAttendanceProvider, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider) {
				userP.On("FindAllEmployeeActive", mock.Anything).Return([]user.Employee{{ID: 1, UserID: 10, BaseSalary: 5000000}}, nil)
				repo.On("GetExistingEmployeeID", mock.Anything, 6, 2025).Return(map[uint]bool{}, nil)
				attend.On("GetBulkLateDuration", mock.Anything, 6, 2025).Return(map[uint]int{}, nil)
				reimburse.On("GetBulkApprovedAmount", mock.Anything, 6, 2025).Return(map[uint]float64{}, nil)
				loanP.On("GetBulkActiveLoansByEmployeeIds", mock.Anything, mock.Anything).Return(map[uint]loan.Loan{}, nil)
				overtimeP.On("GetBulkActiveOvertimesByEmployeeIds", mock.Anything, 6, 2025, mock.Anything).Return(map[uint]int{}, nil)
				repo.On("CreateBulk", mock.Anything, mock.Anything).Return(errors.New("insert error"))
			},
			wantErr: true,
			errMsg:  "insert error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, userP, reimburse, attend, _, _, _, _, loanP, overtimeP := newTestService()
			tt.setupMocks(repo, userP, reimburse, attend, loanP, overtimeP)

			resp, err := svc.GenerateAll(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				if resp != nil {
					assert.Equal(t, 2025, resp.Year)
					assert.Equal(t, 6, resp.Month)
				}
			}
		})
	}
}

func TestService_GetList(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     *PayrollFilter
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:   "success with data",
			filter: &PayrollFilter{Page: 1, Limit: 10, Month: 6, Year: 2025},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("*payroll.PayrollFilter")).Return([]Payroll{
					{
						ID: 1, PeriodDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
						NetSalary: 5000000, Status: constants.PayrollStatusDraft,
						Employee: &user.Employee{FullName: "John Doe", NIK: "EMP001"},
					},
				}, int64(1), nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "success empty",
			filter: &PayrollFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("*payroll.PayrollFilter")).Return([]Payroll{}, int64(0), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "error repo fails",
			filter: &PayrollFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("*payroll.PayrollFilter")).Return([]Payroll(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			list, meta, err := svc.GetList(ctx, tt.filter)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, list, tt.wantLen)
				if tt.wantLen > 0 {
					assert.NotNil(t, meta)
				}
			}
		})
	}
}

func TestService_GetDetail(t *testing.T) {
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
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Payroll{
					ID: 1, BaseSalary: 5000000, TotalAllowance: 5500000, TotalDeduction: 500000, NetSalary: 5000000,
					PeriodDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
					Status:     constants.PayrollStatusDraft,
					Employee: &user.Employee{
						ID: 1, FullName: "John Doe", NIK: "EMP001",
						BankName: "BCA", BankAccountNumber: "1234567890", BankAccountHolder: "John Doe",
					},
					Details: []PayrollDetail{
						{ID: 1, PayrollID: 1, Title: "Base Salary", Type: constants.DetailTypeAllowance, Amount: 5000000},
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   999,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name: "error employee nil",
			id:   2,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&Payroll{
					ID:       2,
					Employee: nil,
				}, nil)
			},
			wantErr: true,
			errMsg:  "employee not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			resp, err := svc.GetDetail(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.id, resp.ID)
				assert.Equal(t, "John Doe", resp.EmployeeName)
			}
		})
	}
}

func TestService_MarkAsPaid(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		id         uint
		setupMocks func(*mockRepo, *mockLoanProvider, *mockOvertimeProvider, *mockNotificationProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success draft to paid",
			id:   1,
			setupMocks: func(repo *mockRepo, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider, notif *mockNotificationProvider) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Payroll{
					ID: 1, EmployeeID: 1,
					PeriodDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
					Status:     constants.PayrollStatusDraft,
					Employee:   &user.Employee{UserID: 10},
					Details:    []PayrollDetail{},
				}, nil)
				repo.On("UpdateStatus", mock.Anything, uint(1), constants.PayrollStatusPaid).Return(nil)
				overtimeP.On("UpdateBulkStatusByEmployeeId", mock.Anything, uint(1), 6, 2025, constants.OvertimeStatusPaid).Return(nil)
				notif.On("SendNotification", mock.Anything, uint(10), mock.Anything, mock.Anything, mock.Anything, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "already paid",
			id:   2,
			setupMocks: func(repo *mockRepo, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider, notif *mockNotificationProvider) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&Payroll{
					ID:       2,
					Status:   constants.PayrollStatusPaid,
					Employee: &user.Employee{UserID: 10},
					Details:  []PayrollDetail{},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "success with loan deduction",
			id:   3,
			setupMocks: func(repo *mockRepo, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider, notif *mockNotificationProvider) {
				repo.On("FindByID", mock.Anything, uint(3)).Return(&Payroll{
					ID: 3, EmployeeID: 5,
					PeriodDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
					Status:     constants.PayrollStatusDraft,
					Employee:   &user.Employee{UserID: 10},
					Details: []PayrollDetail{
						{Title: "Potongan Kasbon", Type: constants.DetailTypeDeduction, Amount: 500000},
					},
				}, nil)
				repo.On("UpdateStatus", mock.Anything, uint(3), constants.PayrollStatusPaid).Return(nil)
				loanP.On("GetBulkActiveLoansByEmployeeIds", mock.Anything, []uint{5}).Return(map[uint]loan.Loan{
					5: {ID: 1, InstallmentAmount: 500000, RemainingAmount: 1000000},
				}, nil)
				loanP.On("Update", mock.Anything, mock.AnythingOfType("*loan.Loan")).Return(nil)
				overtimeP.On("UpdateBulkStatusByEmployeeId", mock.Anything, uint(5), 6, 2025, constants.OvertimeStatusPaid).Return(nil)
				notif.On("SendNotification", mock.Anything, uint(10), mock.Anything, mock.Anything, mock.Anything, uint(3)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error find by id",
			id:   999,
			setupMocks: func(repo *mockRepo, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider, notif *mockNotificationProvider) {
				repo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name: "error update status",
			id:   1,
			setupMocks: func(repo *mockRepo, loanP *mockLoanProvider, overtimeP *mockOvertimeProvider, notif *mockNotificationProvider) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Payroll{
					ID: 1, EmployeeID: 1,
					PeriodDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
					Status:     constants.PayrollStatusDraft,
					Employee:   &user.Employee{UserID: 10},
					Details:    []PayrollDetail{},
				}, nil)
				repo.On("UpdateStatus", mock.Anything, uint(1), constants.PayrollStatusPaid).Return(errors.New("update error"))
			},
			wantErr: true,
			errMsg:  "update error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, notif, _, _, loanP, overtimeP := newTestService()
			tt.setupMocks(repo, loanP, overtimeP, notif)

			err := svc.MarkAsPaid(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			time.Sleep(50 * time.Millisecond)
		})
	}
}
