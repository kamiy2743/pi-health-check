package main

import (
	"time"

	"health-check/internal/config"
	"health-check/internal/healthcheck"
)

func main() {
	interval := config.MustGetHealthCheckInterval()
	healthcheck.RunPeriodic(interval, time.Now)
}
