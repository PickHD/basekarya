package contract

import (
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

func newTestContractService() (Service, *mockRepo, *mockStorageProvider, *mockNotificationProvider, *mockUserProvider, *mockExcel) {
	repo := new(mockRepo)
	storage := new(mockStorageProvider)
	notif := new(mockNotificationProvider)
	userProv := new(mockUserProvider)
	excel := new(mockExcel)

	svc := NewService(repo, storage, notif, userProv, excel)
	return svc, repo, storage, notif, userProv, excel
}

func TestService_Upsert(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *UpsertContractRequest
		setupMocks func(*mockRepo, *mockStorageProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success PKWT",
			req: &UpsertContractRequest{
				EmployeeID:     1,
				ContractType:   constants.ContractTypePKWT,
				ContractNumber: "CTR-001",
				StartDate:      "2026-01-01",
				EndDate:        "2026-12-31",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorageProvider) {
				repo.On("FindByEmployeeID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
				repo.On("Upsert", mock.Anything, mock.AnythingOfType("*contract.Contract")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success PKWTT without end date",
			req: &UpsertContractRequest{
				EmployeeID:     1,
				ContractType:   constants.ContractTypePKWTT,
				ContractNumber: "CTR-002",
				StartDate:      "2026-01-01",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorageProvider) {
				repo.On("FindByEmployeeID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
				repo.On("Upsert", mock.Anything, mock.AnythingOfType("*contract.Contract")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error invalid start date",
			req: &UpsertContractRequest{
				EmployeeID:   1,
				ContractType: constants.ContractTypePKWT,
				StartDate:    "not-a-date",
				EndDate:      "2026-12-31",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorageProvider) {},
			wantErr:    true,
			errMsg:     "invalid start date format",
		},
		{
			name: "error PKWT without end date",
			req: &UpsertContractRequest{
				EmployeeID:   1,
				ContractType: constants.ContractTypePKWT,
				StartDate:    "2026-01-01",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorageProvider) {},
			wantErr:    true,
			errMsg:     "end date is required for PKWT",
		},
		{
			name: "error invalid end date format",
			req: &UpsertContractRequest{
				EmployeeID:   1,
				ContractType: constants.ContractTypePKWT,
				StartDate:    "2026-01-01",
				EndDate:      "not-a-date",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorageProvider) {},
			wantErr:    true,
			errMsg:     "invalid end date format",
		},
		{
			name: "error end date before start date",
			req: &UpsertContractRequest{
				EmployeeID:   1,
				ContractType: constants.ContractTypePKWT,
				StartDate:    "2026-12-31",
				EndDate:      "2026-01-01",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorageProvider) {},
			wantErr:    true,
			errMsg:     "end date must be after start date",
		},
		{
			name: "error repo upsert fails",
			req: &UpsertContractRequest{
				EmployeeID:     1,
				ContractType:   constants.ContractTypePKWT,
				ContractNumber: "CTR-003",
				StartDate:      "2026-01-01",
				EndDate:        "2026-12-31",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorageProvider) {
				repo.On("FindByEmployeeID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
				repo.On("Upsert", mock.Anything, mock.AnythingOfType("*contract.Contract")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, storage, _, _, _ := newTestContractService()
			tt.setupMocks(repo, storage)

			err := svc.Upsert(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
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
		filter     *ContractFilter
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:   "success with data",
			filter: &ContractFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("*contract.ContractFilter")).Return([]Contract{
					{
						ID: 1, EmployeeID: 1, ContractType: constants.ContractTypePKWT,
						Employee: &user.Employee{FullName: "John", NIK: "001"},
					},
				}, int64(1), nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "success empty",
			filter: &ContractFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("*contract.ContractFilter")).Return([]Contract{}, int64(0), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "error repo fails",
			filter: &ContractFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("*contract.ContractFilter")).Return([]Contract(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestContractService()
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

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

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
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Contract{
					ID: 1, EmployeeID: 1, ContractType: constants.ContractTypePKWT,
					StartDate: startDate,
					Employee:  &user.Employee{FullName: "John", NIK: "001"},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestContractService()
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

func TestService_GetByEmployeeID(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		employeeID uint
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name:       "success",
			employeeID: 1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByEmployeeID", mock.Anything, uint(1)).Return(&Contract{
					ID: 1, EmployeeID: 1, ContractType: constants.ContractTypePKWT,
					Employee: &user.Employee{FullName: "John", NIK: "001"},
				}, nil)
			},
			wantErr: false,
		},
		{
			name:       "not found",
			employeeID: 99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByEmployeeID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestContractService()
			tt.setupMocks(repo)

			resp, err := svc.GetByEmployeeID(ctx, tt.employeeID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
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
				repo.On("SoftDelete", mock.Anything, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("SoftDelete", mock.Anything, uint(99)).Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestContractService()
			tt.setupMocks(repo)

			err := svc.Delete(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
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
		filter     *ContractFilter
		setupMocks func(*mockRepo, *mockExcel)
		wantErr    bool
	}{
		{
			name:   "success",
			filter: &ContractFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("*contract.ContractFilter")).Return([]Contract{
					{
						ID: 1, EmployeeID: 1, ContractType: constants.ContractTypePKWT,
						StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   &endDate,
						Employee:  &user.Employee{FullName: "John", NIK: "001"},
					},
				}, int64(1), nil)
				excel.On("GenerateSimpleExcel", "Contracts", mock.Anything, mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "error repo fails",
			filter: &ContractFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAll", mock.Anything, mock.AnythingOfType("*contract.ContractFilter")).Return([]Contract(nil), int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, excel := newTestContractService()
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

func TestService_CheckExpiringContracts(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	endDate := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		setupMocks func(*mockRepo, *mockNotificationProvider, *mockUserProvider)
		wantErr    bool
	}{
		{
			name: "success with expiring contracts",
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider, userProv *mockUserProvider) {
				repo.On("FindExpiringContracts", mock.Anything, 30).Return([]Contract{
					{
						ID: 1, ContractNumber: "CTR-001", EndDate: &endDate,
						Employee: &user.Employee{FullName: "John"},
					},
				}, nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.VIEW_CONTRACT)).Return([]uint{10, 11}, nil)
				notif.On("BlastNotification", mock.Anything, []uint{10, 11}, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("MarkAlerted", mock.Anything, []uint{1}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "no expiring contracts",
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider, userProv *mockUserProvider) {
				repo.On("FindExpiringContracts", mock.Anything, 30).Return([]Contract{}, nil)
			},
			wantErr: false,
		},
		{
			name: "no approval users",
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider, userProv *mockUserProvider) {
				repo.On("FindExpiringContracts", mock.Anything, 30).Return([]Contract{
					{ID: 1, ContractNumber: "CTR-001", EndDate: &endDate, Employee: &user.Employee{FullName: "John"}},
				}, nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.VIEW_CONTRACT)).Return([]uint{}, nil)
			},
			wantErr: false,
		},
		{
			name: "error find expiring contracts",
			setupMocks: func(repo *mockRepo, notif *mockNotificationProvider, userProv *mockUserProvider) {
				repo.On("FindExpiringContracts", mock.Anything, 30).Return([]Contract(nil), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, notif, userProv, _ := newTestContractService()
			tt.setupMocks(repo, notif, userProv)

			err := svc.CheckExpiringContracts(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
