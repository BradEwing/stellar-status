package launches

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL       = "https://ll.thespacedevs.com/2.2.0"
	vandenbergID  = 11
	defaultLimit  = 5
	requestTimeout = 10 * time.Second
)

// Launch represents a single upcoming launch.
type Launch struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`      // "Rocket | Mission"
	NET      time.Time `json:"net"`       // No Earlier Than
	Status   Status    `json:"status"`
	Location string    `json:"location"`
	Pad      string    `json:"pad"`
}

// Status represents launch readiness.
type Status struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Abbrev string `json:"abbrev"` // "Go", "TBC", "TBD"
}

// apiResponse is the LL2 list-mode response envelope.
type apiResponse struct {
	Count   int              `json:"count"`
	Results []apiLaunchEntry `json:"results"`
}

// apiLaunchEntry maps the LL2 list-mode launch object.
type apiLaunchEntry struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	NET      string    `json:"net"`
	Status   Status    `json:"status"`
	Location string    `json:"location"`
	Pad      string    `json:"pad"`
}

// Client fetches launch data from the Launch Library 2 API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new API client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: requestTimeout},
		baseURL:    baseURL,
	}
}

// FetchUpcoming returns upcoming launches from Vandenberg.
func (c *Client) FetchUpcoming(ctx context.Context) ([]Launch, error) {
	url := fmt.Sprintf("%s/launch/upcoming/?format=json&location__ids=%d&limit=%d&mode=list",
		c.baseURL, vandenbergID, defaultLimit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "stellar-status/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching launches: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from launch API", resp.StatusCode)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	launches := make([]Launch, 0, len(apiResp.Results))
	for _, entry := range apiResp.Results {
		t, err := time.Parse(time.RFC3339, entry.NET)
		if err != nil {
			continue
		}
		launches = append(launches, Launch{
			ID:       entry.ID,
			Name:     entry.Name,
			NET:      t,
			Status:   entry.Status,
			Location: entry.Location,
			Pad:      entry.Pad,
		})
	}

	return launches, nil
}
