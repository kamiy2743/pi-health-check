package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const dockerSecretsDir = "/run/secrets"

func mustGetSecretString(secretName string) string {
	path := filepath.Join(dockerSecretsDir, secretName)
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("secret %s の読み込みに失敗しました: %v", path, err)
	}

	value := strings.TrimSpace(string(raw))
	if value == "" {
		log.Fatalf("secret %s の内容が空です。", path)
	}

	return value
}
