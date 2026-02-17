package constants

import "time"

const (
	WriteWait  time.Duration = 10 * time.Second
	PongWait   time.Duration = 60 * time.Second
	PingPeriod time.Duration = (PongWait * 9) / 10
)
