package overtime

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestOvertimeService() (Service, *mockRepo, *mockNotification, *mockUserProvider, *testutil.MockTransactionManager, *mockExcel) {
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
		req        *OvertimeRequest
		setupMocks func(*mockRepo, *mockNotification, *mockUserProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req: &OvertimeRequest{
				UserID:     1,
				EmployeeID: 1,
				Date:       "2026-06-01",
				StartTime:  "18:00",
				EndTime:    "20:00",
				Reason:     "Project deadline",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*overtime.Overtime")).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_OVERTIME)).Return([]uint{10}, nil)
				notif.On("BlastNotification", mock.Anything, []uint{10}, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error zero user and employee",
			req: &OvertimeRequest{
				UserID:     0,
				EmployeeID: 0,
				Date:       "2026-06-01",
				StartTime:  "18:00",
				EndTime:    "20:00",
				Reason:     "Test",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {},
			wantErr:    true,
			errMsg:     "user not found",
		},
		{
			name: "error invalid start time",
			req: &OvertimeRequest{
				UserID:     1,
				EmployeeID: 1,
				Date:       "2026-06-01",
				StartTime:  "invalid",
				EndTime:    "20:00",
				Reason:     "Test",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {},
			wantErr:    true,
			errMsg:     "invalid start time format",
		},
		{
			name: "error invalid end time",
			req: &OvertimeRequest{
				UserID:     1,
				EmployeeID: 1,
				Date:       "2026-06-01",
				StartTime:  "18:00",
				EndTime:    "invalid",
				Reason:     "Test",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {},
			wantErr:    true,
			errMsg:     "invalid end time format",
		},
		{
			name: "error repo create fails",
			req: &OvertimeRequest{
				UserID:     1,
				EmployeeID: 1,
				Date:       "2026-06-01",
				StartTime:  "18:00",
				EndTime:    "20:00",
				Reason:     "Test",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*overtime.Overtime")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, notif, userProv, _, _ := newTestOvertimeService()
			tt.setupMocks(repo, notif, userProv)

			err := svc.Create(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
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
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Overtime{
					ID:              1,
					EmployeeID:      1,
					Employee:        user.Employee{FullName: "John", NIK: "001"},
					User:            user.User{ID: 1},
					Date:            "2026-06-01",
					StartTime:       "18:00",
					EndTime:         "20:00",
					DurationMinutes: 120,
					Reason:          "Project",
					Status:          constants.OvertimeStatusPending,
					RejectionReason: sql.NullString{},
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
			name: "error user data not found",
			id:   2,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&Overtime{
					ID:         2,
					EmployeeID: 1,
					Employee:   user.Employee{},
					User:       user.User{},
				}, nil)
			},
			wantErr: true,
			errMsg:  "data user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestOvertimeService()
			tt.setupMocks(repo)

			resp, err := svc.GetDetail(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, resp.ID)
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
				Action:       string(constants.OvertimeActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Overtime{
					ID:     1,
					UserID: 1,
					Status: constants.OvertimeStatusPending,
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*overtime.Overtime")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "reject success",
			req: &ActionRequest{
				ID:              2,
				SuperAdminID:    10,
				Action:          string(constants.OvertimeActionReject),
				RejectionReason: "Not eligible",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&Overtime{
					ID:     2,
					UserID: 1,
					Status: constants.OvertimeStatusPending,
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*overtime.Overtime")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not pending",
			req: &ActionRequest{
				ID:           3,
				SuperAdminID: 10,
				Action:       string(constants.OvertimeActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(3)).Return(&Overtime{
					ID:     3,
					Status: constants.OvertimeStatusApproved,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "error rejection reason required",
			req: &ActionRequest{
				ID:           4,
				SuperAdminID: 10,
				Action:       string(constants.OvertimeActionReject),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(4)).Return(&Overtime{
					ID:     4,
					UserID: 1,
					Status: constants.OvertimeStatusPending,
				}, nil)
			},
			wantErr: true,
			errMsg:  "rejection reason is required",
		},
		{
			name: "error invalid action",
			req: &ActionRequest{
				ID:           5,
				SuperAdminID: 10,
				Action:       "INVALID",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(5)).Return(&Overtime{
					ID:     5,
					UserID: 1,
					Status: constants.OvertimeStatusPending,
				}, nil)
			},
			wantErr: true,
			errMsg:  "invalid action",
		},
		{
			name: "error find fails",
			req: &ActionRequest{
				ID:           6,
				SuperAdminID: 10,
				Action:       string(constants.OvertimeActionApprove),
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(6)).Return(nil, errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, notif, _, _, _ := newTestOvertimeService()
			tt.setupMocks(repo)

			if !tt.wantErr && (tt.req.Action == string(constants.OvertimeActionApprove) || tt.req.Action == string(constants.OvertimeActionReject)) {
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}

			err := svc.ProcessAction(ctx, tt.req)

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

func TestService_GetList(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     OvertimeFilter
		setupMocks func(*mockRepo)
		wantLen    int
	}{
		{
			name:   "success with data",
			filter: OvertimeFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("overtime.OvertimeFilter")).
					Return([]Overtime{
						{
							ID: 1, EmployeeID: 1,
							Employee: user.Employee{FullName: "John", NIK: "001"},
							Status: constants.OvertimeStatusPending,
						},
					}, int64(1), nil)
			},
			wantLen: 1,
		},
		{
			name:   "success empty",
			filter: OvertimeFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("overtime.OvertimeFilter")).
					Return([]Overtime{}, int64(0), nil)
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestOvertimeService()
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

func TestService_Export(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     OvertimeFilter
		setupMocks func(*mockRepo, *mockExcel)
		wantErr    bool
	}{
		{
			name:   "success",
			filter: OvertimeFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("overtime.OvertimeFilter")).
					Return([]Overtime{
						{
							ID: 1, EmployeeID: 1,
							Employee:        user.Employee{FullName: "John"},
							Date:            "2026-06-01",
							StartTime:       "18:00",
							EndTime:         "20:00",
							DurationMinutes: 120,
							Reason:          "Project",
							Status:          constants.OvertimeStatusApproved,
							CreatedAt:       time.Now(),
						},
					}, int64(1), nil)
				excel.On("GenerateSimpleExcel", "Overtimes", mock.Anything, mock.Anything).Return([]byte("fake"), nil)
			},
			wantErr: false,
		},
		{
			name:   "error repo fails",
			filter: OvertimeFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("overtime.OvertimeFilter")).
					Return([]Overtime(nil), int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, excel := newTestOvertimeService()
			tt.setupMocks(repo, excel)

			data, err := svc.Export(ctx, tt.filter)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, data)
			}
		})
	}
}
