package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	confPath       = "/home/kamiy2743/workspace/health-check/.conf"
	statePath      = "/run/healthcheck/state.json"
	checkTimeout   = 10 * time.Second
	webhookTimeout = 10 * time.Second
)

func main() {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		fmt.Fprintln(os.Stderr, "DISCORD_WEBHOOK_URL が未設定です")
		os.Exit(1)
	}

	urls, err := readURLs(confPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "configを読み込めませんでした: %s: %v\n", confPath, err)
		os.Exit(1)
	}

	prevState, err := readState(statePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "状態の読込に失敗しました: %s: %v\n", statePath, err)
		os.Exit(1)
	}

	failures := checkURLs(urls)

	client := webhookClient{
		Client: &http.Client{Timeout: webhookTimeout},
		URL:    webhookURL,
	}

	if err := handleResult(prevState, client, failures); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func handleResult(prevState checkState, client webhookClient, failures []failure) error {
	if len(failures) == 0 {
		fmt.Println("すべてのヘルスチェックが成功しました")
	} else {
		fmt.Fprintln(os.Stderr, "失敗したURL:")
		for _, f := range failures {
			fmt.Fprintf(os.Stderr, "- %s (%s)\n", f.url, f.err)
		}
	}

	failedURLs := make([]string, 0, len(failures))
	for _, f := range failures {
		failedURLs = append(failedURLs, f.url)
	}

	added, resolved := diffFailedURLs(prevState.FailedURLs, failedURLs)
	if err := notifyDiff(client, failedURLs, added, resolved); err != nil {
		return fmt.Errorf("Discord通知に失敗しました: %v", err)
	}

	if err := writeState(statePath, checkState{
		FailedURLs: failedURLs,
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("状態ファイルの書込に失敗しました: %s: %v", statePath, err)
	}
	return nil
}

func diffFailedURLs(prev, current []string) (added []string, resolved []string) {
	prevSet := make(map[string]struct{}, len(prev))
	for _, u := range prev {
		prevSet[u] = struct{}{}
	}

	currentSet := make(map[string]struct{}, len(current))
	for _, u := range current {
		currentSet[u] = struct{}{}
		if _, ok := prevSet[u]; !ok {
			added = append(added, u)
		}
	}

	for _, u := range prev {
		if _, ok := currentSet[u]; !ok {
			resolved = append(resolved, u)
		}
	}

	return added, resolved
}
