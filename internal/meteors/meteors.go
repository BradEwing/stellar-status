package meteors

import (
	"fmt"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
)

type Shower struct {
	Name       string
	PeakMonth  time.Month
	PeakDay    int
	StartMonth time.Month
	StartDay   int
	EndMonth   time.Month
	EndDay     int
	ZHR        int
}

type Status struct {
	Shower      *Shower
	DaysToPeak  int
	Active      bool
	PeakTonight bool
	Timestamp   time.Time
}

var showers = []Shower{
	{"Quadrantids", time.January, 4, time.December, 28, time.January, 12, 110},
	{"Lyrids", time.April, 22, time.April, 14, time.April, 30, 18},
	{"Eta Aquariids", time.May, 6, time.April, 19, time.May, 28, 50},
	{"Southern Delta Aquariids", time.July, 30, time.July, 12, time.August, 23, 25},
	{"Alpha Capricornids", time.July, 30, time.July, 3, time.August, 15, 5},
	{"Perseids", time.August, 12, time.July, 17, time.August, 24, 100},
	{"Draconids", time.October, 8, time.October, 6, time.October, 10, 10},
	{"Orionids", time.October, 21, time.October, 2, time.November, 7, 20},
	{"Leonids", time.November, 17, time.November, 6, time.November, 30, 15},
	{"Geminids", time.December, 14, time.December, 4, time.December, 20, 150},
	{"Ursids", time.December, 22, time.December, 17, time.December, 26, 10},
	{"Taurids", time.November, 5, time.September, 10, time.November, 20, 5},
}

func Current(loc astro.Location) Status {
	return ForTime(time.Now(), loc)
}

func ForTime(t time.Time, _ astro.Location) Status {
	result := Status{Timestamp: t}

	var bestActive *Shower
	bestActiveDist := 366

	var bestUpcoming *Shower
	bestUpcomingDist := 366

	for i := range showers {
		s := &showers[i]
		days := daysToPeak(t, s)
		active := isActive(t, s)

		if active {
			absDays := days
			if absDays < 0 {
				absDays = -absDays
			}
			if absDays < bestActiveDist {
				bestActiveDist = absDays
				bestActive = s
			}
		}

		if days >= 0 && days <= 7 && days < bestUpcomingDist {
			bestUpcomingDist = days
			bestUpcoming = s
		}
	}

	if bestActive != nil {
		result.Shower = bestActive
		result.DaysToPeak = daysToPeak(t, bestActive)
		result.Active = true
		result.PeakTonight = result.DaysToPeak == 0
	} else if bestUpcoming != nil {
		result.Shower = bestUpcoming
		result.DaysToPeak = daysToPeak(t, bestUpcoming)
		result.Active = false
	}

	return result
}

func (s Status) FormatStatus() string {
	if s.Shower == nil {
		return ""
	}

	if s.PeakTonight {
		return fmt.Sprintf("☄️ %s peak tonight (ZHR ~%d)", s.Shower.Name, s.Shower.ZHR)
	}

	if s.Active && s.DaysToPeak > 0 && s.DaysToPeak <= 3 {
		return fmt.Sprintf("☄️ %s peak in %dd (ZHR ~%d)", s.Shower.Name, s.DaysToPeak, s.Shower.ZHR)
	}

	if s.Active {
		peakDate := fmt.Sprintf("%s %d", s.Shower.PeakMonth.String()[:3], s.Shower.PeakDay)
		return fmt.Sprintf("☄️ %s active (peak %s)", s.Shower.Name, peakDate)
	}

	if s.DaysToPeak > 0 && s.DaysToPeak <= 7 {
		return fmt.Sprintf("☄️ %s in %dd", s.Shower.Name, s.DaysToPeak)
	}

	return ""
}

func daysToPeak(t time.Time, s *Shower) int {
	peak := time.Date(t.Year(), s.PeakMonth, s.PeakDay, 0, 0, 0, 0, t.Location())
	current := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	diff := int(peak.Sub(current).Hours() / 24)
	if diff < -182 {
		diff += 365
	} else if diff > 182 {
		diff -= 365
	}
	return diff
}

func isActive(t time.Time, s *Shower) bool {
	year := t.Year()
	start := time.Date(year, s.StartMonth, s.StartDay, 0, 0, 0, 0, t.Location())
	end := time.Date(year, s.EndMonth, s.EndDay, 23, 59, 59, 0, t.Location())
	current := time.Date(year, t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if start.Before(end) || start.Equal(end) {
		return !current.Before(start) && !current.After(end)
	}
	return !current.Before(start) || !current.After(end)
}
