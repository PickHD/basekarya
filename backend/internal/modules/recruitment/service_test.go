package recruitment

import (
	"errors"
	"testing"
	"time"

	"basekarya-backend/internal/modules/department"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateRequisition(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CreateRequisitionRequest
		setupMocks func(*mockRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req: &CreateRequisitionRequest{
				DepartmentID:   1,
				Title:          "Senior Go Developer",
				Description:    "Build services",
				Quantity:       2,
				EmploymentType: "PKWTT",
				Priority:       "HIGH",
				TargetDate:     "2026-07-01",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateRequisition", mock.Anything, mock.AnythingOfType("*recruitment.JobRequisition")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success without target date",
			req: &CreateRequisitionRequest{
				DepartmentID:   1,
				Title:          "Junior Dev",
				Quantity:       1,
				EmploymentType: "PKWT",
				Priority:       "MEDIUM",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateRequisition", mock.Anything, mock.AnythingOfType("*recruitment.JobRequisition")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error invalid target date",
			req: &CreateRequisitionRequest{
				DepartmentID:   1,
				Title:          "Dev",
				Quantity:       1,
				EmploymentType: "PKWTT",
				Priority:       "MEDIUM",
				TargetDate:     "invalid-date",
			},
			setupMocks: func(repo *mockRepo) {},
			wantErr:    true,
			errMsg:     "invalid target_date format (expected YYYY-MM-DD)",
		},
		{
			name: "error repo create fails",
			req: &CreateRequisitionRequest{
				DepartmentID:   1,
				Title:          "Dev",
				Quantity:       1,
				EmploymentType: "PKWTT",
				Priority:       "MEDIUM",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("CreateRequisition", mock.Anything, mock.AnythingOfType("*recruitment.JobRequisition")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			err := svc.CreateRequisition(ctx, 1, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_SubmitRequisition(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		id         uint
		setupMocks func(*mockRepo, *mockUserProvider, *mockNotification)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, RequesterID: 1, Status: constants.RequisitionStatusDraft, Title: "Dev",
				}, nil)
				repo.On("UpdateRequisitionStatus", mock.Anything, uint(1), constants.RequisitionStatusPending, (*uint)(nil), "").Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, string(constants.APPROVAL_REQUISITION)).Return([]uint{10, 11}, nil)
				notif.On("BlastNotification", mock.Anything, []uint{10, 11}, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   99,
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "requisition not found",
		},
		{
			name: "error not requester",
			id:   1,
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, RequesterID: 2, Status: constants.RequisitionStatusDraft,
				}, nil)
			},
			wantErr: true,
			errMsg:  "you are not the requester of this requisition",
		},
		{
			name: "error not draft status",
			id:   1,
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, RequesterID: 1, Status: constants.RequisitionStatusPending,
				}, nil)
			},
			wantErr: true,
			errMsg:  "requisition cannot be submitted from status 'PENDING'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, notif, userProv, _ := newTestService()
			tt.setupMocks(repo, userProv, notif)

			err := svc.SubmitRequisition(ctx, tt.id, 1)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_RequisitionAction(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		id         uint
		req        *RequisitionActionRequest
		setupMocks func(*mockRepo, *mockNotification)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "approve success",
			id:   1,
			req:  &RequisitionActionRequest{Action: "APPROVE"},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, RequesterID: 1, Status: constants.RequisitionStatusPending, Title: "Dev",
				}, nil)
				repo.On("UpdateRequisitionStatus", mock.Anything, uint(1), constants.RequisitionStatusApproved, mock.AnythingOfType("*uint"), "").Return(nil)
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "reject success",
			id:   2,
			req:  &RequisitionActionRequest{Action: "REJECT", RejectionReason: "Not enough budget"},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(2)).Return(&JobRequisition{
					ID: 2, RequesterID: 1, Status: constants.RequisitionStatusPending, Title: "Dev",
				}, nil)
				repo.On("UpdateRequisitionStatus", mock.Anything, uint(2), constants.RequisitionStatusRejected, mock.AnythingOfType("*uint"), "Not enough budget").Return(nil)
				notif.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   99,
			req:  &RequisitionActionRequest{Action: "APPROVE"},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "requisition not found",
		},
		{
			name: "error not pending",
			id:   1,
			req:  &RequisitionActionRequest{Action: "APPROVE"},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, Status: constants.RequisitionStatusDraft,
				}, nil)
			},
			wantErr: true,
			errMsg:  "requisition is not in PENDING status (current: DRAFT)",
		},
		{
			name: "error rejection reason required",
			id:   1,
			req:  &RequisitionActionRequest{Action: "REJECT", RejectionReason: ""},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, Status: constants.RequisitionStatusPending, RequesterID: 1,
				}, nil)
			},
			wantErr: true,
			errMsg:  "rejection reason is required",
		},
		{
			name: "error invalid action",
			id:   1,
			req:  &RequisitionActionRequest{Action: "INVALID"},
			setupMocks: func(repo *mockRepo, notif *mockNotification) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, Status: constants.RequisitionStatusPending, RequesterID: 1,
				}, nil)
			},
			wantErr: true,
			errMsg:  "invalid action: INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, notif, _, _ := newTestService()
			tt.setupMocks(repo, notif)

			err := svc.RequisitionAction(ctx, tt.id, 10, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetRequisitions(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     *RequisitionFilter
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:   "success with data",
			filter: &RequisitionFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllRequisitions", mock.Anything, mock.AnythingOfType("*recruitment.RequisitionFilter")).
					Return([]JobRequisition{
						{ID: 1, Title: "Dev", Status: "DRAFT", RequesterID: 1,
							Requester:  &user.User{Employee: &user.Employee{FullName: "John Doe"}},
							Department: &department.Department{Name: "Engineering"}},
					}, int64(1), nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "success empty",
			filter: &RequisitionFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllRequisitions", mock.Anything, mock.AnythingOfType("*recruitment.RequisitionFilter")).
					Return([]JobRequisition{}, int64(0), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "repo error",
			filter: &RequisitionFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllRequisitions", mock.Anything, mock.AnythingOfType("*recruitment.RequisitionFilter")).
					Return([]JobRequisition(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:   "defaults page and limit",
			filter: &RequisitionFilter{},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllRequisitions", mock.Anything, mock.MatchedBy(func(f *RequisitionFilter) bool {
					return f.Page == 1 && f.Limit == 10
				})).Return([]JobRequisition{}, int64(0), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			list, meta, err := svc.GetRequisitions(ctx, tt.filter)

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

func TestService_GetRequisitionDetail(t *testing.T) {
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
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, Title: "Dev", Status: "DRAFT", RequesterID: 1,
					Requester:  &user.User{Employee: &user.Employee{FullName: "John Doe"}},
					Approver:   &user.User{Employee: &user.Employee{FullName: "Admin"}},
					Department: &department.Department{Name: "Engineering"},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequisitionByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "requisition not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			result, err := svc.GetRequisitionDetail(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
				assert.Equal(t, "John Doe", result.RequesterName)
				assert.Equal(t, "Admin", result.ApproverName)
				assert.Equal(t, "Engineering", result.DepartmentName)
			}
		})
	}
}

