package user

import (
	"encoding/json"
	"errors"
	"testing"

	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetProfile(t *testing.T) {
	tests := []struct {
		name       string
		userID     uint
		setupMocks func(*mockRepo, *mockCache)
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "cache hit",
			userID: 1,
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				cachedData, _ := json.Marshal(&UserProfileResponse{
					ID:       1,
					Username: "john.doe",
					Role:     "EMPLOYEE",
					FullName: "John Doe",
				})
				cache.On("Get", mock.Anything, "user:1").Return(string(cachedData), nil)
			},
			wantErr: false,
		},
		{
			name:   "cache miss - db lookup success",
			userID: 2,
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				cache.On("Get", mock.Anything, "user:2").Return("", redis.Nil)
				repo.On("FindByID", mock.Anything, uint(2)).Return(&User{
					ID:                 2,
					Username:           "jane.doe",
					MustChangePassword: false,
					Role:               &rbac.Role{Name: "EMPLOYEE"},
					Employee: &Employee{
						FullName: "Jane Doe",
						NIK:      "EMP002",
						Department: &master.Department{Name: "Engineering"},
						Shift:      &master.Shift{Name: "Day", StartTime: "09:00", EndTime: "17:00"},
					},
				}, nil)
				cache.On("Set", mock.Anything, "user:2", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "cache miss - db lookup error",
			userID: 99,
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				cache.On("Get", mock.Anything, "user:99").Return("", redis.Nil)
				repo.On("FindByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name:   "cache miss - super admin no employee",
			userID: 1,
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				cache.On("Get", mock.Anything, "user:1").Return("", redis.Nil)
				repo.On("FindByID", mock.Anything, uint(1)).Return(&User{
					ID:                 1,
					Username:           "admin",
					MustChangePassword: false,
					Role:               &rbac.Role{Name: "SUPERADMIN"},
					Employee:           nil,
				}, nil)
				cache.On("Set", mock.Anything, "user:1", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "cache error",
			userID: 1,
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				cache.On("Get", mock.Anything, "user:1").Return("", errors.New("redis down"))
			},
			wantErr: true,
			errMsg:  "redis down",
		},
		{
			name:   "cache miss - set cache error",
			userID: 1,
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				cache.On("Get", mock.Anything, "user:1").Return("", redis.Nil)
				repo.On("FindByID", mock.Anything, uint(1)).Return(&User{
					ID:       1,
					Username: "john.doe",
					Role:     &rbac.Role{Name: "EMPLOYEE"},
					Employee: &Employee{FullName: "John Doe"},
				}, nil)
				cache.On("Set", mock.Anything, "user:1", mock.Anything, mock.Anything).Return(errors.New("set failed"))
			},
			wantErr: true,
			errMsg:  "set failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, cache, _, _, _ := newTestUserService()
			tt.setupMocks(repo, cache)

			resp, err := svc.GetProfile(tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.userID, resp.ID)
			}
		})
	}
}

