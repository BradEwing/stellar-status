package astro

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSunEquatorial_SummerSolstice(t *testing.T) {
	// Around June 21, the sun's declination should be near +23.4 degrees
	solstice := time.Date(2024, 6, 21, 12, 0, 0, 0, time.UTC)
	_, dec := SunEquatorial(solstice)
	assert.InDelta(t, 23.4, dec, 1.0)
}

func TestSunEquatorial_WinterSolstice(t *testing.T) {
	// Around Dec 21, the sun's declination should be near -23.4 degrees
	solstice := time.Date(2024, 12, 21, 12, 0, 0, 0, time.UTC)
	_, dec := SunEquatorial(solstice)
	assert.InDelta(t, -23.4, dec, 1.0)
}

func TestSunEquatorial_Equinox(t *testing.T) {
	// Around March 20, the sun's declination should be near 0
	equinox := time.Date(2024, 3, 20, 12, 0, 0, 0, time.UTC)
	_, dec := SunEquatorial(equinox)
	assert.InDelta(t, 0.0, dec, 1.5)
}

func TestSunPosition_NoonAtEquator(t *testing.T) {
	// At solar noon on the equinox at the equator, altitude should be near 90
	// Use March 20 at ~12:00 UTC, longitude 0 (Greenwich)
	equinox := time.Date(2024, 3, 20, 12, 0, 0, 0, time.UTC)
	loc := Location{Latitude: 0, Longitude: 0}
	alt, _ := SunPosition(equinox, loc)
	assert.InDelta(t, 90.0, alt, 5.0) // within 5 degrees is fine
}

func TestSunPosition_MidnightBelowHorizon(t *testing.T) {
	// At midnight at moderate latitudes, sun should be below horizon
	midnight := time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)
	loc := Location{Latitude: 45, Longitude: 0}
	alt, _ := SunPosition(midnight, loc)
	assert.Less(t, alt, 0.0)
}

func TestSunPosition_Vandenberg_Daytime(t *testing.T) {
	// Vandenberg at local noon (~20:00 UTC in winter) should have positive altitude
	noon := time.Date(2024, 1, 15, 20, 0, 0, 0, time.UTC)
	loc := Location{Latitude: 34.742, Longitude: -120.572}
	alt, _ := SunPosition(noon, loc)
	assert.Greater(t, alt, 10.0)
}
