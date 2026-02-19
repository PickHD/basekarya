package leave

import (
	"context"
	"hris-backend/internal/infrastructure"
	"hris-backend/pkg/logger"
)

type Scheduler interface {
	Start()
	Stop()
}

type scheduler struct {
	cronProvider *infrastructure.CronProvider
	service      Service
}

func NewScheduler(cronProvider *infrastructure.CronProvider, service Service) Scheduler {
	return &scheduler{cronProvider, service}
}

func (sch *scheduler) Start() {
	logger.Info("Leave Scheduler Started...")

	_, err := sch.cronProvider.GetCron().AddFunc("0 0 1 1 *", func() {
		logger.Info("[SCHEDULER] Starting Annual Leave Balance Generation...")

		if err := sch.service.GenerateAnnualBalance(context.Background()); err != nil {
			logger.Errorf("[SCHEDULER] Failed: %v\n", err)
		} else {
			logger.Info("[SCHEDULER] Success! Annual balances generated.")
		}
	})

	if err != nil {
		logger.Errorf("Failed to start scheduler ", err)
	}

	sch.cronProvider.GetCron().Start()
}

func (sch *scheduler) Stop() {
	if sch.cronProvider != nil && sch.cronProvider.GetCron() != nil {
		sch.cronProvider.GetCron().Stop()
		logger.Info("Leave Scheduler Stopped.")
	}
}
