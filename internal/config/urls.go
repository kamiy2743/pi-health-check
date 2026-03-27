package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func mustGetCheckURLs(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("URLファイルを開けませんでした: %v", err)
	}
	defer f.Close()

	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("URLファイルの読み取り中にエラーが発生しました: %v", err)
	}

	return urls
}
