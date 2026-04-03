package apod

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 6 * time.Hour

type cachedData struct {
	FetchedAt time.Time `json:"fetched_at"`
	APOD      APOD      `json:"apod"`
}

func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".cache", "stellar-status")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "apod.json"), nil
}

func loadCache() (*APOD, error) {
	path, err := cachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cached cachedData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	if time.Since(cached.FetchedAt) > cacheTTL {
		return nil, nil
	}

	return &cached.APOD, nil
}

func saveCache(a *APOD) error {
	path, err := cachePath()
	if err != nil {
		return err
	}

	cached := cachedData{
		FetchedAt: time.Now(),
		APOD:      *a,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
