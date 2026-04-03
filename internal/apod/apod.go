package apod

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const apiURL = "https://api.nasa.gov/planetary/apod"

type APOD struct {
	Title     string    `json:"title"`
	Date      string    `json:"date"`
	MediaType string    `json:"media_type"`
	FetchedAt time.Time `json:"fetched_at"`
}

func (a *APOD) FormatStatus() string {
	if a == nil || a.Title == "" {
		return ""
	}
	title := a.Title
	if len(title) > 40 {
		title = title[:39] + "…"
	}
	return fmt.Sprintf("🔭 APOD: \"%s\"", title)
}

func Fetch(ctx context.Context, apiKey string, useCache bool) (*APOD, error) {
	if useCache {
		cached, err := loadCache()
		if err == nil && cached != nil {
			return cached, nil
		}
	}

	result, err := fetchFromAPI(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	if useCache {
		_ = saveCache(result)
	}

	return result, nil
}

func fetchFromAPI(ctx context.Context, apiKey string) (*APOD, error) {
	url := fmt.Sprintf("%s?api_key=%s", apiURL, apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "stellar-status/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching APOD: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("APOD API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return parseResponse(body)
}

func parseResponse(data []byte) (*APOD, error) {
	var result APOD
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	result.FetchedAt = time.Now()
	return &result, nil
}
