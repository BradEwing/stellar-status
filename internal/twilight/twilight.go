package twilight

import (
	"fmt"
	"math"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
)

const (
	sunriseAlt      = -0.833 // accounts for refraction + solar radius
	civilAlt        = -6.0
	astroAlt        = -18.0
	goldenHourAlt   = 6.0
)

// SunTimes holds sunrise, sunset, and twilight times for a date and location.
type SunTimes struct {
	Sunrise         time.Time
	Sunset          time.Time
	CivilDawn       time.Time
	CivilDusk       time.Time
	AstroDawn       time.Time
	AstroDusk       time.Time
	GoldenHourStart time.Time
	NeverRises      bool
	NeverSets       bool
	Location        astro.Location
	Date            time.Time
}

// Today returns sun times for today at the given location.
func Today(loc astro.Location) SunTimes {
	return ForDate(time.Now(), loc)
}

// ForDate returns sun times for the given date and location.
func ForDate(date time.Time, loc astro.Location) SunTimes {
	noon := time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, time.UTC)
	jdNoon := astro.JulianDate(noon)
	T := astro.JulianCentury(jdNoon)

	ra, dec := astro.SunEquatorial(noon)
	lst := astro.LocalSiderealTime(jdNoon, loc.Longitude)
	transitHA := astro.NormalizeAngle(lst - ra)
	if transitHA > 180 {
		transitHA -= 360
	}
	transitOffset := -transitHA / 360.0 * 24.0
	transit := noon.Add(time.Duration(transitOffset * float64(time.Hour)))

	_ = T

	st := SunTimes{
		Location: loc,
		Date:     date,
	}

	rise, set, ok := hourAngleToTimes(transit, dec, loc.Latitude, sunriseAlt)
	if !ok {
		alt, _ := astro.SunPosition(noon, loc)
		if alt > 0 {
			st.NeverSets = true
		} else {
			st.NeverRises = true
		}
		return st
	}
	st.Sunrise = rise
	st.Sunset = set

	if dawn, dusk, ok := hourAngleToTimes(transit, dec, loc.Latitude, civilAlt); ok {
		st.CivilDawn = dawn
		st.CivilDusk = dusk
	}

	if dawn, dusk, ok := hourAngleToTimes(transit, dec, loc.Latitude, astroAlt); ok {
		st.AstroDawn = dawn
		st.AstroDusk = dusk
	}

	if ghStart, _, ok := hourAngleToTimes(transit, dec, loc.Latitude, goldenHourAlt); ok {
		_ = ghStart
		st.GoldenHourStart = time.Date(
			set.Year(), set.Month(), set.Day(),
			set.Hour(), set.Minute(), set.Second(), 0, set.Location(),
		)
		haH := solveHourAngle(dec, loc.Latitude, goldenHourAlt)
		ghDuration := time.Duration(haH / 180.0 * 12.0 * float64(time.Hour))
		st.GoldenHourStart = set.Add(-ghDuration + time.Duration(solveHourAngle(dec, loc.Latitude, sunriseAlt)/180.0*12.0*float64(time.Hour)))
	}

	return st
}

func solveHourAngle(dec, lat, targetAlt float64) float64 {
	decR := astro.Deg2Rad(dec)
	latR := astro.Deg2Rad(lat)
	altR := astro.Deg2Rad(targetAlt)

	cosH := (math.Sin(altR) - math.Sin(latR)*math.Sin(decR)) /
		(math.Cos(latR) * math.Cos(decR))

	if cosH < -1 || cosH > 1 {
		return -1
	}
	return astro.Rad2Deg(math.Acos(cosH))
}

func hourAngleToTimes(transit time.Time, dec, lat, targetAlt float64) (rise, set time.Time, ok bool) {
	H := solveHourAngle(dec, lat, targetAlt)
	if H < 0 {
		return time.Time{}, time.Time{}, false
	}

	offset := time.Duration(H / 360.0 * 24.0 * float64(time.Hour))
	return transit.Add(-offset), transit.Add(offset), true
}

// FormatStatus returns a context-dependent status string based on the current time.
func (s SunTimes) FormatStatus(now time.Time) string {
	if s.NeverSets {
		return "☀️ midnight sun"
	}
	if s.NeverRises {
		return "🌃 polar night"
	}

	localNow := now.In(time.Local)

	switch {
	case !s.AstroDawn.IsZero() && localNow.Before(s.AstroDawn):
		return fmt.Sprintf("🌃 dark til %s", fmtTime(s.AstroDawn))
	case !s.CivilDawn.IsZero() && localNow.Before(s.CivilDawn):
		return fmt.Sprintf("🌆 dawn %s", fmtTime(s.CivilDawn))
	case localNow.Before(s.Sunrise):
		return fmt.Sprintf("🌅 sunrise %s", fmtTime(s.Sunrise))
	case !s.GoldenHourStart.IsZero() && localNow.Before(s.GoldenHourStart):
		return fmt.Sprintf("🌅 sunset %s", fmtTime(s.Sunset))
	case localNow.Before(s.Sunset):
		return "🌇 golden hour"
	case !s.CivilDusk.IsZero() && localNow.Before(s.CivilDusk):
		return fmt.Sprintf("🌆 dusk til %s", fmtTime(s.CivilDusk))
	case !s.AstroDusk.IsZero() && localNow.Before(s.AstroDusk):
		return fmt.Sprintf("🌃 twilight til %s", fmtTime(s.AstroDusk))
	default:
		return "🌃 night"
	}
}

func fmtTime(t time.Time) string {
	local := t.In(time.Local)
	return local.Format("3:04pm")
}
