package meteors

import (
	"testing"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
	"github.com/stretchr/testify/assert"
)

var loc = astro.Location{Latitude: 34.742, Longitude: -120.572}

func TestForTime_PerseidsPeak(t *testing.T) {
	peak := time.Date(2025, 8, 12, 22, 0, 0, 0, time.UTC)
	s := ForTime(peak, loc)
	assert.NotNil(t, s.Shower)
	assert.True(t, s.PeakTonight)
	assert.Contains(t, s.FormatStatus(), "peak tonight")
	assert.Contains(t, s.FormatStatus(), "Perseids")
}

func TestForTime_GeminidsTwoDaysBefore(t *testing.T) {
	before := time.Date(2025, 12, 12, 12, 0, 0, 0, time.UTC)
	s := ForTime(before, loc)
	assert.NotNil(t, s.Shower)
	assert.Contains(t, s.FormatStatus(), "in 2d")
	assert.Contains(t, s.FormatStatus(), "Geminids")
}

func TestForTime_ActiveButNotNearPeak(t *testing.T) {
	early := time.Date(2025, 10, 3, 12, 0, 0, 0, time.UTC)
	s := ForTime(early, loc)
	assert.NotNil(t, s.Shower)
	assert.True(t, s.Active)
	assert.Contains(t, s.FormatStatus(), "active")
	assert.Contains(t, s.FormatStatus(), "peak")
}

func TestForTime_NoActiveShowers(t *testing.T) {
	quiet := time.Date(2025, 3, 1, 12, 0, 0, 0, time.UTC)
	s := ForTime(quiet, loc)
	assert.Equal(t, "", s.FormatStatus())
}

func TestFormatStatus_PeakTonight(t *testing.T) {
	s := Status{
		Shower:      &showers[5],
		DaysToPeak:  0,
		Active:      true,
		PeakTonight: true,
	}
	result := s.FormatStatus()
	assert.Contains(t, result, "☄️")
	assert.Contains(t, result, "peak tonight")
	assert.Contains(t, result, "ZHR ~100")
}

func TestFormatStatus_PeakSoon(t *testing.T) {
	s := Status{
		Shower:     &showers[9],
		DaysToPeak: 2,
		Active:     true,
	}
	result := s.FormatStatus()
	assert.Contains(t, result, "peak in 2d")
	assert.Contains(t, result, "ZHR ~150")
}

func TestFormatStatus_ActiveNotImminent(t *testing.T) {
	s := Status{
		Shower:     &showers[7],
		DaysToPeak: 10,
		Active:     true,
	}
	result := s.FormatStatus()
	assert.Contains(t, result, "active")
	assert.Contains(t, result, "Orionids")
}

func TestFormatStatus_UpcomingNotActive(t *testing.T) {
	s := Status{
		Shower:     &showers[1],
		DaysToPeak: 5,
		Active:     false,
	}
	result := s.FormatStatus()
	assert.Contains(t, result, "in 5d")
	assert.Contains(t, result, "Lyrids")
}

func TestFormatStatus_NilShower(t *testing.T) {
	s := Status{}
	assert.Equal(t, "", s.FormatStatus())
}

func TestCurrent_DoesNotPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		Current(loc)
	})
}