func TestService_UpdateProfile(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		userID     uint
		req        *UpdateProfileRequest
		setupMocks func(*mockRepo, *mockStorage, *mockCache)
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "success",
			userID: 1,
			req: &UpdateProfileRequest{
				FullName:    "John Updated",
				PhoneNumber: "081234567890",
			},
			setupMocks: func(repo *mockRepo, storage *mockStorage, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&User{
					ID: 1,
					Employee: &Employee{
						FullName: "John Doe",
					},
				}, nil)
				repo.On("UpdateEmployee", mock.Anything, mock.AnythingOfType("*user.Employee")).Return(nil)
				cache.On("Del", mock.Anything, "user:1").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error user not found",
			userID: 99,
			req:    &UpdateProfileRequest{FullName: "Test"},
			setupMocks: func(repo *mockRepo, storage *mockStorage, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name:   "error no employee data",
			userID: 2,
			req:    &UpdateProfileRequest{FullName: "Test"},
			setupMocks: func(repo *mockRepo, storage *mockStorage, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&User{
					ID:       2,
					Employee: nil,
				}, nil)
			},
			wantErr: true,
			errMsg:  "employee data not found",
		},
		{
			name:   "error update employee fails",
			userID: 1,
			req:    &UpdateProfileRequest{FullName: "John Updated"},
			setupMocks: func(repo *mockRepo, storage *mockStorage, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&User{
					ID:       1,
					Employee: &Employee{FullName: "John Doe"},
				}, nil)
				repo.On("UpdateEmployee", mock.Anything, mock.AnythingOfType("*user.Employee")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name:   "error cache delete fails",
			userID: 1,
			req:    &UpdateProfileRequest{FullName: "John Updated"},
			setupMocks: func(repo *mockRepo, storage *mockStorage, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&User{
					ID:       1,
					Employee: &Employee{FullName: "John Doe"},
				}, nil)
				repo.On("UpdateEmployee", mock.Anything, mock.AnythingOfType("*user.Employee")).Return(nil)
				cache.On("Del", mock.Anything, "user:1").Return(errors.New("cache error"))
			},
			wantErr: true,
			errMsg:  "cache error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, storage, cache, _, _, _ := newTestUserService()
			tt.setupMocks(repo, storage, cache)

			err := svc.UpdateProfile(ctx, tt.userID, tt.req, nil)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_ChangePassword(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		userID     uint
		req        *ChangePasswordRequest
		setupMocks func(*mockRepo, *mockHasher, *mockCache)
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "success",
			userID: 1,
			req: &ChangePasswordRequest{
				OldPassword:     "oldpass",
				NewPassword:     "newpass",
				ConfirmPassword: "newpass",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&User{
					ID:           1,
					PasswordHash: "oldhash",
				}, nil)
				hasher.On("CheckPasswordHash", "oldpass", "oldhash").Return(true)
				hasher.On("HashPassword", "newpass").Return("newhash", nil)
				repo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
				cache.On("Del", mock.Anything, "user:1").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error user not found",
			userID: 99,
			req: &ChangePasswordRequest{
				OldPassword: "oldpass", NewPassword: "newpass", ConfirmPassword: "newpass",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name:   "error invalid old password",
			userID: 1,
			req: &ChangePasswordRequest{
				OldPassword: "wrongpass", NewPassword: "newpass", ConfirmPassword: "newpass",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&User{
					ID:           1,
					PasswordHash: "oldhash",
				}, nil)
				hasher.On("CheckPasswordHash", "wrongpass", "oldhash").Return(false)
			},
			wantErr: true,
			errMsg:  "invalid old password",
		},
		{
			name:   "error hash password fails",
			userID: 1,
			req: &ChangePasswordRequest{
				OldPassword: "oldpass", NewPassword: "newpass", ConfirmPassword: "newpass",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&User{
					ID:           1,
					PasswordHash: "oldhash",
				}, nil)
				hasher.On("CheckPasswordHash", "oldpass", "oldhash").Return(true)
				hasher.On("HashPassword", "newpass").Return("", errors.New("hash error"))
			},
			wantErr: true,
			errMsg:  "hash error",
		},
		{
			name:   "error update user fails",
			userID: 1,
			req: &ChangePasswordRequest{
				OldPassword: "oldpass", NewPassword: "newpass", ConfirmPassword: "newpass",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, cache *mockCache) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&User{
					ID:           1,
					PasswordHash: "oldhash",
				}, nil)
				hasher.On("CheckPasswordHash", "oldpass", "oldhash").Return(true)
				hasher.On("HashPassword", "newpass").Return("newhash", nil)
				repo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, hasher, _, cache, _, _, _ := newTestUserService()
			tt.setupMocks(repo, hasher, cache)

			err := svc.ChangePassword(ctx, tt.userID, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetAllEmployees(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		page       int
		limit      int
		search     string
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:   "success with data",
			page:   1,
			limit:  10,
			search: "",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllEmployees", mock.Anything, 1, 10, "").Return([]User{
					{
						ID:       1,
						Username: "john.doe",
						Role:     &rbac.Role{ID: 1, Name: "EMPLOYEE"},
						Employee: &Employee{
							ID: 1, FullName: "John Doe", NIK: "EMP001",
							Department: &master.Department{Name: "Engineering"},
							Shift:      &master.Shift{Name: "Day"},
							BaseSalary: 5000000,
							Email:      "john@example.com",
							Position:   "Developer",
						},
					},
				}, int64(1), nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "success empty list",
			page:   1,
			limit:  10,
			search: "",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllEmployees", mock.Anything, 1, 10, "").Return([]User{}, int64(0), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "error repo fails",
			page:   1,
			limit:  10,
			search: "",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllEmployees", mock.Anything, 1, 10, "").Return([]User(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _, _ := newTestUserService()
			tt.setupMocks(repo)

			list, meta, err := svc.GetAllEmployees(ctx, tt.page, tt.limit, tt.search)

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

func TestService_CreateEmployee(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CreateEmployeeRequest
		setupMocks func(*mockRepo, *mockHasher, *mockLeaveGen, *mockSubscription)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req: &CreateEmployeeRequest{
				FullName:     "Jane Doe",
				NIK:          "EMP002",
				DepartmentID: 1,
				ShiftID:      1,
				RoleID:       1,
				BaseSalary:   5000000,
				Email:        "jane@example.com",
				Position:     "Designer",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, leaveGen *mockLeaveGen, sub *mockSubscription) {
				sub.On("CheckEmployeeLimit", mock.Anything).Return(true, nil)
				hasher.On("HashPassword", "BaseKarya2024").Return("hashedpass", nil)
				repo.On("FindRoleByID", mock.Anything, uint(1)).Return(&rbac.Role{ID: 1, Name: "EMPLOYEE"}, nil)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
				repo.On("CreateEmployee", mock.Anything, mock.AnythingOfType("*user.Employee")).Return(nil)
				leaveGen.On("GenerateInitialBalance", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error subscription limit reached",
			req: &CreateEmployeeRequest{
				FullName: "Jane Doe", NIK: "EMP002",
				DepartmentID: 1, ShiftID: 1, RoleID: 1,
				BaseSalary: 5000000, Email: "jane@example.com", Position: "Designer",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, leaveGen *mockLeaveGen, sub *mockSubscription) {
				sub.On("CheckEmployeeLimit", mock.Anything).Return(false, nil)
			},
			wantErr: true,
			errMsg:  "employee limit reached for your subscription plan. please upgrade to add more employees",
		},
		{
			name: "error subscription check fails",
			req: &CreateEmployeeRequest{
				FullName: "Jane Doe", NIK: "EMP002",
				DepartmentID: 1, ShiftID: 1, RoleID: 1,
				BaseSalary: 5000000, Email: "jane@example.com", Position: "Designer",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, leaveGen *mockLeaveGen, sub *mockSubscription) {
				sub.On("CheckEmployeeLimit", mock.Anything).Return(false, errors.New("subscription error"))
			},
			wantErr: true,
			errMsg:  "failed to check employee limit: subscription error",
		},
		{
			name: "error role not found",
			req: &CreateEmployeeRequest{
				FullName: "Jane Doe", NIK: "EMP002",
				DepartmentID: 1, ShiftID: 1, RoleID: 99,
				BaseSalary: 5000000, Email: "jane@example.com", Position: "Designer",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, leaveGen *mockLeaveGen, sub *mockSubscription) {
				sub.On("CheckEmployeeLimit", mock.Anything).Return(true, nil)
				hasher.On("HashPassword", "BaseKarya2024").Return("hashedpass", nil)
				repo.On("FindRoleByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "role not found",
		},
		{
			name: "error create user fails",
			req: &CreateEmployeeRequest{
				FullName: "Jane Doe", NIK: "EMP002",
				DepartmentID: 1, ShiftID: 1, RoleID: 1,
				BaseSalary: 5000000, Email: "jane@example.com", Position: "Designer",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, leaveGen *mockLeaveGen, sub *mockSubscription) {
				sub.On("CheckEmployeeLimit", mock.Anything).Return(true, nil)
				hasher.On("HashPassword", "BaseKarya2024").Return("hashedpass", nil)
				repo.On("FindRoleByID", mock.Anything, uint(1)).Return(&rbac.Role{ID: 1}, nil)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name: "error create employee fails",
			req: &CreateEmployeeRequest{
				FullName: "Jane Doe", NIK: "EMP002",
				DepartmentID: 1, ShiftID: 1, RoleID: 1,
				BaseSalary: 5000000, Email: "jane@example.com", Position: "Designer",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, leaveGen *mockLeaveGen, sub *mockSubscription) {
				sub.On("CheckEmployeeLimit", mock.Anything).Return(true, nil)
				hasher.On("HashPassword", "BaseKarya2024").Return("hashedpass", nil)
				repo.On("FindRoleByID", mock.Anything, uint(1)).Return(&rbac.Role{ID: 1}, nil)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
				repo.On("CreateEmployee", mock.Anything, mock.AnythingOfType("*user.Employee")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name: "error leave generation fails",
			req: &CreateEmployeeRequest{
				FullName: "Jane Doe", NIK: "EMP002",
				DepartmentID: 1, ShiftID: 1, RoleID: 1,
				BaseSalary: 5000000, Email: "jane@example.com", Position: "Designer",
			},
			setupMocks: func(repo *mockRepo, hasher *mockHasher, leaveGen *mockLeaveGen, sub *mockSubscription) {
				sub.On("CheckEmployeeLimit", mock.Anything).Return(true, nil)
				hasher.On("HashPassword", "BaseKarya2024").Return("hashedpass", nil)
				repo.On("FindRoleByID", mock.Anything, uint(1)).Return(&rbac.Role{ID: 1}, nil)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
				repo.On("CreateEmployee", mock.Anything, mock.AnythingOfType("*user.Employee")).Return(nil)
				leaveGen.On("GenerateInitialBalance", mock.Anything, mock.Anything).Return(errors.New("leave error"))
			},
			wantErr: true,
			errMsg:  "leave error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, hasher, _, _, leaveGen, _, sub := newTestUserService()
			tt.setupMocks(repo, hasher, leaveGen, sub)

			resp, err := svc.CreateEmployee(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp.Username)
			}
		})
	}
}

func TestService_UpdateEmployee(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		id         uint
		req        *UpdateEmployeeRequest
		setupMocks func(*mockRepo, *mockCache)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			id:   1,
			req: &UpdateEmployeeRequest{
				FullName: "John Updated",
				Position: "Senior Developer",
			},
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				repo.On("FindEmployeeByID", mock.Anything, uint(1)).Return(&Employee{
					ID: 1, FullName: "John Doe", Position: "Developer", UserID: 10,
				}, nil)
				repo.On("UpdateEmployee", mock.Anything, mock.AnythingOfType("*user.Employee")).Return(nil)
				cache.On("Del", mock.Anything, "user:10").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success with role update",
			id:   1,
			req: &UpdateEmployeeRequest{
				FullName: "John Updated",
				RoleID:   2,
			},
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				repo.On("FindEmployeeByID", mock.Anything, uint(1)).Return(&Employee{
					ID: 1, FullName: "John Doe", Position: "Developer", UserID: 10,
					User: User{ID: 10, RoleID: 1},
				}, nil)
				repo.On("UpdateEmployee", mock.Anything, mock.AnythingOfType("*user.Employee")).Return(nil)
				repo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
				cache.On("Del", mock.Anything, "user:10").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error employee not found",
			id:   99,
			req:  &UpdateEmployeeRequest{FullName: "Test"},
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				repo.On("FindEmployeeByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "employee not found",
		},
		{
			name: "error update fails",
			id:   1,
			req:  &UpdateEmployeeRequest{FullName: "John Updated"},
			setupMocks: func(repo *mockRepo, cache *mockCache) {
				repo.On("FindEmployeeByID", mock.Anything, uint(1)).Return(&Employee{
					ID: 1, FullName: "John Doe",
				}, nil)
				repo.On("UpdateEmployee", mock.Anything, mock.AnythingOfType("*user.Employee")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, cache, _, _, _ := newTestUserService()
			tt.setupMocks(repo, cache)

			err := svc.UpdateEmployee(ctx, tt.id, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_DeleteEmployee(t *testing.T) {
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
				repo.On("FindEmployeeByID", mock.Anything, uint(1)).Return(&Employee{
					ID: 1, UserID: 1,
				}, nil)
				repo.On("DeleteUser", mock.Anything, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error employee not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindEmployeeByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "employee not found",
		},
		{
			name: "error delete user fails",
			id:   1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindEmployeeByID", mock.Anything, uint(1)).Return(&Employee{
					ID: 1, UserID: 1,
				}, nil)
				repo.On("DeleteUser", mock.Anything, uint(1)).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _, _ := newTestUserService()
			tt.setupMocks(repo)

			err := svc.DeleteEmployee(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
