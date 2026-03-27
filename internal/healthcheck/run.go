package healthcheck

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"health-check/internal/config"
)

var jst = time.FixedZone("JST", 9*60*60)

const logTimeFormat = "2006-01-02 15:04"

func RunPeriodic(interval time.Duration, now func() time.Time) {
	var previousFailed []string
	for {
		next := now().Truncate(interval).Add(interval)
		wait := time.Until(next)
		if wait > 0 {
			time.Sleep(wait)
		}

		failedURLs, err := run(previousFailed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		previousFailed = failedURLs
	}
}

func run(previousFailed []string) ([]string, error) {
	fmt.Printf("[%s] ヘルスチェックを開始します\n", time.Now().In(jst).Format(logTimeFormat))

	urls := config.MustGetCheckURLs()
	healthCheckTimeout := config.MustGetHealthCheckTimeout()
	healthCheckRetries := config.MustGetHealthCheckRetries()
	failures := checkURLs(urls, healthCheckTimeout, healthCheckRetries)

	webhookTimeout := config.MustGetDiscordWebhookTimeout()
	client := webhookClient{
		Client: &http.Client{Timeout: webhookTimeout},
		URL:    config.MustGetDiscordWebhookURL(),
	}

	failedURLs, err := handleResult(previousFailed, client, failures)
	if err != nil {
		return nil, err
	}

	return failedURLs, nil
}

func handleResult(previousFailed []string, client webhookClient, failures []failure) ([]string, error) {
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

	added, resolved := diffFailedURLs(previousFailed, failedURLs)
	if err := notifyDiff(client, failedURLs, added, resolved); err != nil {
		return nil, fmt.Errorf("Discord通知に失敗しました: %v", err)
	}

	return failedURLs, nil
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
