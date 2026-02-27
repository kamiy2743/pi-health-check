package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type failure struct {
	url string
	err error
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

func checkURLs(urls []string) []failure {
	client := &http.Client{Timeout: checkTimeout}
	failures := make([]failure, 0, len(urls))
	for _, url := range urls {
		fmt.Printf("チェック中: %s\n", url)
		if err := checkURL(client, url); err != nil {
			failures = append(failures, failure{url: url, err: err})
		}
	}
	return failures
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
