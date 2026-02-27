package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type webhookClient struct {
	Client *http.Client
	URL    string
}

type webhookPayload struct {
	Embeds []embed `json:"embeds"`
}

type embed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       int    `json:"color,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

func notifyDiff(client webhookClient, currentFailedURLs, added, resolved []string) error {
	if len(added) == 0 && len(resolved) == 0 {
		if len(currentFailedURLs) > 0 {
			fmt.Fprintln(os.Stderr, "疎通不可が継続中のため通知をスキップします")
		}
		return nil
	}

	var title string
	var color int
	var description string

	switch {
	case len(currentFailedURLs) == 0 && len(resolved) > 0:
		title = "疎通不可がすべて解消しました"
		color = 0x2ECC71
		fmt.Fprintln(os.Stderr, "状態変化: 全回復 (通知送信)")
	case len(resolved) > 0:
		title = "疎通不可が一部解消しました"
		color = 0xF39C12
		description = buildPartialRecoveryDescription(currentFailedURLs, resolved)
		fmt.Fprintln(os.Stderr, "状態変化: 一部回復 (通知送信)")
	case len(added) > 0 && len(resolved) == 0 && len(currentFailedURLs) == len(added):
		title = "疎通不可が発生しました"
		color = 0xE74C3C
		description = buildCurrentFailuresDescription(currentFailedURLs)
		fmt.Fprintln(os.Stderr, "状態変化: 疎通不可発生 (通知送信)")
	default:
		title = "疎通不可が増加しました"
		color = 0xE74C3C
		description = buildCurrentFailuresDescription(currentFailedURLs)
		fmt.Fprintln(os.Stderr, "状態変化: 疎通不可増加 (通知送信)")
	}

	return sendDiscord(client, webhookPayload{
		Embeds: []embed{
			{
				Title:       title,
				Description: description,
				Color:       color,
				Timestamp:   time.Now().Format(time.RFC3339),
			},
		},
	})
}

func buildCurrentFailuresDescription(currentFailedURLs []string) string {
	if len(currentFailedURLs) == 0 {
		return ""
	}

	lines := []string{"疎通不可"}
	for _, u := range currentFailedURLs {
		lines = append(lines, fmt.Sprintf("⚠️ <%s>", u))
	}
	return strings.Join(lines, "\n")
}

func buildPartialRecoveryDescription(currentFailedURLs, resolved []string) string {
	lines := []string{}
	if len(currentFailedURLs) > 0 {
		lines = append(lines, "疎通不可")
		for _, u := range currentFailedURLs {
			lines = append(lines, fmt.Sprintf("⚠️ <%s>", u))
		}
	}
	if len(resolved) > 0 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, "回復")
		for _, u := range resolved {
			lines = append(lines, fmt.Sprintf("✅ <%s>", u))
		}
	}
	return strings.Join(lines, "\n")
}

func sendDiscord(client webhookClient, payload webhookPayload) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Client.Post(client.URL, "application/json", bytes.NewReader(payloadJSON))
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
