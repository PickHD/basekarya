package notification

import (
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
	logger.Info("Notification Scheduler Started...")

	_, err := sch.cronProvider.Cron.AddFunc("0 0 3 * *", func() {
		logger.Info("[SCHEDULER] Starting Remove Old Notification...")

		if err := sch.service.DeleteReadOlderThan(3); err != nil {
			logger.Errorf("[SCHEDULER] Failed: %v\n", err)
		} else {
			logger.Info("[SCHEDULER] Success! notification old removed.")
		}
	})

	if err != nil {
		logger.Errorf("Failed to start scheduler ", err)
	}

	sch.cronProvider.Cron.Start()
}

func (sch *scheduler) Stop() {
	if sch.cronProvider != nil && sch.cronProvider.Cron != nil {
		sch.cronProvider.Cron.Stop()
		logger.Info("Notification Scheduler Stopped.")
	}
}