func TestService_CloseRequisition(t *testing.T) {
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
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, Status: constants.RequisitionStatusApproved,
				}, nil)
				repo.On("UpdateRequisitionStatus", mock.Anything, uint(1), constants.RequisitionStatusClosed, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequisitionByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "requisition not found",
		},
		{
			name: "error already closed",
			id:   1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					ID: 1, Status: constants.RequisitionStatusClosed,
				}, nil)
			},
			wantErr: true,
			errMsg:  "requisition is already closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			err := svc.CloseRequisition(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_DeleteRequisition(t *testing.T) {
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
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{ID: 1}, nil)
				repo.On("SoftDeleteRequisition", mock.Anything, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequisitionByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "requisition not found",
		},
		{
			name: "error delete fails",
			id:   1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{ID: 1}, nil)
				repo.On("SoftDeleteRequisition", mock.Anything, uint(1)).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			err := svc.DeleteRequisition(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_AddApplicant(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		reqID      uint
		req        *CreateApplicantRequest
		setupMocks func(*mockRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name:  "success",
			reqID: 1,
			req: &CreateApplicantRequest{
				FullName: "Jane Smith", Email: "jane@example.com", PhoneNumber: "0812345678",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{ID: 1}, nil)
				repo.On("CountApplicantsByRequisitionAndStage", mock.Anything, uint(1), constants.ApplicantStageScreening).Return(int64(0), nil)
				repo.On("CreateApplicant", mock.Anything, mock.AnythingOfType("*recruitment.Applicant")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "error requisition not found",
			reqID: 99,
			req: &CreateApplicantRequest{
				FullName: "Jane Smith", Email: "jane@example.com",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequisitionByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "requisition not found",
		},
		{
			name:  "error create fails",
			reqID: 1,
			req: &CreateApplicantRequest{
				FullName: "Jane Smith", Email: "jane@example.com",
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{ID: 1}, nil)
				repo.On("CountApplicantsByRequisitionAndStage", mock.Anything, uint(1), constants.ApplicantStageScreening).Return(int64(0), nil)
				repo.On("CreateApplicant", mock.Anything, mock.AnythingOfType("*recruitment.Applicant")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			err := svc.AddApplicant(ctx, tt.reqID, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_UpdateStage(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	targetDate := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		id         uint
		req        *UpdateApplicantStageRequest
		setupMocks func(*mockRepo, *mockOnboarding)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success interview",
			id:   1,
			req:  &UpdateApplicantStageRequest{Stage: constants.ApplicantStageInterview, Notes: "Passed screening"},
			setupMocks: func(repo *mockRepo, onboard *mockOnboarding) {
				repo.On("FindApplicantByID", mock.Anything, uint(1)).Return(&Applicant{
					ID: 1, JobRequisitionID: 1, Stage: constants.ApplicantStageScreening, CompanyID: 1,
				}, nil)
				repo.On("CountApplicantsByRequisitionAndStage", mock.Anything, uint(1), constants.ApplicantStageInterview).Return(int64(0), nil)
				repo.On("UpdateApplicantStage", mock.Anything, uint(1), constants.ApplicantStageInterview, 0, "Passed screening", "").Return(nil)
				repo.On("CreateStageHistory", mock.Anything, mock.AnythingOfType("*recruitment.ApplicantStageHistory")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success hired triggers onboarding",
			id:   1,
			req:  &UpdateApplicantStageRequest{Stage: constants.ApplicantStageHired},
			setupMocks: func(repo *mockRepo, onboard *mockOnboarding) {
				repo.On("FindApplicantByID", mock.Anything, uint(1)).Return(&Applicant{
					ID: 1, JobRequisitionID: 1, Stage: constants.ApplicantStageOffering,
					FullName: "Jane", Email: "jane@example.com", CompanyID: 1,
					JobRequisition: &JobRequisition{
						Title: "Dev", TargetDate: &targetDate,
						Department: &department.Department{Name: "Engineering"},
					},
				}, nil)
				repo.On("CountApplicantsByRequisitionAndStage", mock.Anything, uint(1), constants.ApplicantStageHired).Return(int64(0), nil)
				repo.On("UpdateApplicantStage", mock.Anything, uint(1), constants.ApplicantStageHired, 0, "", "").Return(nil)
				repo.On("FindRequisitionByID", mock.Anything, uint(1)).Return(&JobRequisition{
					Title: "Dev", TargetDate: &targetDate,
					Department: &department.Department{Name: "Engineering"},
				}, nil)
				onboard.On("CreateWorkflow", mock.Anything, mock.AnythingOfType("*onboarding.CreateWorkflowRequest")).Return(nil)
				repo.On("CreateStageHistory", mock.Anything, mock.AnythingOfType("*recruitment.ApplicantStageHistory")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error applicant not found",
			id:   99,
			req:  &UpdateApplicantStageRequest{Stage: constants.ApplicantStageInterview},
			setupMocks: func(repo *mockRepo, onboard *mockOnboarding) {
				repo.On("FindApplicantByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "applicant not found",
		},
		{
			name: "error update stage fails",
			id:   1,
			req:  &UpdateApplicantStageRequest{Stage: constants.ApplicantStageInterview},
			setupMocks: func(repo *mockRepo, onboard *mockOnboarding) {
				repo.On("FindApplicantByID", mock.Anything, uint(1)).Return(&Applicant{
					ID: 1, JobRequisitionID: 1, Stage: constants.ApplicantStageScreening, CompanyID: 1,
				}, nil)
				repo.On("CountApplicantsByRequisitionAndStage", mock.Anything, uint(1), constants.ApplicantStageInterview).Return(int64(0), nil)
				repo.On("UpdateApplicantStage", mock.Anything, uint(1), constants.ApplicantStageInterview, 0, "", "").Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, onboard := newTestService()
			tt.setupMocks(repo, onboard)

			err := svc.UpdateStage(ctx, tt.id, 1, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetApplicantsByRequisition(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		reqID      uint
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name:  "success with applicants",
			reqID: 1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindApplicantsByRequisitionID", mock.Anything, uint(1)).Return([]Applicant{
					{ID: 1, FullName: "Jane", Stage: constants.ApplicantStageScreening, StageOrder: 0},
					{ID: 2, FullName: "Bob", Stage: constants.ApplicantStageInterview, StageOrder: 0},
				}, nil)
			},
			wantErr: false,
		},
		{
			name:  "success empty",
			reqID: 1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindApplicantsByRequisitionID", mock.Anything, uint(1)).Return([]Applicant{}, nil)
			},
			wantErr: false,
		},
		{
			name:  "repo error",
			reqID: 1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindApplicantsByRequisitionID", mock.Anything, uint(1)).Return([]Applicant(nil), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			board, err := svc.GetApplicantsByRequisition(ctx, tt.reqID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, board)
			}
		})
	}
}

func TestService_GetApplicantDetail(t *testing.T) {
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
				repo.On("FindApplicantByID", mock.Anything, uint(1)).Return(&Applicant{
					ID: 1, FullName: "Jane", Email: "jane@example.com",
					Stage: constants.ApplicantStageScreening,
					StageHistories: []ApplicantStageHistory{
						{
							ID: 1, FromStage: constants.ApplicantStageScreening,
							ToStage: constants.ApplicantStageInterview,
							ChangedByUser: &user.User{Employee: &user.Employee{FullName: "Admin"}},
						},
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindApplicantByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "applicant not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestService()
			tt.setupMocks(repo)

			result, err := svc.GetApplicantDetail(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
				assert.Len(t, result.StageHistories, 1)
			}
		})
	}
}
