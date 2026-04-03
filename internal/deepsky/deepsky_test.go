package deepsky

import (
	"testing"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
	"github.com/stretchr/testify/assert"
)

var vandenberg = astro.Location{Latitude: 34.742, Longitude: -120.572}

func TestForTime_WinterNight_ObjectVisible(t *testing.T) {
	dt := time.Date(2024, 1, 15, 4, 0, 0, 0, time.UTC)
	obs := ForTime(dt, vandenberg)

	assert.True(t, obs.IsDark)
	assert.NotNil(t, obs.Best)
	assert.Greater(t, obs.Altitude, 15.0)
}

func TestForTime_SummerNight_M13Visible(t *testing.T) {
	dt := time.Date(2024, 7, 15, 6, 0, 0, 0, time.UTC)
	obs := ForTime(dt, vandenberg)

	assert.True(t, obs.IsDark)
	assert.NotNil(t, obs.Best)
	assert.Equal(t, "M13", obs.Best.ID)
	assert.Greater(t, obs.Altitude, 15.0)
}

func TestForTime_Daytime_NothingVisible(t *testing.T) {
	noon := time.Date(2024, 6, 21, 20, 0, 0, 0, time.UTC)
	obs := ForTime(noon, vandenberg)

	assert.False(t, obs.IsDark)
	assert.Nil(t, obs.Best)
}

func TestFormatStatus_WithVisibleObject(t *testing.T) {
	obs := Observation{
		Best:     &Object{ID: "M42", Name: "Orion Nebula"},
		Altitude: 62.3,
	}
	assert.Equal(t, "🌌 M42 (Orion Nebula) alt 62°", obs.FormatStatus())
}

func TestFormatStatus_NilBest_EmptyString(t *testing.T) {
	obs := Observation{Best: nil}
	assert.Equal(t, "", obs.FormatStatus())
}

func TestCurrent_DoesNotPanic(t *testing.T) {
	obs := Current(vandenberg)
	assert.NotZero(t, obs.Timestamp)
}
