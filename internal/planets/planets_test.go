package planets

import (
	"testing"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
	"github.com/stretchr/testify/assert"
)

var vandenberg = astro.Location{Latitude: 34.742, Longitude: -120.572}

func TestForTime_AllPlanetsComputed(t *testing.T) {
	dt := time.Date(2024, 6, 21, 4, 0, 0, 0, time.UTC) // nighttime at VBG
	v := ForTime(dt, vandenberg)

	assert.Len(t, v.Planets, 5)
	names := make(map[PlanetName]bool)
	for _, p := range v.Planets {
		names[p.Name] = true
		assert.InDelta(t, 0, p.Altitude, 90.0) // reasonable range
	}
	assert.True(t, names[Mercury])
	assert.True(t, names[Venus])
	assert.True(t, names[Mars])
	assert.True(t, names[Jupiter])
	assert.True(t, names[Saturn])
}

func TestForTime_DaytimeNoneVisible(t *testing.T) {
	noon := time.Date(2024, 6, 21, 20, 0, 0, 0, time.UTC) // ~1pm PDT
	v := ForTime(noon, vandenberg)

	for _, p := range v.Planets {
		assert.False(t, p.Visible, "%s should not be visible during daytime", p.Name)
	}
}

func TestForTime_AltitudeInRange(t *testing.T) {
	dt := time.Date(2024, 3, 15, 4, 0, 0, 0, time.UTC)
	v := ForTime(dt, vandenberg)

	for _, p := range v.Planets {
		assert.GreaterOrEqual(t, p.Altitude, -90.0)
		assert.LessOrEqual(t, p.Altitude, 90.0)
	}
}

func TestFormatStatus_NoneVisible(t *testing.T) {
	v := Visibility{
		Planets: []PlanetStatus{
			{Name: Mercury, Visible: false},
			{Name: Venus, Visible: false},
		},
	}
	assert.Equal(t, "🔭 —", v.FormatStatus())
}

func TestFormatStatus_SomeVisible(t *testing.T) {
	v := Visibility{
		Planets: []PlanetStatus{
			{Name: Mercury, Visible: false},
			{Name: Venus, Visible: true},
			{Name: Mars, Visible: false},
			{Name: Jupiter, Visible: true},
			{Name: Saturn, Visible: false},
		},
	}
	assert.Equal(t, "🔭 Venus Jupiter", v.FormatStatus())
}

func TestCurrent_DoesNotPanic(t *testing.T) {
	v := Current(vandenberg)
	assert.NotZero(t, v.Timestamp)
	assert.Len(t, v.Planets, 5)
}

func TestHeliocentricPosition_EarthReasonable(t *testing.T) {
	T := 0.0 // J2000.0
	_, _, r := heliocentricPosition(Earth, T)
	assert.InDelta(t, 1.0, r, 0.02) // Earth ~1 AU from sun
}

func TestHeliocentricPosition_JupiterReasonable(t *testing.T) {
	T := 0.0
	_, _, r := heliocentricPosition(Jupiter, T)
	assert.InDelta(t, 5.2, r, 0.3) // Jupiter ~5.2 AU
}
