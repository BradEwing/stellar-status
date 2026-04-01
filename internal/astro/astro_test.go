package astro

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJulianDate_J2000Epoch(t *testing.T) {
	// J2000.0 = 2000-01-01T12:00:00Z = JD 2451545.0
	j2000 := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	jd := JulianDate(j2000)
	assert.InDelta(t, 2451545.0, jd, 0.0001)
}

func TestJulianDate_KnownDate(t *testing.T) {
	// 1999-01-01T00:00:00Z = JD 2451179.5
	dt := time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
	jd := JulianDate(dt)
	assert.InDelta(t, 2451179.5, jd, 0.0001)
}

func TestJulianCentury_J2000(t *testing.T) {
	T := JulianCentury(2451545.0)
	assert.InDelta(t, 0.0, T, 1e-10)
}

func TestNormalizeAngle(t *testing.T) {
	assert.InDelta(t, 0.0, NormalizeAngle(0), 1e-10)
	assert.InDelta(t, 0.0, NormalizeAngle(360), 1e-10)
	assert.InDelta(t, 270.0, NormalizeAngle(-90), 1e-10)
	assert.InDelta(t, 45.0, NormalizeAngle(405), 1e-10)
}

func TestMeanObliquity_J2000(t *testing.T) {
	obl := MeanObliquity(0)
	// Should be approximately 23.439 degrees at J2000.0
	assert.InDelta(t, 23.439, obl, 0.01)
}

func TestEclipticToEquatorial_VernalEquinox(t *testing.T) {
	// At ecliptic longitude 0, latitude 0: RA=0, Dec=0
	ra, dec := EclipticToEquatorial(0, 0, 23.439)
	assert.InDelta(t, 0.0, ra, 0.01)
	assert.InDelta(t, 0.0, dec, 0.01)
}

func TestEclipticToEquatorial_SummerSolstice(t *testing.T) {
	// At ecliptic longitude 90: RA=90-ish, Dec=obliquity
	obl := 23.439
	ra, dec := EclipticToEquatorial(90, 0, obl)
	assert.InDelta(t, 90.0, ra, 0.5)
	assert.InDelta(t, obl, dec, 0.5)
}

func TestAltitudeAzimuth_Zenith(t *testing.T) {
	// A star at the zenith: HA=0, Dec=Lat
	alt, _ := AltitudeAzimuth(0, 45, 45)
	assert.InDelta(t, 90.0, alt, 0.01)
}

func TestAltitudeAzimuth_Horizon(t *testing.T) {
	// Star on celestial equator at HA=90 (6h west), observer at equator
	alt, _ := AltitudeAzimuth(90, 0, 0)
	assert.InDelta(t, 0.0, alt, 0.5)
}

func TestLocalSiderealTime_Deterministic(t *testing.T) {
	jd := JulianDate(time.Date(2024, 6, 21, 12, 0, 0, 0, time.UTC))
	lst1 := LocalSiderealTime(jd, 0)
	lst2 := LocalSiderealTime(jd, 0)
	assert.Equal(t, lst1, lst2)

	// LST at east longitude should be greater
	lstEast := LocalSiderealTime(jd, 90)
	expected := NormalizeAngle(lst1 + 90)
	assert.InDelta(t, expected, lstEast, 0.01)
}
