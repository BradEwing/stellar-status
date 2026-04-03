package aurora

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const swpcURL = "https://services.swpc.noaa.gov/products/noaa-planetary-k-index-forecast.json"

type Status struct {
	Kp        float64
	Timestamp time.Time
	Observed  bool
}

func (s *Status) FormatStatus() string {
	kpInt := int(s.Kp)
	switch {
	case s.Kp >= 7:
		return fmt.Sprintf("\U0001f30c Aurora likely! (Kp=%d)", kpInt)
	case s.Kp >= 5:
		return fmt.Sprintf("\U0001f30c Aurora possible (Kp=%d)", kpInt)
	case s.Kp >= 4:
		return fmt.Sprintf("\U0001f30c Geomagnetic active (Kp=%d)", kpInt)
	default:
		return fmt.Sprintf("\u2600\ufe0f Solar quiet (Kp=%d)", kpInt)
	}
}

func Fetch(ctx context.Context, useCache bool) (*Status, error) {
	if useCache {
		cached, err := loadCache()
		if err == nil && cached != nil {
			return cached, nil
		}
	}

	status, err := fetchFromAPI(ctx)
	if err != nil {
		return nil, err
	}

	if useCache {
		_ = saveCache(status)
	}

	return status, nil
}

func fetchFromAPI(ctx context.Context) (*Status, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, swpcURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching Kp index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SWPC API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return parseResponse(body)
}

func parseResponse(data []byte) (*Status, error) {
	var rows [][]string
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("no Kp data in response")
	}

	now := time.Now().UTC()
	var best *Status

	for _, row := range rows[1:] {
		if len(row) < 3 {
			continue
		}

		t, err := time.Parse("2006-01-02 15:04:05", row[0])
		if err != nil {
			continue
		}

		kp, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			continue
		}

		observed := row[2] == "observed"

		if t.After(now) && !observed {
			continue
		}

		if best == nil || t.After(best.Timestamp) {
			best = &Status{
				Kp:        kp,
				Timestamp: t,
				Observed:  observed,
			}
		}
	}

	if best == nil {
		return nil, fmt.Errorf("no current Kp observation found")
	}

	return best, nil
}
