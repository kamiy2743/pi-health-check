package healthcheck

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type failure struct {
	url string
	err error
}

func checkURLs(urls []string, timeout time.Duration, retries int) []failure {
	client := &http.Client{Timeout: timeout}
	failures := make([]failure, 0, len(urls))
	for _, url := range urls {
		healthURL := url + "/health"
		fmt.Printf("チェック中: %s\n", healthURL)
		if err := checkURL(client, healthURL, retries); err != nil {
			failures = append(failures, failure{url: url, err: err})
		}
	}
	return failures
}

func checkURL(client *http.Client, url string, retries int) error {
	for attempt := 0; attempt <= retries; attempt++ {
		resp, err := client.Get(url)
		if err != nil {
			if attempt == retries {
				return err
			}
			continue
		}

		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		if attempt == retries {
			return fmt.Errorf("status: %s", resp.Status)
		}
	}
	return nil
}
