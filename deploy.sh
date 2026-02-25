#!/usr/bin/env bash

set -euo pipefail

SERVICE_SRC="/home/kamiy2743/workspace/health-check/healthcheck.service"
SERVICE_DST="/etc/systemd/system/healthcheck.service"
TIMER_SRC="/home/kamiy2743/workspace/health-check/healthcheck.timer"
TIMER_DST="/etc/systemd/system/healthcheck.timer"

/usr/local/go/bin/go build -o check-go check.go

sudo rsync -av --delete $SERVICE_SRC $SERVICE_DST
sudo chmod 644 $SERVICE_DST
sudo rsync -av --delete $TIMER_SRC $TIMER_DST
sudo chmod 644 $TIMER_DST

sudo systemctl daemon-reload
