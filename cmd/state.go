package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type checkState struct {
	FailedURLs []string `json:"failed_urls"`
	UpdatedAt  string   `json:"updated_at"`
}

func readState(path string) (checkState, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return checkState{FailedURLs: []string{}}, nil
		}
		return checkState{}, err
	}

	var state checkState
	if err := json.Unmarshal(b, &state); err != nil {
		return checkState{}, err
	}
	if state.FailedURLs == nil {
		state.FailedURLs = []string{}
	}
	return state, nil
}

func writeState(path string, state checkState) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	b, err := json.Marshal(state)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, "state-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(b); err != nil {
		tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
