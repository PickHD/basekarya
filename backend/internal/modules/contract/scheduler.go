package contract

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/logger"
	"context"
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
	logger.Info("Contract Scheduler Started...")

	// Run every day at 08:00
	_, err := sch.cronProvider.GetCron().AddFunc("0 8 * * *", func() {
		logger.Info("[SCHEDULER] Checking for expiring contracts...")

		if err := sch.service.CheckExpiringContracts(context.Background()); err != nil {
			logger.Errorf("[SCHEDULER] Failed: %v\n", err)
		}
	})

	if err != nil {
		logger.Errorf("Failed to start contract scheduler ", err)
	}

	sch.cronProvider.GetCron().Start()
}

func (sch *scheduler) Stop() {
	if sch.cronProvider != nil && sch.cronProvider.GetCron() != nil {
		sch.cronProvider.GetCron().Stop()
		logger.Info("Contract Scheduler Stopped.")
	}
}
