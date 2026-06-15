package subscription

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PlanCacheService struct {
	db    *gorm.DB
	redis *infrastructure.RedisClientProvider
}

type planFeatures struct {
	Modules []string `json:"modules"`
}

func NewPlanCacheService(db *gorm.DB, redis *infrastructure.RedisClientProvider) *PlanCacheService {
	return &PlanCacheService{db: db, redis: redis}
}

func (s *PlanCacheService) HasAccess(ctx context.Context, companyID uint, module string) (bool, error) {
	key := fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID)

	featuresJSON, err := s.redis.Get(ctx, key)
	if err == nil && featuresJSON != "" {
		var features planFeatures
		if err := json.Unmarshal([]byte(featuresJSON), &features); err != nil {
			return false, fmt.Errorf("failed to parse cached features: %w", err)
		}
		return moduleInFeatures(features.Modules, module), nil
	}

	var dbJSON string
	err = s.db.Table("subscription_plans").
		Select("subscription_plans.features").
		Joins("JOIN companies ON companies.subscription_plan_id = subscription_plans.id").
		Where("companies.id = ?", companyID).
		Scan(&dbJSON).Error
	if err != nil || dbJSON == "" {
		return false, fmt.Errorf("subscription plan not found")
	}

	if setErr := s.redis.Set(ctx, key, dbJSON, 24*time.Hour); setErr != nil {
		logger.Errorw("PlanCacheService: failed to cache features", "key", key, "err", setErr)
	}

	var features planFeatures
	if err := json.Unmarshal([]byte(dbJSON), &features); err != nil {
		return false, fmt.Errorf("failed to parse features: %w", err)
	}
	return moduleInFeatures(features.Modules, module), nil
}

func (s *PlanCacheService) CheckEmployeeLimit(ctx context.Context) (bool, error) {
	companyID := utils.GetCompanyIDFromCtx(ctx)
	if companyID == 0 {
		return true, nil
	}

	var maxEmployees int
	err := s.db.Table("subscription_plans").
		Select("subscription_plans.max_employees").
		Joins("JOIN companies ON companies.subscription_plan_id = subscription_plans.id").
		Where("companies.id = ?", companyID).
		Scan(&maxEmployees).Error
	if err != nil {
		return true, err
	}

	if maxEmployees == 0 {
		return true, nil
	}

	var count int64
	s.db.Table("users").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.company_id = ? AND roles.name = ? AND users.is_active = ?", companyID, "EMPLOYEE", true).
		Count(&count)

	return count < int64(maxEmployees), nil
}

func (s *PlanCacheService) Del(ctx context.Context, key string) error {
	return s.redis.Del(ctx, key)
}

func (s *PlanCacheService) Invalidate(ctx context.Context, companyID uint) {
	key := fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID)
	if err := s.redis.Del(ctx, key); err != nil {
		logger.Errorw("PlanCacheService: failed to delete features cache", "key", key, "err", err)
	}
	profileKey := fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, companyID)
	if err := s.redis.Del(ctx, profileKey); err != nil {
		logger.Errorw("PlanCacheService: failed to delete company profile cache", "key", profileKey, "err", err)
	}
}

func moduleInFeatures(modules []string, target string) bool {
	for _, m := range modules {
		if m == target {
			return true
		}
	}
	return false
}
