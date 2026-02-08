package leave

import (
	"context"
	"hris-backend/internal/infrastructure"
	"hris-backend/pkg/logger"
)

type LeaveScheduler interface {
	Start()
}

type leaveScheduler struct {
	cronProvider *infrastructure.CronProvider
	leaveService Service
}

func NewLeaveScheduler(cronProvider *infrastructure.CronProvider, leaveService Service) LeaveScheduler {
	return &leaveScheduler{cronProvider, leaveService}
}

func (sch *leaveScheduler) Start() {
	logger.Info("Leave Scheduler Started...")

	_, err := sch.cronProvider.Cron.AddFunc("0 0 1 1 *", func() {
		logger.Info("[SCHEDULER] Starting Annual Leave Balance Generation...")

		if err := sch.leaveService.GenerateAnnualBalance(context.Background()); err != nil {
			logger.Errorf("[SCHEDULER] Failed: %v\n", err)
		} else {
			logger.Info("[SCHEDULER] Success! Annual balances generated.")
		}
	})

	if err != nil {
		logger.Errorf("Failed to start scheduler ", err)
	}

	sch.cronProvider.Cron.Start()
	defer sch.cronProvider.Cron.Stop()
}
