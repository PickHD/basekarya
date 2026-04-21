package onboarding

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"time"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
)

type Service interface {
	// Templates
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest) error
	GetTemplates(ctx context.Context) ([]TemplateResponse, error)
	GetTemplateByID(ctx context.Context, id uint) (*TemplateResponse, error)
	UpdateTemplate(ctx context.Context, id uint, req *UpdateTemplateRequest) error
	DeleteTemplate(ctx context.Context, id uint) error

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
	transaction  infrastructure.TransactionManager
}

func NewService(
	repo Repository,
	notification NotificationProvider,
	user UserProvider,
	email EmailProvider,
	company CompanyProvider,
	transaction infrastructure.TransactionManager,
) Service {
	return &service{repo, notification, user, email, company, transaction}
}

// ── Templates ─────────────────────────────────────────────────────────────────

func (s *service) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) error {
	return s.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		t := &OnboardingTemplate{
			Name:       req.Name,
			Department: req.Department,
		}

		for _, item := range req.Items {
			t.Items = append(t.Items, OnboardingTemplateItem{
				TaskName:    item.TaskName,
				Description: item.Description,
				SortOrder:   item.SortOrder,
			})
		}

		err := s.repo.CreateTemplate(ctx, t)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetTemplates(ctx context.Context) ([]TemplateResponse, error) {
	templates, err := s.repo.FindAllTemplates(ctx)
	if err != nil {
		return nil, err
	}
	var result []TemplateResponse
	for _, t := range templates {
		result = append(result, toTemplateResponse(t))
	}
	return result, nil
}

func (s *service) GetTemplateByID(ctx context.Context, id uint) (*TemplateResponse, error) {
	t, err := s.repo.FindTemplateByID(ctx, id)
	if err != nil {
		return nil, errors.New("template not found")
	}
	r := toTemplateResponse(*t)
	return &r, nil
}

func (s *service) UpdateTemplate(ctx context.Context, id uint, req *UpdateTemplateRequest) error {
	return s.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		existing, err := s.repo.FindTemplateByID(ctx, id)
		if err != nil {
			return errors.New("template not found")
		}

		existing.Name = req.Name
		existing.Department = req.Department
		existing.Items = nil
		for _, item := range req.Items {
			existing.Items = append(existing.Items, OnboardingTemplateItem{
				TemplateID:  id,
				TaskName:    item.TaskName,
				Description: item.Description,
				SortOrder:   item.SortOrder,
			})
		}

		err = s.repo.UpdateTemplate(ctx, existing)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) DeleteTemplate(ctx context.Context, id uint) error {
	return s.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repo.FindTemplateByID(ctx, id)
		if err != nil {
			return errors.New("template not found")
		}

		err = s.repo.DeleteTemplate(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

// ── Workflows ─────────────────────────────────────────────────────────────────

func (s *service) CreateWorkflow(ctx context.Context, req *CreateWorkflowRequest) error {
	return s.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		workflow := &OnboardingWorkflow{
			ApplicantID:  req.ApplicantID,
			EmployeeID:   req.EmployeeID,
			NewHireName:  req.NewHireName,
			NewHireEmail: req.NewHireEmail,
			Position:     req.Position,
			Department:   req.Department,
			Status:       WorkflowStatusInProgress,
		}

		if req.StartDate != "" {
			t, err := time.Parse("2006-01-02", req.StartDate)
			if err == nil {
				workflow.StartDate = &t
			}
		}

		if err := s.repo.CreateWorkflow(ctx, workflow); err != nil {
			return err
		}

		// Copy all template items into workflow tasks
		templates, err := s.repo.FindAllTemplates(ctx)
		if err != nil {
			logger.Errorf("onboarding: failed to load templates: %v", err)
		} else {
			var tasks []OnboardingTask
			for _, tmpl := range templates {
				for _, item := range tmpl.Items {
					itemID := item.ID
					tasks = append(tasks, OnboardingTask{
						OnboardingWorkflowID: workflow.ID,
						TemplateItemID:       &itemID,
						TaskName:             item.TaskName,
						Description:          item.Description,
						Department:           tmpl.Department,
						SortOrder:            item.SortOrder,
					})
				}
			}
			if err := s.repo.CreateTasks(ctx, tasks); err != nil {
				logger.Errorf("onboarding: failed to create tasks: %v", err)
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
		ITTasks:          []TaskResponse{},
		HRTasks:          []TaskResponse{},
		OtherTasks:       []TaskResponse{},
	}

	for _, t := range w.Tasks {
		tr := toTaskResponse(t)
		switch t.Department {
		case "IT":
			detail.ITTasks = append(detail.ITTasks, tr)
		case "HR":
			detail.HRTasks = append(detail.HRTasks, tr)
		default:
			detail.OtherTasks = append(detail.OtherTasks, tr)
		}
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

func toTemplateResponse(t OnboardingTemplate) TemplateResponse {
	r := TemplateResponse{
		ID:         t.ID,
		Name:       t.Name,
		Department: t.Department,
		CreatedAt:  t.CreatedAt,
		Items:      []TemplateItemResponse{},
	}
	for _, item := range t.Items {
		r.Items = append(r.Items, TemplateItemResponse{
			ID:          item.ID,
			TaskName:    item.TaskName,
			Description: item.Description,
			SortOrder:   item.SortOrder,
		})
	}
	return r
}

func toTaskResponse(t OnboardingTask) TaskResponse {
	tr := TaskResponse{
		ID:          t.ID,
		TaskName:    t.TaskName,
		Description: t.Description,
		Department:  t.Department,
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
