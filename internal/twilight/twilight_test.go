package twilight

import (
	"testing"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
	"github.com/stretchr/testify/assert"
)

var vandenberg = astro.Location{Latitude: 34.742, Longitude: -120.572}

func TestForDate_SummerSolstice_Vandenberg(t *testing.T) {
	date := time.Date(2024, 6, 21, 12, 0, 0, 0, time.UTC)
	st := ForDate(date, vandenberg)

	assert.False(t, st.NeverRises)
	assert.False(t, st.NeverSets)
	assert.False(t, st.Sunrise.IsZero())
	assert.False(t, st.Sunset.IsZero())
	assert.True(t, st.Sunset.After(st.Sunrise))

	sunriseHour := st.Sunrise.UTC().Hour()
	assert.InDelta(t, 12, sunriseHour, 2) // ~12-13 UTC = ~5-6am PDT

	sunsetHour := st.Sunset.UTC().Hour()
	assert.InDelta(t, 3, float64(sunsetHour), 2) // ~3 UTC next day = ~8pm PDT
}

func TestForDate_WinterSolstice_Vandenberg(t *testing.T) {
	date := time.Date(2024, 12, 21, 12, 0, 0, 0, time.UTC)
	st := ForDate(date, vandenberg)

	assert.False(t, st.NeverRises)
	assert.False(t, st.NeverSets)
	assert.True(t, st.Sunset.After(st.Sunrise))
}

func TestForDate_Equinox_Equator(t *testing.T) {
	equator := astro.Location{Latitude: 0, Longitude: 0}
	date := time.Date(2024, 3, 20, 12, 0, 0, 0, time.UTC)
	st := ForDate(date, equator)

	assert.False(t, st.NeverRises)
	assert.False(t, st.NeverSets)

	dayLength := st.Sunset.Sub(st.Sunrise).Hours()
	assert.InDelta(t, 12.0, dayLength, 0.5)
}

func TestForDate_TwilightOrder(t *testing.T) {
	date := time.Date(2024, 6, 21, 12, 0, 0, 0, time.UTC)
	st := ForDate(date, vandenberg)

	if !st.AstroDawn.IsZero() && !st.CivilDawn.IsZero() {
		assert.True(t, st.AstroDawn.Before(st.CivilDawn))
	}
	if !st.CivilDawn.IsZero() {
		assert.True(t, st.CivilDawn.Before(st.Sunrise))
	}
	if !st.CivilDusk.IsZero() {
		assert.True(t, st.Sunset.Before(st.CivilDusk))
	}
	if !st.CivilDusk.IsZero() && !st.AstroDusk.IsZero() {
		assert.True(t, st.CivilDusk.Before(st.AstroDusk))
	}
}

func TestForDate_PolarNeverSets(t *testing.T) {
	arctic := astro.Location{Latitude: 80, Longitude: 0}
	date := time.Date(2024, 6, 21, 12, 0, 0, 0, time.UTC)
	st := ForDate(date, arctic)

	assert.True(t, st.NeverSets)
	assert.False(t, st.NeverRises)
}

func TestForDate_PolarNeverRises(t *testing.T) {
	arctic := astro.Location{Latitude: 80, Longitude: 0}
	date := time.Date(2024, 12, 21, 12, 0, 0, 0, time.UTC)
	st := ForDate(date, arctic)

	assert.True(t, st.NeverRises)
	assert.False(t, st.NeverSets)
}

func TestFormatStatus_Day(t *testing.T) {
	st := SunTimes{
		Sunrise:         time.Date(2024, 6, 21, 12, 50, 0, 0, time.UTC),
		Sunset:          time.Date(2024, 6, 22, 3, 20, 0, 0, time.UTC),
		GoldenHourStart: time.Date(2024, 6, 22, 2, 30, 0, 0, time.UTC),
	}
	now := time.Date(2024, 6, 21, 18, 0, 0, 0, time.UTC)
	s := st.FormatStatus(now)
	assert.Contains(t, s, "🌅 sunset")
}

func TestFormatStatus_NeverSets(t *testing.T) {
	st := SunTimes{NeverSets: true}
	assert.Equal(t, "☀️ midnight sun", st.FormatStatus(time.Now()))
}

func TestFormatStatus_NeverRises(t *testing.T) {
	st := SunTimes{NeverRises: true}
	assert.Equal(t, "🌃 polar night", st.FormatStatus(time.Now()))
}

func TestFormatStatus_GoldenHour(t *testing.T) {
	sunset := time.Date(2024, 6, 22, 3, 20, 0, 0, time.UTC)
	st := SunTimes{
		Sunrise:         time.Date(2024, 6, 21, 12, 50, 0, 0, time.UTC),
		Sunset:          sunset,
		GoldenHourStart: time.Date(2024, 6, 22, 2, 30, 0, 0, time.UTC),
		CivilDusk:       time.Date(2024, 6, 22, 3, 50, 0, 0, time.UTC),
	}
	now := time.Date(2024, 6, 22, 2, 45, 0, 0, time.UTC)
	assert.Equal(t, "🌇 golden hour", st.FormatStatus(now))
}

func TestToday_DoesNotPanic(t *testing.T) {
	st := Today(vandenberg)
	assert.NotZero(t, st.Date)
}
