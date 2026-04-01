package solar

import (
	"fmt"
	"math"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
)

// Position holds the sun's current position relative to the observer.
type Position struct {
	Altitude  float64
	Azimuth   float64
	IsDay     bool
	Timestamp time.Time
}

// Current returns the sun's position for the current time.
func Current(loc astro.Location) Position {
	return ForTime(time.Now(), loc)
}

// ForTime returns the sun's position for the given time and location.
func ForTime(t time.Time, loc astro.Location) Position {
	alt, az := astro.SunPosition(t, loc)
	return Position{
		Altitude:  alt,
		Azimuth:   az,
		IsDay:     alt > 0,
		Timestamp: t,
	}
}

// FormatStatus returns a compact status string for the status line.
func (p Position) FormatStatus() string {
	altInt := int(math.Round(p.Altitude))
	switch {
	case p.Altitude > 6:
		return fmt.Sprintf("☀️ %d°", altInt)
	case p.Altitude > -6:
		return fmt.Sprintf("🌅 %d°", altInt)
	default:
		return fmt.Sprintf("🌙 %d°", altInt)
	}
}
