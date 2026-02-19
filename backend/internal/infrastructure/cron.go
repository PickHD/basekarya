package infrastructure

import "github.com/robfig/cron/v3"

type CronProvider struct {
	cron *cron.Cron
}

func NewCronProvider() *CronProvider {
	return &CronProvider{
		cron: cron.New(),
	}
}

func (c *CronProvider) GetCron() *cron.Cron {
	return c.cron
}
