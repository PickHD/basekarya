package subscription

import (
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/utils"
	"context"
	"errors"
	"fmt"
	"time"
)

type Service interface {
	ListPlans(ctx context.Context) ([]PlanResponse, error)
	RequestUpgrade(ctx context.Context, req *UpgradeRequest) (*UpgradeResponse, error)
	ListPendingRequests(ctx context.Context) ([]SubscriptionRequestResponse, error)
	ListAllRequests(ctx context.Context) ([]SubscriptionRequestResponse, error)
	ReviewRequest(ctx context.Context, requestID uint, req *ReviewRequest) error
	ListCompanies(ctx context.Context, search string) ([]CompanyListItem, error)
	GetCompanyDetail(ctx context.Context, id uint) (*CompanyDetail, error)
	UpdateCompanyStatus(ctx context.Context, companyID uint, req *UpdateCompanyStatusRequest) error
	RefreshCompanyCache(ctx context.Context, companyID uint) error
	GetDashboardStats(ctx context.Context) (*DashboardStatsResponse, error)
}

type service struct {
	repo    Repository
	company company.Repository
	role    RoleProvider
	user    UserProvider
	cache   CacheProvider
}

func NewService(repo Repository, companyRepo company.Repository, role RoleProvider, user UserProvider, cache CacheProvider) Service {
	return &service{repo: repo, company: companyRepo, role: role, user: user, cache: cache}
}

func (s *service) ListPlans(ctx context.Context) ([]PlanResponse, error) {
	plans, err := s.repo.FindAllPlans(ctx)
	if err != nil {
		return nil, err
	}

	var result []PlanResponse
	for _, p := range plans {
		result = append(result, PlanResponse{
			ID:           p.ID,
			Name:         p.Name,
			Slug:         p.Slug,
			MaxEmployees: p.MaxEmployees,
			PriceMonthly: p.PriceMonthly,
			Features:     p.Features,
		})
	}
	return result, nil
}

func (s *service) RequestUpgrade(ctx context.Context, req *UpgradeRequest) (*UpgradeResponse, error) {
	companyID := utils.GetCompanyIDFromCtx(ctx)
	userID := utils.GetUserIDFromCtx(ctx)

	if companyID == 0 {
		return nil, errors.New("company not found in context")
	}

	targetPlan, err := s.repo.FindPlanBySlug(ctx, req.PlanSlug)
	if err != nil {
		return nil, errors.New("plan not found")
	}

	comp, err := s.company.FindByID(ctx, companyID)
	if err != nil {
		return nil, errors.New("company not found")
	}

	if comp.SubscriptionPlanID != nil && *comp.SubscriptionPlanID == targetPlan.ID {
		return nil, errors.New("already on this plan")
	}

	if comp.SubscriptionPlanID != nil {
		currentPlan, _ := s.repo.FindPlanByID(ctx, *comp.SubscriptionPlanID)
		if currentPlan != nil && currentPlan.PriceMonthly >= targetPlan.PriceMonthly {
			return nil, errors.New("can only upgrade to a higher plan")
		}
	}

	pending, _ := s.repo.FindPendingRequestByCompanyID(ctx, companyID)
	if pending != nil && pending.ID != 0 {
		return nil, errors.New("you already have a pending upgrade request")
	}

	currentPlanID := uint(0)
	if comp.SubscriptionPlanID != nil {
		currentPlanID = *comp.SubscriptionPlanID
	}

	subReq := &SubscriptionRequest{
		CompanyID:       companyID,
		CurrentPlanID:   currentPlanID,
		RequestedPlanID: targetPlan.ID,
		Status:          constants.SubReqStatusPending,
		RequestedBy:     &userID,
	}
	if err := s.repo.CreateRequest(ctx, subReq); err != nil {
		return nil, fmt.Errorf("failed to create upgrade request: %w", err)
	}

	return &UpgradeResponse{
		ID:              subReq.ID,
		RequestedPlanID: targetPlan.ID,
		Status:          subReq.Status,
	}, nil
}

