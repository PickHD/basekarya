package leave

import (
	"errors"
	"testing"
	"time"

	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestLeaveService() (Service, *mockRepo, *mockStorage, *mockNotification, *mockUserProvider, *testutil.MockTransactionManager, *mockExcel) {
	repo := new(mockRepo)
	storage := new(mockStorage)
	notif := new(mockNotification)
	userProv := new(mockUserProvider)
	tm := testutil.NewMockTransactionManager()
	excel := new(mockExcel)

	svc := NewService(repo, storage, notif, userProv, tm, excel)
	return svc, repo, storage, notif, userProv, tm, excel
}

func TestService_Apply(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *ApplyRequest
		setupMocks func(*mockRepo, *mockStorage, *mockNotification, *mockUserProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success without attachment",
			req: &ApplyRequest{
				EmployeeID:  1,
				LeaveTypeID: 1,
				StartDate:   "2026-06-01",
				EndDate:     "2026-06-02",
				Reason:      "Family event",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("GetBalance", mock.Anything, uint(1), uint(1), 2026).Return(&LeaveBalance{QuotaLeft: 5}, nil)
				repo.On("CreateRequest", mock.Anything, mock.AnythingOfType("*leave.LeaveRequest")).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_LEAVE)).Return([]uint{10, 11}, nil)
				notif.On("BlastNotification", mock.Anything, []uint{10, 11}, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error invalid start date format",
			req: &ApplyRequest{
				EmployeeID:  1,
				LeaveTypeID: 1,
				StartDate:   "not-a-date",
				EndDate:     "2026-06-02",
				Reason:      "Test",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {},
			wantErr: true,
			errMsg: "invalid start date format",
		},
		{
			name: "error invalid end date format",
			req: &ApplyRequest{
				EmployeeID:  1,
				LeaveTypeID: 1,
				StartDate:   "2026-06-01",
				EndDate:     "not-a-date",
				Reason:      "Test",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {},
			wantErr: true,
			errMsg: "invalid end date format",
		},
		{
			name: "error end date before start date",
			req: &ApplyRequest{
				EmployeeID:  1,
				LeaveTypeID: 1,
				StartDate:   "2026-06-05",
				EndDate:     "2026-06-01",
				Reason:      "Test",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {},
			wantErr: true,
			errMsg: "end date must be after start date",
		},
		{
			name: "error insufficient leave balance",
			req: &ApplyRequest{
				EmployeeID:  1,
				LeaveTypeID: 1,
				StartDate:   "2026-06-01",
				EndDate:     "2026-06-05",
				Reason:      "Vacation",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("GetBalance", mock.Anything, uint(1), uint(1), 2026).Return(&LeaveBalance{QuotaLeft: 2}, nil)
			},
			wantErr: true,
			errMsg: "insufficient leave balance",
		},
		{
			name: "error balance lookup fails",
			req: &ApplyRequest{
				EmployeeID:  1,
				LeaveTypeID: 1,
				StartDate:   "2026-06-01",
				EndDate:     "2026-06-02",
				Reason:      "Test",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("GetBalance", mock.Anything, uint(1), uint(1), 2026).Return(nil, errors.New("db error"))
			},
			wantErr: true,
			errMsg: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, storage, notif, userProv, _, _ := newTestLeaveService()
			tt.setupMocks(repo, storage, notif, userProv)

			err := svc.Apply(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_RequestAction(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *LeaveActionRequest
		setupMocks func(*mockRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "approve success",
			req: &LeaveActionRequest{
				RequestID:  1,
				ApproverID: 10,
				Action:     string(constants.LeaveActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&LeaveRequest{
					ID:          1,
					EmployeeID:  1,
					LeaveTypeID: 1,
					TotalDays:   2,
					Status:      constants.LeaveStatusPending,
					StartDate:   time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC),
					LeaveType:   &master.LeaveType{ID: 1, IsDeducted: false},
					Employee:    &user.Employee{ID: 1, ShiftID: 1},
					User:        user.User{ID: 1},
				}, nil)
				repo.On("ApproveRequest", mock.Anything, uint(1), uint(10), mock.Anything, false, 2).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "approve with deduction success",
			req: &LeaveActionRequest{
				RequestID:  2,
				ApproverID: 10,
				Action:     string(constants.LeaveActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequestByID", mock.Anything, uint(2)).Return(&LeaveRequest{
					ID:          2,
					EmployeeID:  1,
					LeaveTypeID: 1,
					TotalDays:   2,
					Status:      constants.LeaveStatusPending,
					StartDate:   time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC),
					LeaveType:   &master.LeaveType{ID: 1, IsDeducted: true},
					Employee:    &user.Employee{ID: 1, ShiftID: 1},
					User:        user.User{ID: 1},
				}, nil)
				repo.On("GetBalance", mock.Anything, uint(1), uint(1), 2026).Return(&LeaveBalance{QuotaLeft: 5}, nil)
				repo.On("ApproveRequest", mock.Anything, uint(2), uint(10), mock.Anything, true, 2).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "reject success",
			req: &LeaveActionRequest{
				RequestID:       3,
				ApproverID:      10,
				Action:          string(constants.LeaveActionReject),
				RejectionReason: "Not eligible",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequestByID", mock.Anything, uint(3)).Return(&LeaveRequest{
					ID:     3,
					Status: constants.LeaveStatusPending,
					User:   user.User{ID: 1},
				}, nil)
				repo.On("RejectRequest", mock.Anything, uint(3), uint(10), "Not eligible").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error request not pending",
			req: &LeaveActionRequest{
				RequestID:  4,
				ApproverID: 10,
				Action:     string(constants.LeaveActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequestByID", mock.Anything, uint(4)).Return(&LeaveRequest{
					ID:     4,
					Status: constants.LeaveStatusApproved,
				}, nil)
			},
			wantErr: true,
			errMsg: "request is not pending",
		},
		{
			name: "error rejection reason required",
			req: &LeaveActionRequest{
				RequestID:       5,
				ApproverID:      10,
				Action:          string(constants.LeaveActionReject),
				RejectionReason: "",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequestByID", mock.Anything, uint(5)).Return(&LeaveRequest{
					ID:     5,
					Status: constants.LeaveStatusPending,
					User:   user.User{ID: 1},
				}, nil)
			},
			wantErr: true,
			errMsg: "rejection reason required",
		},
		{
			name: "error request not found",
			req: &LeaveActionRequest{
				RequestID:  99,
				ApproverID: 10,
				Action:     string(constants.LeaveActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequestByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg: "not found",
		},
		{
			name: "error invalid action",
			req: &LeaveActionRequest{
				RequestID:  6,
				ApproverID: 10,
				Action:     "INVALID",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequestByID", mock.Anything, uint(6)).Return(&LeaveRequest{
					ID:     6,
					Status: constants.LeaveStatusPending,
					User:   user.User{ID: 1},
				}, nil)
			},
			wantErr: true,
			errMsg: "invalid action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, notif, _, _, _ := newTestLeaveService()
			tt.setupMocks(repo)

			// Only set notification mock for success cases
			if !tt.wantErr && (tt.req.Action == string(constants.LeaveActionApprove) || tt.req.Action == string(constants.LeaveActionReject)) {
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}

			err := svc.RequestAction(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
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
				repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&LeaveRequest{
					ID:         1,
					EmployeeID: 1,
					StartDate:  time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
					EndDate:    time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC),
					Status:     constants.LeaveStatusPending,
					Reason:     "Family event",
					Employee:   &user.Employee{FullName: "John", NIK: "001"},
					LeaveType:  &master.LeaveType{ID: 1, Name: "Annual", DefaultQuota: 12, IsDeducted: true},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequestByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _ := newTestLeaveService()
			tt.setupMocks(repo)

			resp, err := svc.GetDetail(ctx, tt.id)

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

func TestService_GetList(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     *LeaveFilter
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name: "success with data",
			filter: &LeaveFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllRequests", mock.Anything, mock.AnythingOfType("*leave.LeaveFilter")).
					Return([]LeaveRequest{
						{
							ID:         1,
							EmployeeID: 1,
							Employee:   &user.Employee{FullName: "John", NIK: "001"},
							LeaveType:  &master.LeaveType{ID: 1, Name: "Annual"},
							Status:     constants.LeaveStatusPending,
						},
					}, int64(1), nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "success empty list",
			filter: &LeaveFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllRequests", mock.Anything, mock.AnythingOfType("*leave.LeaveFilter")).
					Return([]LeaveRequest{}, int64(0), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "repo returns error returns empty",
			filter: &LeaveFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllRequests", mock.Anything, mock.AnythingOfType("*leave.LeaveFilter")).
					Return([]LeaveRequest(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: false, // service returns empty, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _ := newTestLeaveService()
			tt.setupMocks(repo)

			list, meta, err := svc.GetList(ctx, tt.filter)

			require.NoError(t, err)
			assert.Len(t, list, tt.wantLen)
			if tt.wantLen > 0 {
				assert.NotNil(t, meta)
			}
		})
	}
}

func TestService_GenerateInitialBalance(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		employeeID uint
		setupMocks func(*mockRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "success with leave types",
			employeeID: 1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllLeaveTypes", mock.Anything).Return([]master.LeaveType{
					{ID: 1, Name: "Annual", DefaultQuota: 12},
					{ID: 2, Name: "Sick", DefaultQuota: 6},
				}, nil)
				repo.On("CreateLeaveBalances", mock.Anything, mock.AnythingOfType("[]leave.LeaveBalance")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "error find leave types fails",
			employeeID: 1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllLeaveTypes", mock.Anything).Return([]master.LeaveType(nil), errors.New("db error"))
			},
			wantErr: true,
			errMsg: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _ := newTestLeaveService()
			tt.setupMocks(repo)

			err := svc.GenerateInitialBalance(ctx, tt.employeeID)

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
		filter     *LeaveFilter
		setupMocks func(*mockRepo, *mockExcel)
		wantErr    bool
	}{
		{
			name:   "success",
			filter: &LeaveFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAllRequests", mock.Anything, mock.AnythingOfType("*leave.LeaveFilter")).
					Return([]LeaveRequest{
						{
							ID:         1,
							EmployeeID: 1,
							Employee:   &user.Employee{FullName: "John", NIK: "001"},
							LeaveType:  &master.LeaveType{Name: "Annual"},
							StartDate:  time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
							EndDate:    time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC),
							Status:     constants.LeaveStatusPending,
							Reason:     "Test",
							CreatedAt:  time.Now(),
						},
					}, int64(1), nil)
				excel.On("GenerateSimpleExcel", "Leaves", mock.Anything, mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "error repo fails",
			filter: &LeaveFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAllRequests", mock.Anything, mock.AnythingOfType("*leave.LeaveFilter")).
					Return([]LeaveRequest(nil), int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, excel := newTestLeaveService()
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
