package asset

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

func newTestAssetService() (Service, *mockRepo, *mockNotification, *mockUserProvider, *testutil.MockTransactionManager, *mockExcel) {
	repo := new(mockRepo)
	notif := new(mockNotification)
	userProv := new(mockUserProvider)
	tm := testutil.NewMockTransactionManager()
	excel := new(mockExcel)

	svc := NewService(repo, notif, userProv, tm, excel)
	return svc, repo, notif, userProv, tm, excel
}

func TestService_CreateCategory(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CreateAssetCategoryRequest
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success",
			req: &CreateAssetCategoryRequest{
				Name:        "Laptop",
				Description: "Company laptops",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateCategory", mock.Anything, mock.AnythingOfType("*asset.AssetCategory")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error repo fails",
			req:  &CreateAssetCategoryRequest{Name: "Laptop"},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateCategory", mock.Anything, mock.AnythingOfType("*asset.AssetCategory")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestAssetService()
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
		filter     AssetCategoryFilter
		setupMocks func(*mockRepo)
		wantLen    int
	}{
		{
			name:   "success with data",
			filter: AssetCategoryFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllCategories", mock.Anything, mock.AnythingOfType("asset.AssetCategoryFilter")).
					Return([]AssetCategory{
						{ID: 1, Name: "Laptop", Description: "Company laptops"},
					}, int64(1), nil)
			},
			wantLen: 1,
		},
		{
			name:   "success empty list",
			filter: AssetCategoryFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllCategories", mock.Anything, mock.AnythingOfType("asset.AssetCategoryFilter")).
					Return([]AssetCategory{}, int64(0), nil)
			},
			wantLen: 0,
		},
		{
			name:   "repo returns error returns empty",
			filter: AssetCategoryFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllCategories", mock.Anything, mock.AnythingOfType("asset.AssetCategoryFilter")).
					Return([]AssetCategory(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestAssetService()
			tt.setupMocks(repo)
			list, meta, err := svc.GetCategories(ctx, tt.filter)
			require.NoError(t, err)
			assert.Len(t, list, tt.wantLen)
			if tt.wantLen > 0 {
				assert.NotNil(t, meta)
			}
		})
	}
}