func (s *service) ListPendingRequests(ctx context.Context) ([]SubscriptionRequestResponse, error) {
	return s.repo.FindAllPendingRequests(ctx)
}

func (s *service) ReviewRequest(ctx context.Context, requestID uint, req *ReviewRequest) error {
	subReq, err := s.repo.FindRequestByID(ctx, requestID)
	if err != nil {
		return errors.New("request not found")
	}

	if subReq.Status != constants.SubReqStatusPending {
		return errors.New("request already reviewed")
	}

	reviewerID := utils.GetUserIDFromCtx(ctx)
	now := time.Now()

	subReq.Status = req.Status
	subReq.ReviewedBy = &reviewerID
	subReq.ReviewedAt = &now
	subReq.Notes = req.Notes

	if err := s.repo.UpdateRequest(ctx, subReq); err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	if req.Status == constants.SubReqStatusApproved {
		comp, err := s.company.FindByID(ctx, subReq.CompanyID)
		if err != nil {
			return errors.New("company not found")
		}

		comp.SubscriptionPlanID = &subReq.RequestedPlanID
		comp.SubscriptionStatus = constants.SubStatusActive
		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		comp.SubscriptionExpiresAt = &expiresAt

		if err := s.company.Update(ctx, comp); err != nil {
			return fmt.Errorf("failed to update company plan: %w", err)
		}

		plan, err := s.repo.FindPlanByID(ctx, subReq.RequestedPlanID)
		if err == nil && plan != nil {
			allowedGroups := buildAllowedGroups(plan.Slug)
			permissionIDs, err := s.role.FindPermissionIDsByGroupNames(ctx, allowedGroups)
			if err == nil && len(permissionIDs) > 0 {
				roleIDs, _ := s.role.FindRoleIDsByCompanyID(ctx, subReq.CompanyID)
				for _, roleID := range roleIDs {
					_ = s.role.AssignPermissions(ctx, roleID, permissionIDs, subReq.CompanyID)
				}
			}
		}

		_ = s.cache.Del(ctx, fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, subReq.CompanyID))
		_ = s.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, subReq.CompanyID))

		_ = s.user.ForceResetPasswordByCompanyID(ctx, subReq.CompanyID)
	}

	return nil
}

func (s *service) ListAllRequests(ctx context.Context) ([]SubscriptionRequestResponse, error) {
	return s.repo.FindAllRequests(ctx)
}

func (s *service) ListCompanies(ctx context.Context, search string) ([]CompanyListItem, error) {
	return s.repo.FindAllCompanies(ctx, search)
}

func (s *service) GetCompanyDetail(ctx context.Context, id uint) (*CompanyDetail, error) {
	return s.repo.FindCompanyDetailByID(ctx, id)
}

func (s *service) UpdateCompanyStatus(ctx context.Context, companyID uint, req *UpdateCompanyStatusRequest) error {
	if err := s.repo.UpdateCompanyStatus(ctx, companyID, req.SubscriptionStatus); err != nil {
		return err
	}

	_ = s.cache.Del(ctx, fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID))
	_ = s.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, companyID))

	return nil
}

func (s *service) GetDashboardStats(ctx context.Context) (*DashboardStatsResponse, error) {
	return s.repo.GetDashboardStats(ctx)
}

func (s *service) RefreshCompanyCache(ctx context.Context, companyID uint) error {
	s.cache.Del(ctx, fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID))
	s.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, companyID))
	return nil
}

func buildAllowedGroups(planSlug string) []string {
	groups := make([]string, len(constants.AlwaysAvailableGroups))
	copy(groups, constants.AlwaysAvailableGroups)

	planModules, ok := constants.PlanModules[planSlug]
	if !ok {
		return groups
	}

	for _, mod := range planModules {
		if modGroups, exists := constants.ModulePermissionGroups[mod]; exists {
			groups = append(groups, modGroups...)
		}
	}

	return groups
}
