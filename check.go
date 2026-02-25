package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	confPath       = "/home/kamiy2743/workspace/health-check/.conf"
	checkTimeout   = 3 * time.Second
	webhookTimeout = 10 * time.Second
)

type failure struct {
	url string
	err error
}

type webhookPayload struct {
	Embeds []embed `json:"embeds"`
}

type embed struct {
	Title     string       `json:"title,omitempty"`
	Color     int          `json:"color,omitempty"`
	Fields    []embedField `json:"fields,omitempty"`
	Timestamp string       `json:"timestamp,omitempty"`
}

type embedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

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

	checkClient := &http.Client{Timeout: checkTimeout}
	webhookClient := &http.Client{Timeout: webhookTimeout}

	failures := make([]failure, 0, len(urls))
	for _, url := range urls {
		fmt.Printf("チェック中: %s\n", url)
		if err := checkURL(checkClient, url); err != nil {
			failures = append(failures, failure{url: url, err: err})
		}
	}

	if len(failures) > 0 {
		reportFailures(webhookClient, webhookURL, failures)
		return
	}

	fmt.Println("all checks passed")
}

func readURLs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	urls := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 空行やコメントはスキップ
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func checkURL(client *http.Client, url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}
	return nil
}

func reportFailures(webhookClient *http.Client, webhookURL string, failures []failure) {
	// 標準エラー出力
	fmt.Fprintln(os.Stderr, "失敗したURL:")
	for _, f := range failures {
		fmt.Fprintf(os.Stderr, "- %s (%s)\n", f.url, f.err)
	}

	// Discord通知
	fields := make([]embedField, 0, len(failures)+1)
	for _, f := range failures {
		fields = append(fields, embedField{
			Name:  fmt.Sprintf("URL : **<%s>**", f.url),
			Value: fmt.Sprintf("エラー : %s", f.err.Error()),
		})
	}
	payload := webhookPayload{
		Embeds: []embed{
			{
				Title:     "ヘルスチェックに失敗しました！ :warning:",
				Color:     0xE74C3C,
				Fields:    fields,
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
	}
	if err := sendDiscord(webhookClient, webhookURL, payload); err != nil {
		fmt.Fprintf(os.Stderr, "Discord通知に失敗しました: %v\n", err)
	}
}

func sendDiscord(client *http.Client, webhookURL string, payload webhookPayload) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(payloadJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status: %s", resp.Status)
	}
	return nil
}
