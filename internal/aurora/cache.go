package aurora

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 30 * time.Minute

type cachedData struct {
	FetchedAt time.Time `json:"fetched_at"`
	Status    Status    `json:"status"`
}

func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", "stellar-status", "aurora.json"), nil
}

func loadCache() (*Status, error) {
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

	return &cached.Status, nil
}

func saveCache(s *Status) error {
	path, err := cachePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	cached := cachedData{
		FetchedAt: time.Now(),
		Status:    *s,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
