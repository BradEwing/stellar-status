package moon

import (
	"math"
	"time"
)

// Synodic month length in days
const synodicMonth = 29.53058770576

var referenceNewMoon = time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

// PhaseName represents the name of a lunar phase.
type PhaseName string

const (
	NewMoon        PhaseName = "New Moon"
	WaxingCrescent PhaseName = "Waxing Crescent"
	FirstQuarter   PhaseName = "First Quarter"
	WaxingGibbous  PhaseName = "Waxing Gibbous"
	FullMoon       PhaseName = "Full Moon"
	WaningGibbous  PhaseName = "Waning Gibbous"
	ThirdQuarter   PhaseName = "Third Quarter"
	WaningCrescent PhaseName = "Waning Crescent"
)

// Phase holds computed lunar phase data.
type Phase struct {
	Name         PhaseName
	Emoji        string
	Illumination float64 // 0.0 to 1.0
	DayOfCycle   float64 // 0.0 to ~29.53
}

// Current returns the moon phase for the current time.
func Current() Phase {
	return ForTime(time.Now())
}

// ForTime returns the moon phase for the given time.
func ForTime(t time.Time) Phase {
	daysSinceRef := t.Sub(referenceNewMoon).Hours() / 24.0
	cycle := math.Mod(daysSinceRef, synodicMonth)
	if cycle < 0 {
		cycle += synodicMonth
	}

	// Illumination: 0 at new moon, 1 at full moon, back to 0.
	// Uses cosine curve mapped from cycle position.
	illumination := (1 - math.Cos(2*math.Pi*cycle/synodicMonth)) / 2

	name, emoji := phaseNameAndEmoji(cycle)

	return Phase{
		Name:         name,
		Emoji:        emoji,
		Illumination: illumination,
		DayOfCycle:   cycle,
	}
}

var asciiPhases = map[PhaseName][3]string{
	NewMoon:        {" . . ", ".   .", " ' ' "},
	WaxingCrescent: {" .*  ", ".(*  ", " '*  "},
	FirstQuarter:   {" .** ", ".(** ", " '** "},
	WaxingGibbous:  {" .***", ".(***", " '***"},
	FullMoon:       {" *** ", "*****", " *** "},
	WaningGibbous:  {"***. ", "***) ", "***' "},
	ThirdQuarter:   {" **. ", " **) ", " **' "},
	WaningCrescent: {"  *. ", "  *) ", "  *' "},
}

// ASCII returns a 3-line minimal-dots ASCII representation of the moon phase.
func (p Phase) ASCII() [3]string {
	if art, ok := asciiPhases[p.Name]; ok {
		return art
	}
	return asciiPhases[NewMoon]
}

func phaseNameAndEmoji(dayOfCycle float64) (PhaseName, string) {
	// Divide the synodic month into 8 equal segments.
	segment := synodicMonth / 8.0
	switch {
	case dayOfCycle < segment*0.5 || dayOfCycle >= segment*7.5:
		return NewMoon, "🌑"
	case dayOfCycle < segment*1.5:
		return WaxingCrescent, "🌒"
	case dayOfCycle < segment*2.5:
		return FirstQuarter, "🌓"
	case dayOfCycle < segment*3.5:
		return WaxingGibbous, "🌔"
	case dayOfCycle < segment*4.5:
		return FullMoon, "🌕"
	case dayOfCycle < segment*5.5:
		return WaningGibbous, "🌖"
	case dayOfCycle < segment*6.5:
		return ThirdQuarter, "🌗"
	default:
		return WaningCrescent, "🌘"
	}
}
