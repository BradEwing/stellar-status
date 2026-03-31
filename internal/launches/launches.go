package launches

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// NextLaunchResult holds the formatted next launch info.
type NextLaunchResult struct {
	Launch    Launch
	Countdown string // e.g. "3d 4h"
}

// FormatStatus returns a compact status string like "🚀[VBG] Falcon 9 in 3d 4h [Go]".
func (r NextLaunchResult) FormatStatus() string {
	name := r.Launch.Name
	if parts := strings.SplitN(name, "|", 2); len(parts) > 1 {
		name = strings.TrimSpace(parts[0])
	}
	return fmt.Sprintf("🚀[VBG] %s in %s [%s]", name, r.Countdown, r.Launch.Status.Abbrev)
}

// Tracker fetches and caches VBG launch data.
type Tracker struct {
	client *Client
	cache  *Cache
}

// NewTracker creates a tracker with default client and cache.
func NewTracker() (*Tracker, error) {
	cache, err := NewCache()
	if err != nil {
		return nil, fmt.Errorf("initializing cache: %w", err)
	}
	return &Tracker{
		client: NewClient(),
		cache:  cache,
	}, nil
}

// NewTrackerWithDeps creates a tracker with injected dependencies.
func NewTrackerWithDeps(client *Client, cache *Cache) *Tracker {
	return &Tracker{client: client, cache: cache}
}

// NextLaunch returns the next upcoming VBG launch, using cache when available.
func (t *Tracker) NextLaunch(ctx context.Context) (*NextLaunchResult, error) {
	launches := t.cache.Get()
	if launches == nil {
		var err error
		launches, err = t.client.FetchUpcoming(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching launches: %w", err)
		}
		_ = t.cache.Set(launches)
	}

	now := time.Now()
	for _, l := range launches {
		if l.NET.After(now) {
			return &NextLaunchResult{
				Launch:    l,
				Countdown: formatCountdown(l.NET, now),
			}, nil
		}
	}

	return nil, nil
}

// formatCountdown returns a human-readable duration like "3d 4h" or "2h 15m".
func formatCountdown(target, now time.Time) string {
	d := target.Sub(now)
	if d < 0 {
		return "past"
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}

	if len(parts) == 0 {
		return "<1m"
	}
	return strings.Join(parts, " ")
}
