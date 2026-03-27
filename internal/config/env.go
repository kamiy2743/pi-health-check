package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

func mustGetEnvString(envName string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Fatalf(".env に %s が未設定です。", envName)
	}
	return value
}

func mustGetEnvInt(envName string) int {
	env := mustGetEnvString(envName)
	value, err := strconv.Atoi(env)
	if err != nil {
		log.Fatalf("%s の値 %s は整数に変換できません。", envName, env)
	}
	return value
}

func mustGetEnvDuration(envName string) time.Duration {
	env := mustGetEnvString(envName)
	duration, err := time.ParseDuration(env)
	if err != nil {
		log.Fatalf("%s の値 %s は時間の形式として解釈できません。例: 5s, 10m, 1h", envName, env)
	}
	return duration
}
