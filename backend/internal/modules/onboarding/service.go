package onboarding

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"strings"
	"time"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
)

type Service interface {
	// Workflows
	CreateWorkflow(ctx context.Context, req *CreateWorkflowRequest) error
	GetWorkflows(ctx context.Context, filter *WorkflowFilter) ([]WorkflowListResponse, *response.Meta, error)
	GetWorkflowDetail(ctx context.Context, id uint) (*WorkflowDetailResponse, error)

	// Tasks
	CompleteTask(ctx context.Context, taskID uint, completedByID uint, req *CompleteTaskRequest) error
}

type service struct {
	repo         Repository
	notification NotificationProvider
	user         UserProvider
	email        EmailProvider
	company      CompanyProvider
	role         RoleProvider
	department   DepartmentProvider
	master       MasterProvider
	transaction  infrastructure.TransactionManager
}

func NewService(
	repo Repository,
	notification NotificationProvider,
	user UserProvider,
	email EmailProvider,
	company CompanyProvider,
	role RoleProvider,
	department DepartmentProvider,
	master MasterProvider,
	transaction infrastructure.TransactionManager,
) Service {
	return &service{repo, notification, user, email, company, role, department, master, transaction}
}


// ── Workflows ─────────────────────────────────────────────────────────────────

func (s *service) CreateWorkflow(ctx context.Context, req *CreateWorkflowRequest) error {
	return s.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		workflow := &OnboardingWorkflow{
			CompanyID:    utils.GetCompanyIDFromCtx(ctx),
			ApplicantID:  req.ApplicantID,
			EmployeeID:   req.EmployeeID,
			NewHireName:  req.NewHireName,
			NewHireEmail: req.NewHireEmail,
			Position:     req.Position,
			Department:   req.Department,
			Status:       WorkflowStatusInProgress,
		}

		if req.StartDate != "" {
			t, err := time.Parse(constants.DefaultTimeFormat, req.StartDate)
			if err == nil {
				workflow.StartDate = &t
			}
		}

		if err := s.repo.CreateWorkflow(ctx, workflow); err != nil {
			return err
		}

		// Create tasks from request
		if len(req.Tasks) > 0 {
			var tasks []OnboardingTask
			for _, t := range req.Tasks {
				tasks = append(tasks, OnboardingTask{
					CompanyID:            utils.GetCompanyIDFromCtx(ctx),
					OnboardingWorkflowID: workflow.ID,
					TaskName:             t.TaskName,
					Description:          t.Description,
					SortOrder:            t.SortOrder,
				})
			}
			if err := s.repo.CreateTasks(ctx, tasks); err != nil {
				return err
			}
		}

		// Send welcome email asynchronously
		go func() {
			if err := s.sendWelcomeEmail(workflow); err != nil {
				logger.Errorf("onboarding: welcome email failed for %s: %v", workflow.NewHireEmail, err)
				return
			}
			_ = s.repo.MarkWorkflowEmailSent(context.Background(), workflow.ID)
		}()

		// Notify IT and HR users about pending tasks
		go func() {
			itUsers, _ := s.user.FindApprovalUsers(context.Background(), constants.UPDATE_ONBOARDING_TASK)
			if len(itUsers) > 0 {
			_ = s.notification.BlastNotification(
				utils.DetachContext(ctx),
				itUsers,
					string(constants.NotificationTypeOnboardingTask),
					"New Onboarding Tasks",
					"New hire "+req.NewHireName+" has joined. Please complete onboarding tasks.",
					workflow.ID,
				)
			}
		}()

		return nil
	})
}