func TestService_CreateAsset(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CreateAssetRequest
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success",
			req: &CreateAssetRequest{
				AssetCategoryID: 1,
				Name:            "MacBook Pro",
				Description:     "14 inch",
				SerialNumber:    "SN001",
				Condition:       constants.AssetConditionGood,
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateAsset", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success with default condition",
			req: &CreateAssetRequest{
				AssetCategoryID: 1,
				Name:            "Monitor",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateAsset", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error repo fails",
			req: &CreateAssetRequest{
				AssetCategoryID: 1,
				Name:            "MacBook Pro",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateAsset", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestAssetService()
			tt.setupMocks(repo)
			err := svc.CreateAsset(ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetAssetDetail(t *testing.T) {
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
				repo.On("FindAssetByID", mock.Anything, uint(1)).Return(&Asset{
					ID:              1,
					Name:            "MacBook Pro",
					AssetCategoryID: 1,
					Status:          constants.AssetStatusAvailable,
					Condition:       constants.AssetConditionGood,
					AssetCategory:   AssetCategory{ID: 1, Name: "Laptop"},
				}, nil)
				repo.On("FindActiveAssignmentByAssetID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
			},
			wantErr: false,
		},
		{
			name: "success with active assignment",
			id:   2,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAssetByID", mock.Anything, uint(2)).Return(&Asset{
					ID:              2,
					Name:            "Monitor",
					AssetCategoryID: 1,
					Status:          constants.AssetStatusAssigned,
					Condition:       constants.AssetConditionGood,
					AssetCategory:   AssetCategory{ID: 1, Name: "Laptop"},
				}, nil)
				repo.On("FindActiveAssignmentByAssetID", mock.Anything, uint(2)).Return(&AssetAssignment{
					ID:         1,
					AssetID:    2,
					EmployeeID: 1,
					Employee:   user.Employee{ID: 1, FullName: "John Doe"},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAssetByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestAssetService()
			tt.setupMocks(repo)
			resp, err := svc.GetAssetDetail(ctx, tt.id)
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

func TestService_GetAssets(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     AssetFilter
		setupMocks func(*mockRepo)
		wantLen    int
	}{
		{
			name:   "success with data",
			filter: AssetFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllAssets", mock.Anything, mock.AnythingOfType("asset.AssetFilter")).
					Return([]Asset{
						{ID: 1, Name: "MacBook Pro", AssetCategory: AssetCategory{ID: 1, Name: "Laptop"}, Status: constants.AssetStatusAvailable, Condition: constants.AssetConditionGood},
					}, int64(1), nil)
			},
			wantLen: 1,
		},
		{
			name:   "success empty list",
			filter: AssetFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllAssets", mock.Anything, mock.AnythingOfType("asset.AssetFilter")).
					Return([]Asset{}, int64(0), nil)
			},
			wantLen: 0,
		},
		{
			name:   "repo returns error returns empty",
			filter: AssetFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllAssets", mock.Anything, mock.AnythingOfType("asset.AssetFilter")).
					Return([]Asset(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestAssetService()
			tt.setupMocks(repo)
			list, meta, err := svc.GetAssets(ctx, tt.filter)
			require.NoError(t, err)
			assert.Len(t, list, tt.wantLen)
			if tt.wantLen > 0 {
				assert.NotNil(t, meta)
			}
		})
	}
}

func TestService_UpdateAsset(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *UpdateAssetRequest
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success update condition",
			req:  &UpdateAssetRequest{ID: 1, Condition: constants.AssetConditionDamaged},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAssetByID", mock.Anything, uint(1)).Return(&Asset{ID: 1, Name: "MacBook Pro", Condition: constants.AssetConditionGood}, nil)
				repo.On("UpdateAsset", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			req:  &UpdateAssetRequest{ID: 99},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAssetByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestAssetService()
			tt.setupMocks(repo)
			err := svc.UpdateAsset(ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_CreateAssignment(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CreateAssetAssignmentRequest
		setupMocks func(*mockRepo, *mockNotification, *mockUserProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req: &CreateAssetAssignmentRequest{
				UserID:             1,
				EmployeeID:         1,
				AssetID:            1,
				Purpose:            "Need for presentation",
				ExpectedReturnDate: "2025-12-31",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindAssetByID", mock.Anything, uint(1)).Return(&Asset{ID: 1, Name: "MacBook Pro", Status: constants.AssetStatusAvailable}, nil)
				repo.On("CreateAssignment", mock.Anything, mock.AnythingOfType("*asset.AssetAssignment")).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_ASSET)).Return([]uint{10, 11}, nil)
				notif.On("BlastNotification", mock.Anything, []uint{10, 11}, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error user not found",
			req: &CreateAssetAssignmentRequest{
				UserID:     0,
				EmployeeID: 0,
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {},
			wantErr:    true,
			errMsg:     "user not found",
		},
		{
			name: "error asset not found",
			req: &CreateAssetAssignmentRequest{
				UserID:     1,
				EmployeeID: 1,
				AssetID:    99,
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindAssetByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "asset not found",
		},
		{
			name: "error asset not available",
			req: &CreateAssetAssignmentRequest{
				UserID:     1,
				EmployeeID: 1,
				AssetID:    1,
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindAssetByID", mock.Anything, uint(1)).Return(&Asset{ID: 1, Status: constants.AssetStatusAssigned}, nil)
			},
			wantErr: true,
			errMsg:  "asset is not available for assignment",
		},
		{
			name: "error create assignment fails",
			req: &CreateAssetAssignmentRequest{
				UserID:     1,
				EmployeeID: 1,
				AssetID:    1,
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindAssetByID", mock.Anything, uint(1)).Return(&Asset{ID: 1, Status: constants.AssetStatusAvailable}, nil)
				repo.On("CreateAssignment", mock.Anything, mock.AnythingOfType("*asset.AssetAssignment")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name: "error find approval users fails",
			req: &CreateAssetAssignmentRequest{
				UserID:     1,
				EmployeeID: 1,
				AssetID:    1,
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification, userProv *mockUserProvider) {
				repo.On("FindAssetByID", mock.Anything, uint(1)).Return(&Asset{ID: 1, Status: constants.AssetStatusAvailable}, nil)
				repo.On("CreateAssignment", mock.Anything, mock.AnythingOfType("*asset.AssetAssignment")).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_ASSET)).Return([]uint(nil), errors.New("user service error"))
			},
			wantErr: true,
			errMsg:  "user service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, notif, userProv, _, _ := newTestAssetService()
			tt.setupMocks(repo, notif, userProv)
			err := svc.CreateAssignment(ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetAssignmentDetail(t *testing.T) {
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
				repo.On("FindAssignmentByID", mock.Anything, uint(1)).Return(&AssetAssignment{
					ID:         1,
					AssetID:    1,
					EmployeeID: 1,
					UserID:     1,
					Status:     constants.AssetAssignmentStatusActive,
					Purpose:    "Need for presentation",
					Asset:      Asset{ID: 1, Name: "MacBook Pro"},
					User:       user.User{ID: 1},
					Employee:   user.Employee{ID: 1, FullName: "John Doe", NIK: "EMP001"},
					CreatedAt:  time.Now(),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "success with rejection reason",
			id:   2,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAssignmentByID", mock.Anything, uint(2)).Return(&AssetAssignment{
					ID:              2,
					AssetID:         1,
					EmployeeID:      1,
					UserID:          1,
					Status:          constants.AssetAssignmentStatusRejected,
					RejectionReason: sql.NullString{String: "Not available", Valid: true},
					Asset:           Asset{ID: 1, Name: "MacBook Pro"},
					User:            user.User{ID: 1},
					Employee:        user.Employee{ID: 1, FullName: "John Doe", NIK: "EMP001"},
					CreatedAt:       time.Now(),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error data user not found",
			id:   3,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAssignmentByID", mock.Anything, uint(3)).Return(&AssetAssignment{
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
			svc, repo, _, _, _, _ := newTestAssetService()
			tt.setupMocks(repo)
			resp, err := svc.GetAssignmentDetail(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}

func TestService_GetAssignments(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     AssetAssignmentFilter
		setupMocks func(*mockRepo)
		wantLen    int
	}{
		{
			name:   "success with data",
			filter: AssetAssignmentFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllAssignments", mock.Anything, mock.AnythingOfType("asset.AssetAssignmentFilter")).
					Return([]AssetAssignment{
						{ID: 1, AssetID: 1, Employee: user.Employee{FullName: "John Doe", NIK: "EMP001"}, Asset: Asset{Name: "MacBook Pro"}, Status: constants.AssetAssignmentStatusPending, CreatedAt: time.Now()},
					}, int64(1), nil)
			},
			wantLen: 1,
		},
		{
			name:   "success empty list",
			filter: AssetAssignmentFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllAssignments", mock.Anything, mock.AnythingOfType("asset.AssetAssignmentFilter")).
					Return([]AssetAssignment{}, int64(0), nil)
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestAssetService()
			tt.setupMocks(repo)
			list, meta, err := svc.GetAssignments(ctx, tt.filter)
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
				Action:       string(constants.AssetAssignmentActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				assignment := &AssetAssignment{
					ID:     1,
					UserID: 1,
					Status: constants.AssetAssignmentStatusPending,
				}
				repo.On("FindAssignmentByID", mock.Anything, uint(1)).Return(assignment, nil)
				repo.On("UpdateAssignment", mock.Anything, mock.AnythingOfType("*asset.AssetAssignment")).Return(nil)
				repo.On("FindAssetByID", mock.Anything, mock.Anything).Return(&Asset{ID: 1, Status: constants.AssetStatusAvailable}, nil)
				repo.On("UpdateAsset", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(nil)
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "reject success",
			req: &ActionRequest{
				ID:              2,
				SuperAdminID:    10,
				Action:          string(constants.AssetAssignmentActionReject),
				RejectionReason: "Not available",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindAssignmentByID", mock.Anything, uint(2)).Return(&AssetAssignment{
					ID:     2,
					UserID: 1,
					Status: constants.AssetAssignmentStatusPending,
				}, nil)
				repo.On("UpdateAssignment", mock.Anything, mock.AnythingOfType("*asset.AssetAssignment")).Return(nil)
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not pending status",
			req: &ActionRequest{
				ID:           3,
				SuperAdminID: 10,
				Action:       string(constants.AssetAssignmentActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindAssignmentByID", mock.Anything, uint(3)).Return(&AssetAssignment{
					ID:     3,
					Status: constants.AssetAssignmentStatusActive,
				}, nil)
			},
			wantErr: true,
			errMsg:  "cannot process assignment with status ACTIVE",
		},
		{
			name: "error rejection reason required",
			req: &ActionRequest{
				ID:              4,
				SuperAdminID:    10,
				Action:          string(constants.AssetAssignmentActionReject),
				RejectionReason: "",
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindAssignmentByID", mock.Anything, uint(4)).Return(&AssetAssignment{
					ID:     4,
					UserID: 1,
					Status: constants.AssetAssignmentStatusPending,
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
				Action:       string(constants.AssetAssignmentActionApprove),
			},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindAssignmentByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
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
				repo.On("FindAssignmentByID", mock.Anything, uint(5)).Return(&AssetAssignment{
					ID:     5,
					UserID: 1,
					Status: constants.AssetAssignmentStatusPending,
				}, nil)
			},
			wantErr: true,
			errMsg:  "invalid action: INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, notif, _, _, _ := newTestAssetService()
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

func TestService_ProcessReturn(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *ReturnRequest
		setupMocks func(*mockRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req:  &ReturnRequest{ID: 1, UserID: 1},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAssignmentByID", mock.Anything, uint(1)).Return(&AssetAssignment{
					ID:      1,
					AssetID: 1,
					Status:  constants.AssetAssignmentStatusActive,
				}, nil)
				repo.On("UpdateAssignment", mock.Anything, mock.AnythingOfType("*asset.AssetAssignment")).Return(nil)
				repo.On("FindAssetByID", mock.Anything, uint(1)).Return(&Asset{ID: 1, Status: constants.AssetStatusAssigned}, nil)
				repo.On("UpdateAsset", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not active status",
			req:  &ReturnRequest{ID: 2},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAssignmentByID", mock.Anything, uint(2)).Return(&AssetAssignment{
					ID:     2,
					Status: constants.AssetAssignmentStatusReturned,
				}, nil)
			},
			wantErr: true,
			errMsg:  "cannot return assignment with status RETURNED",
		},
		{
			name: "error not found",
			req:  &ReturnRequest{ID: 99},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAssignmentByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestAssetService()
			tt.setupMocks(repo)
			err := svc.ProcessReturn(ctx, tt.req)
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
		filter     AssetFilter
		setupMocks func(*mockRepo, *mockExcel)
		wantErr    bool
	}{
		{
			name:   "success",
			filter: AssetFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAllAssets", mock.Anything, mock.AnythingOfType("asset.AssetFilter")).
					Return([]Asset{
						{ID: 1, Name: "MacBook Pro", AssetCategory: AssetCategory{ID: 1, Name: "Laptop"}, Status: constants.AssetStatusAvailable, Condition: constants.AssetConditionGood, CreatedAt: time.Now()},
					}, int64(1), nil)
				excel.On("GenerateSimpleExcel", "Assets", mock.Anything, mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "error repo fails",
			filter: AssetFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo, excel *mockExcel) {
				repo.On("FindAllAssets", mock.Anything, mock.AnythingOfType("asset.AssetFilter")).
					Return([]Asset(nil), int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, excel := newTestAssetService()
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
