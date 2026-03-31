package moon

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestForTime_KnownNewMoon(t *testing.T) {
	newMoon := time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)
	p := ForTime(newMoon)

	assert.Equal(t, NewMoon, p.Name)
	assert.Less(t, p.Illumination, 0.01)
	assert.Equal(t, "🌑", p.Emoji)
}

func TestForTime_KnownFullMoon(t *testing.T) {
	fullMoon := referenceNewMoon.Add(time.Duration(synodicMonth / 2 * 24 * float64(time.Hour)))
	p := ForTime(fullMoon)

	assert.Equal(t, FullMoon, p.Name)
	assert.Greater(t, p.Illumination, 0.99)
	assert.Equal(t, "🌕", p.Emoji)
}

func TestForTime_KnownFullMoon2024(t *testing.T) {
	// October 17, 2024 — Hunter's Moon.
	fullMoon := time.Date(2024, 10, 17, 11, 26, 0, 0, time.UTC)
	p := ForTime(fullMoon)

	assert.Equal(t, FullMoon, p.Name)
	assert.Greater(t, p.Illumination, 0.95)
}

func TestForTime_FirstQuarter(t *testing.T) {
	fq := referenceNewMoon.Add(time.Duration(synodicMonth / 4 * 24 * float64(time.Hour)))
	p := ForTime(fq)

	assert.Equal(t, FirstQuarter, p.Name)
	assert.InDelta(t, 0.5, p.Illumination, 0.05)
}

func TestForTime_PastDate(t *testing.T) {
	past := time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC)
	p := ForTime(past)

	assert.GreaterOrEqual(t, p.DayOfCycle, 0.0)
	assert.LessOrEqual(t, p.DayOfCycle, synodicMonth)
	assert.GreaterOrEqual(t, p.Illumination, 0.0)
	assert.LessOrEqual(t, p.Illumination, 1.0)
}

func TestCurrent_DoesNotPanic(t *testing.T) {
	p := Current()

	assert.NotEmpty(t, string(p.Name))
	assert.NotEmpty(t, p.Emoji)
}

func TestAllPhasesReachable(t *testing.T) {
	seen := make(map[PhaseName]bool)
	for i := range 30 {
		tm := referenceNewMoon.Add(time.Duration(i*24) * time.Hour)
		p := ForTime(tm)
		seen[p.Name] = true
	}

	phases := []PhaseName{
		NewMoon, WaxingCrescent, FirstQuarter, WaxingGibbous,
		FullMoon, WaningGibbous, ThirdQuarter, WaningCrescent,
	}
	for _, name := range phases {
		assert.True(t, seen[name], "phase %s was never reached in a full cycle walk", name)
	}
}

func TestIlluminationSymmetry(t *testing.T) {
	// Illumination at day X should equal illumination at (synodicMonth - X).
	for _, day := range []float64{3, 7, 10, 12} {
		t1 := referenceNewMoon.Add(time.Duration(day * 24 * float64(time.Hour)))
		t2 := referenceNewMoon.Add(time.Duration((synodicMonth - day) * 24 * float64(time.Hour)))
		p1 := ForTime(t1)
		p2 := ForTime(t2)
		assert.InDelta(t, p1.Illumination, p2.Illumination, 0.01,
			"illumination should be symmetric at day %.0f", day)
	}
}

func TestASCII_AllPhasesHaveArt(t *testing.T) {
	phases := []PhaseName{
		NewMoon, WaxingCrescent, FirstQuarter, WaxingGibbous,
		FullMoon, WaningGibbous, ThirdQuarter, WaningCrescent,
	}
	for _, name := range phases {
		p := Phase{Name: name}
		art := p.ASCII()
		for i, line := range art {
			assert.NotEmpty(t, line, "phase %s line %d should not be empty", name, i)
		}
	}
}

func TestASCII_CurrentPhase(t *testing.T) {
	p := Current()
	art := p.ASCII()
	assert.Len(t, art, 3)
	for _, line := range art {
		assert.NotEmpty(t, line)
	}
}

func TestDayOfCycle_Monotonic(t *testing.T) {
	prev := -1.0
	for h := range 710 { // ~29.5 days in hours
		tm := referenceNewMoon.Add(time.Duration(h) * time.Hour)
		p := ForTime(tm)
		if prev >= 0 {
			diff := p.DayOfCycle - prev
			// Allow wrap-around at cycle boundary.
			if math.Abs(diff) > 1 {
				continue
			}
			assert.GreaterOrEqual(t, diff, 0.0, "day of cycle should be monotonically increasing")
		}
		prev = p.DayOfCycle
	}
}