func (s *service) GetWorkflows(ctx context.Context, filter *WorkflowFilter) ([]WorkflowListResponse, *response.Meta, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	workflows, total, err := s.repo.FindAllWorkflows(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var result []WorkflowListResponse
	for _, w := range workflows {
		total_tasks := len(w.Tasks)
		completed := 0
		for _, t := range w.Tasks {
			if t.IsCompleted {
				completed++
			}
		}
		progress := 0
		if total_tasks > 0 {
			progress = (completed * 100) / total_tasks
		}

		result = append(result, WorkflowListResponse{
			ID:           w.ID,
			NewHireName:  w.NewHireName,
			NewHireEmail: w.NewHireEmail,
			Position:     w.Position,
			Department:   w.Department,
			StartDate:    w.StartDate,
			Status:       w.Status,
			Progress:     progress,
			CreatedAt:    w.CreatedAt,
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return result, meta, nil
}

func (s *service) GetWorkflowDetail(ctx context.Context, id uint) (*WorkflowDetailResponse, error) {
	w, err := s.repo.FindWorkflowByID(ctx, id)
	if err != nil {
		return nil, errors.New("workflow not found")
	}

	total_tasks := len(w.Tasks)
	completed := 0
	for _, t := range w.Tasks {
		if t.IsCompleted {
			completed++
		}
	}
	progress := 0
	if total_tasks > 0 {
		progress = (completed * 100) / total_tasks
	}

	detail := &WorkflowDetailResponse{
		ID:               w.ID,
		NewHireName:      w.NewHireName,
		NewHireEmail:     w.NewHireEmail,
		Position:         w.Position,
		Department:       w.Department,
		StartDate:        w.StartDate,
		Status:           w.Status,
		Progress:         progress,
		WelcomeEmailSent: w.WelcomeEmailSent,
		CreatedAt:        w.CreatedAt,
		Tasks:            []TaskResponse{},
	}

	for _, t := range w.Tasks {
		detail.Tasks = append(detail.Tasks, toTaskResponse(t))
	}

	return detail, nil
}

// ── Tasks ─────────────────────────────────────────────────────────────────────

func (s *service) CompleteTask(ctx context.Context, taskID uint, completedByID uint, req *CompleteTaskRequest) error {
	return s.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		task, err := s.repo.FindTaskByID(ctx, taskID)
		if err != nil {
			return errors.New("task not found")
		}
		if task.IsCompleted {
			return errors.New("task is already completed")
		}

		if err := s.repo.CompleteTask(ctx, taskID, completedByID, req.Notes); err != nil {
			return err
		}

		// Check if all tasks are done → complete the workflow
		pending, errPending := s.repo.CountPendingTasks(ctx, task.OnboardingWorkflowID)
		if errPending == nil && pending == 0 {
			errMark := s.repo.MarkWorkflowCompleted(ctx, task.OnboardingWorkflowID)
			if errMark != nil {
				return errMark
			}

			w, err := s.repo.FindWorkflowByID(ctx, task.OnboardingWorkflowID)
			if err != nil {
				return err
			}

			role, err := s.role.FindRoleByName(ctx, "EMPLOYEE")
			if err != nil {
				return err
			}

			department, err := s.department.FindByName(ctx, "Umum")
			if err != nil {
				return err
			}

			shift, err := s.master.FindShiftByName(ctx, "Regular")
			if err != nil {
				return err
			}

			req := &user.CreateEmployeeRequest{
				NIK:          strings.ReplaceAll(strings.ToLower(w.NewHireName), " ", ""),
				FullName:     w.NewHireName,
				Email:        w.NewHireEmail,
				Position:     w.Position,
				BaseSalary:   0.0,
				RoleID:       role.ID,
				DepartmentID: department.ID,
				ShiftID:      shift.ID,
			}

			if _, err := s.user.CreateEmployee(ctx, req); err != nil {
				return err
			}
		}

		return nil
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func (s *service) sendWelcomeEmail(w *OnboardingWorkflow) error {
	company, err := s.company.FindByID(context.Background(), 1)
	if err != nil {
		return err
	}

	startDate := "-"
	if w.StartDate != nil {
		startDate = w.StartDate.Format("02 January 2006")
	}

	data := struct {
		CompanyName string
		FullName    string
		Position    string
		Department  string
		StartDate   string
	}{
		CompanyName: company.Name,
		FullName:    w.NewHireName,
		Position:    w.Position,
		Department:  w.Department,
		StartDate:   startDate,
	}

	tmpl, err := template.New("welcome").Parse(constants.WelcomeEmailTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return s.email.Send(w.NewHireEmail, "Welcome to "+company.Name+"!", buf.String())
}

func toTaskResponse(t OnboardingTask) TaskResponse {
	tr := TaskResponse{
		ID:          t.ID,
		TaskName:    t.TaskName,
		Description: t.Description,
		IsCompleted: t.IsCompleted,
		CompletedAt: t.CompletedAt,
		Notes:       t.Notes,
		SortOrder:   t.SortOrder,
	}
	if t.CompletedByUser != nil && t.CompletedByUser.Employee != nil {
		tr.CompletedBy = t.CompletedByUser.Employee.FullName
	}
	return tr
}
