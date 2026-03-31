package launches

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixtureResponse(t *testing.T, launches []apiLaunchEntry) *httptest.Server {
	t.Helper()
	resp := apiResponse{Count: len(launches), Results: launches}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func testCache(t *testing.T) *Cache {
	t.Helper()
	path := filepath.Join(t.TempDir(), "launches.json")
	return NewCacheWithPath(path)
}

func TestClient_FetchUpcoming(t *testing.T) {
	future := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
	server := fixtureResponse(t, []apiLaunchEntry{
		{
			ID:       "test-1",
			Name:     "Falcon 9 Block 5 | Starlink Group 12-1",
			NET:      future,
			Status:   Status{ID: 1, Name: "Go for Launch", Abbrev: "Go"},
			Location: "Vandenberg SFB, CA, USA",
			Pad:      "Space Launch Complex 4E",
		},
	})
	defer server.Close()

	client := &Client{httpClient: server.Client(), baseURL: server.URL}
	launches, err := client.FetchUpcoming(context.Background())

	require.NoError(t, err)
	require.Len(t, launches, 1)
	assert.Equal(t, "Falcon 9 Block 5 | Starlink Group 12-1", launches[0].Name)
	assert.Equal(t, "Go", launches[0].Status.Abbrev)
}

func TestClient_FetchUpcoming_EmptyResponse(t *testing.T) {
	server := fixtureResponse(t, []apiLaunchEntry{})
	defer server.Close()

	client := &Client{httpClient: server.Client(), baseURL: server.URL}
	launches, err := client.FetchUpcoming(context.Background())

	require.NoError(t, err)
	assert.Empty(t, launches)
}

func TestClient_FetchUpcoming_BadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := &Client{httpClient: server.Client(), baseURL: server.URL}
	_, err := client.FetchUpcoming(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "429")
}

func TestCache_GetSet(t *testing.T) {
	cache := testCache(t)

	// Empty cache returns nil.
	assert.Nil(t, cache.Get())

	launches := []Launch{{ID: "1", Name: "Test Launch", NET: time.Now().Add(time.Hour)}}
	require.NoError(t, cache.Set(launches))

	got := cache.Get()
	require.Len(t, got, 1)
	assert.Equal(t, "Test Launch", got[0].Name)
}

func TestCache_Expiry(t *testing.T) {
	cache := testCache(t)

	// Write a cache entry with an old timestamp.
	old := cachedData{
		FetchedAt: time.Now().Add(-20 * time.Minute),
		Launches:  []Launch{{ID: "1", Name: "Old Launch"}},
	}
	data, _ := json.Marshal(old)
	os.WriteFile(cache.path, data, 0o644)

	assert.Nil(t, cache.Get(), "expired cache should return nil")
}

func TestTracker_NextLaunch(t *testing.T) {
	future := time.Now().Add(72 * time.Hour)
	server := fixtureResponse(t, []apiLaunchEntry{
		{
			ID:     "launch-1",
			Name:   "Falcon 9 Block 5 | Starlink Group 17-17",
			NET:    future.Format(time.RFC3339),
			Status: Status{ID: 1, Name: "Go for Launch", Abbrev: "Go"},
		},
	})
	defer server.Close()

	client := &Client{httpClient: server.Client(), baseURL: server.URL}
	cache := testCache(t)
	tracker := NewTrackerWithDeps(client, cache)

	result, err := tracker.NextLaunch(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Falcon 9 Block 5 | Starlink Group 17-17", result.Launch.Name)
	assert.Contains(t, result.Countdown, "d")
}

func TestTracker_NextLaunch_UsesCache(t *testing.T) {
	cache := testCache(t)
	future := time.Now().Add(24 * time.Hour)
	cache.Set([]Launch{{ID: "cached", Name: "Cached Launch", NET: future, Status: Status{Abbrev: "TBC"}}})

	// Client points to a server that would fail — proving cache is used.
	client := &Client{httpClient: &http.Client{}, baseURL: "http://localhost:1"}
	tracker := NewTrackerWithDeps(client, cache)

	result, err := tracker.NextLaunch(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Cached Launch", result.Launch.Name)
}

func TestTracker_NextLaunch_SkipsPastLaunches(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	future := time.Now().Add(48 * time.Hour)
	server := fixtureResponse(t, []apiLaunchEntry{
		{ID: "old", Name: "Past Launch", NET: past.Format(time.RFC3339), Status: Status{Abbrev: "Go"}},
		{ID: "new", Name: "Future Launch", NET: future.Format(time.RFC3339), Status: Status{Abbrev: "Go"}},
	})
	defer server.Close()

	client := &Client{httpClient: server.Client(), baseURL: server.URL}
	tracker := NewTrackerWithDeps(client, testCache(t))

	result, err := tracker.NextLaunch(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Future Launch", result.Launch.Name)
}

func TestTracker_NextLaunch_NoUpcoming(t *testing.T) {
	server := fixtureResponse(t, []apiLaunchEntry{})
	defer server.Close()

	client := &Client{httpClient: server.Client(), baseURL: server.URL}
	tracker := NewTrackerWithDeps(client, testCache(t))

	result, err := tracker.NextLaunch(context.Background())

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestFormatCountdown(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		target   time.Time
		expected string
	}{
		{"3 days", now.Add(3*24*time.Hour + 4*time.Hour), "3d 4h"},
		{"3 days with minutes", now.Add(3*24*time.Hour + 4*time.Hour + 15*time.Minute), "3d 4h 15m"},
		{"5 hours", now.Add(5*time.Hour + 30*time.Minute), "5h 30m"},
		{"45 minutes", now.Add(45 * time.Minute), "45m"},
		{"just over a day", now.Add(25 * time.Hour), "1d 1h"},
		{"less than a minute", now.Add(30 * time.Second), "<1m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, formatCountdown(tt.target, now))
		})
	}
}

func TestFormatStatus(t *testing.T) {
	r := NextLaunchResult{
		Launch:    Launch{Name: "Falcon 9 Block 5 | Starlink Group 17-17", Status: Status{Abbrev: "Go"}},
		Countdown: "3d 4h",
	}
	assert.Equal(t, "🚀[VBG] Falcon 9 Block 5 in 3d 4h", r.FormatStatus())
}
