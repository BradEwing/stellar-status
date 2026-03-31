package launches

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const cacheTTL = 10 * time.Minute

// cachedData wraps launches with a timestamp for TTL checks.
type cachedData struct {
	FetchedAt time.Time `json:"fetched_at"`
	Launches  []Launch  `json:"launches"`
}

// Cache provides file-based caching of launch data.
type Cache struct {
	path string
}

// NewCache creates a cache that stores data in ~/.cache/stellar-status/.
// The cache file is named per site to avoid cross-site staleness.
func NewCache(siteAbbrev string) (*Cache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".cache", "stellar-status")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	filename := fmt.Sprintf("launches-%s.json", strings.ToLower(siteAbbrev))
	return &Cache{path: filepath.Join(dir, filename)}, nil
}

// NewCacheWithPath creates a cache at a specific path (useful for testing).
func NewCacheWithPath(path string) *Cache {
	return &Cache{path: path}
}

// Get returns cached launches if the cache is fresh, or nil if stale/missing.
func (c *Cache) Get() []Launch {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return nil
	}

	var cached cachedData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil
	}

	if time.Since(cached.FetchedAt) > cacheTTL {
		return nil
	}

	return cached.Launches
}

// Set writes launches to the cache file.
func (c *Cache) Set(launches []Launch) error {
	cached := cachedData{
		FetchedAt: time.Now(),
		Launches:  launches,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0o644)
}
