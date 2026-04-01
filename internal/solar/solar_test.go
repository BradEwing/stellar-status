package solar

import (
	"testing"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
	"github.com/stretchr/testify/assert"
)

func TestForTime_NoonEquatorEquinox(t *testing.T) {
	equinox := time.Date(2024, 3, 20, 12, 0, 0, 0, time.UTC)
	loc := astro.Location{Latitude: 0, Longitude: 0}
	p := ForTime(equinox, loc)

	assert.InDelta(t, 90.0, p.Altitude, 5.0)
	assert.True(t, p.IsDay)
}

func TestForTime_MidnightBelowHorizon(t *testing.T) {
	midnight := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)
	loc := astro.Location{Latitude: 45, Longitude: 0}
	p := ForTime(midnight, loc)

	assert.Less(t, p.Altitude, 0.0)
	assert.False(t, p.IsDay)
}

func TestForTime_NoonAzimuthSouth(t *testing.T) {
	noon := time.Date(2024, 6, 21, 12, 0, 0, 0, time.UTC)
	loc := astro.Location{Latitude: 45, Longitude: 0}
	p := ForTime(noon, loc)

	assert.InDelta(t, 180.0, p.Azimuth, 20.0)
}

func TestFormatStatus_Day(t *testing.T) {
	p := Position{Altitude: 45.3, IsDay: true}
	assert.Equal(t, "☀️ 45°", p.FormatStatus())
}

func TestFormatStatus_Twilight(t *testing.T) {
	p := Position{Altitude: 2.0, IsDay: true}
	assert.Equal(t, "🌅 2°", p.FormatStatus())
}

func TestFormatStatus_Night(t *testing.T) {
	p := Position{Altitude: -12.0, IsDay: false}
	assert.Equal(t, "🌙 -12°", p.FormatStatus())
}

func TestCurrent_DoesNotPanic(t *testing.T) {
	loc := astro.Location{Latitude: 34.742, Longitude: -120.572}
	p := Current(loc)
	assert.NotZero(t, p.Timestamp)
}
