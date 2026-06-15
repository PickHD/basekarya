package subscription

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"context"
	"fmt"
)

type Scheduler interface {
	Start()
	Stop()
}

type subscriptionScheduler struct {
	cronProvider *infrastructure.CronProvider
	repo         Repository
	cache        CacheProvider
}

func NewScheduler(cronProvider *infrastructure.CronProvider, repo Repository, cache CacheProvider) Scheduler {
	return &subscriptionScheduler{cronProvider, repo, cache}
}

func (sch *subscriptionScheduler) Start() {
	logger.Info("Subscription Expiry Scheduler Started...")

	_, err := sch.cronProvider.GetCron().AddFunc("0 0 * * *", func() {
		logger.Info("[SCHEDULER] Starting subscription expiry check...")

		ctx := context.Background()

		ids, err := sch.repo.FindExpiredCompanies(ctx)
		if err != nil {
			logger.Errorf("[SCHEDULER] Failed to find expired companies: %v", err)
			return
		}

		for _, companyID := range ids {
			if err := sch.repo.UpdateCompanyStatus(ctx, companyID, constants.SubStatusExpired); err != nil {
				logger.Errorf("[SCHEDULER] Failed to update status for company %d: %v", companyID, err)
				continue
			}

			_ = sch.cache.Del(ctx, fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID))
			_ = sch.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, companyID))
		}

		if len(ids) > 0 {
			logger.Infof("[SCHEDULER] Expired %d companies", len(ids))
		} else {
			logger.Info("[SCHEDULER] No expired companies found")
		}
	})

	if err != nil {
		logger.Errorf("Failed to start subscription expiry scheduler: %v", err)
	}

	sch.cronProvider.GetCron().Start()
}

func (sch *subscriptionScheduler) Stop() {
	if sch.cronProvider != nil && sch.cronProvider.GetCron() != nil {
		sch.cronProvider.GetCron().Stop()
		logger.Info("Subscription Expiry Scheduler Stopped.")
	}
}
