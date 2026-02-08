package infrastructure

import "github.com/robfig/cron/v3"

type CronProvider struct {
	Cron *cron.Cron
}

func NewCronProvider() *CronProvider {
	return &CronProvider{
		Cron: cron.New(),
	}
}
