package aurora

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatStatus_AuroraLikely(t *testing.T) {
	s := &Status{Kp: 8.0}
	assert.Equal(t, "\U0001f30c Aurora likely! (Kp=8)", s.FormatStatus())
}

func TestFormatStatus_AuroraPossible(t *testing.T) {
	s := &Status{Kp: 6.0}
	assert.Equal(t, "\U0001f30c Aurora possible (Kp=6)", s.FormatStatus())
}

func TestFormatStatus_GeomagneticActive(t *testing.T) {
	s := &Status{Kp: 4.0}
	assert.Equal(t, "\U0001f30c Geomagnetic active (Kp=4)", s.FormatStatus())
}

func TestFormatStatus_SolarQuiet(t *testing.T) {
	s := &Status{Kp: 2.0}
	assert.Equal(t, "\u2600\ufe0f Solar quiet (Kp=2)", s.FormatStatus())
}

func TestFormatStatus_BoundaryKp7(t *testing.T) {
	s := &Status{Kp: 7.0}
	assert.Equal(t, "\U0001f30c Aurora likely! (Kp=7)", s.FormatStatus())
}

func TestFormatStatus_BoundaryKp5(t *testing.T) {
	s := &Status{Kp: 5.0}
	assert.Equal(t, "\U0001f30c Aurora possible (Kp=5)", s.FormatStatus())
}

func TestParseResponse_ValidData(t *testing.T) {
	now := time.Now().UTC()
	ts := now.Add(-1 * time.Hour).Format("2006-01-02 15:04:05")

	rows := [][]string{
		{"time_tag", "kp", "observed", "noaa_scale"},
		{ts, "3.33", "observed", "1"},
	}

	data, err := json.Marshal(rows)
	assert.NoError(t, err)

	status, err := parseResponse(data)
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.InDelta(t, 3.33, status.Kp, 0.001)
	assert.True(t, status.Observed)
}

func TestParseResponse_PicksLatestObserved(t *testing.T) {
	now := time.Now().UTC()
	older := now.Add(-3 * time.Hour).Format("2006-01-02 15:04:05")
	newer := now.Add(-1 * time.Hour).Format("2006-01-02 15:04:05")

	rows := [][]string{
		{"time_tag", "kp", "observed", "noaa_scale"},
		{older, "2.00", "observed", "0"},
		{newer, "5.00", "observed", "2"},
	}

	data, err := json.Marshal(rows)
	assert.NoError(t, err)

	status, err := parseResponse(data)
	assert.NoError(t, err)
	assert.InDelta(t, 5.0, status.Kp, 0.001)
}

func TestParseResponse_SkipsFutureForecasts(t *testing.T) {
	now := time.Now().UTC()
	past := now.Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
	future := now.Add(3 * time.Hour).Format("2006-01-02 15:04:05")

	rows := [][]string{
		{"time_tag", "kp", "observed", "noaa_scale"},
		{past, "2.00", "observed", "0"},
		{future, "8.00", "estimated", "4"},
	}

	data, err := json.Marshal(rows)
	assert.NoError(t, err)

	status, err := parseResponse(data)
	assert.NoError(t, err)
	assert.InDelta(t, 2.0, status.Kp, 0.001)
}

func TestParseResponse_EmptyData(t *testing.T) {
	rows := [][]string{
		{"time_tag", "kp", "observed", "noaa_scale"},
	}

	data, err := json.Marshal(rows)
	assert.NoError(t, err)

	_, err = parseResponse(data)
	assert.Error(t, err)
}

func TestParseResponse_InvalidJSON(t *testing.T) {
	_, err := parseResponse([]byte("not json"))
	assert.Error(t, err)
}
