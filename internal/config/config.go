package config

import (
	"time"
)

func MustGetHealthCheckTimeout() time.Duration {
	return mustGetEnvDuration("HEALTH_CHECK_TIMEOUT")
}

func MustGetHealthCheckInterval() time.Duration {
	return mustGetEnvDuration("HEALTH_CHECK_INTERVAL")
}

func MustGetHealthCheckRetries() int {
	return mustGetEnvInt("HEALTH_CHECK_RETRIES")
}

func MustGetDiscordWebhookTimeout() time.Duration {
	return mustGetEnvDuration("DISCORD_WEBHOOK_TIMEOUT")
}

func MustGetDiscordWebhookRetries() int {
	return mustGetEnvInt("DISCORD_WEBHOOK_RETRIES")
}

func MustGetDiscordWebhookURL() string {
	return mustGetSecretString("discord_webhook_url")
}

func MustGetCheckURLs() []string {
	return mustGetCheckURLs(mustGetEnvString("URLS_PATH"))
}
