package deepsky

import (
	"fmt"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
)

type Object struct {
	ID   string
	Name string
	RA   float64
	Dec  float64
	Type string
}

type Observation struct {
	Best      *Object
	Altitude  float64
	Azimuth   float64
	IsDark    bool
	Timestamp time.Time
}

var catalog = []Object{
	{ID: "M31", Name: "Andromeda Galaxy", RA: 10.68, Dec: 41.27, Type: "galaxy"},
	{ID: "M42", Name: "Orion Nebula", RA: 83.82, Dec: -5.39, Type: "nebula"},
	{ID: "M45", Name: "Pleiades", RA: 56.75, Dec: 24.12, Type: "cluster"},
	{ID: "M13", Name: "Hercules Cluster", RA: 250.42, Dec: 36.46, Type: "cluster"},
	{ID: "M57", Name: "Ring Nebula", RA: 283.40, Dec: 33.03, Type: "planetary"},
	{ID: "M27", Name: "Dumbbell Nebula", RA: 299.90, Dec: 22.72, Type: "planetary"},
	{ID: "M51", Name: "Whirlpool Galaxy", RA: 202.47, Dec: 47.20, Type: "galaxy"},
	{ID: "M101", Name: "Pinwheel Galaxy", RA: 210.80, Dec: 54.35, Type: "galaxy"},
	{ID: "M8", Name: "Lagoon Nebula", RA: 270.92, Dec: -24.38, Type: "nebula"},
	{ID: "M20", Name: "Trifid Nebula", RA: 270.60, Dec: -23.03, Type: "nebula"},
	{ID: "M44", Name: "Beehive Cluster", RA: 130.03, Dec: 19.67, Type: "cluster"},
	{ID: "M35", Name: "open cluster", RA: 92.25, Dec: 24.33, Type: "cluster"},
	{ID: "NGC 869", Name: "Double Cluster", RA: 34.75, Dec: 57.13, Type: "cluster"},
	{ID: "M1", Name: "Crab Nebula", RA: 83.63, Dec: 22.01, Type: "nebula"},
	{ID: "M33", Name: "Triangulum Galaxy", RA: 23.46, Dec: 30.66, Type: "galaxy"},
}

func Current(loc astro.Location) Observation {
	return ForTime(time.Now(), loc)
}

func ForTime(t time.Time, loc astro.Location) Observation {
	sunAlt, _ := astro.SunPosition(t, loc)
	isDark := sunAlt < -6.0

	jd := astro.JulianDate(t)
	lst := astro.LocalSiderealTime(jd, loc.Longitude)

	obs := Observation{
		IsDark:    isDark,
		Timestamp: t,
	}

	bestAlt := -999.0
	for i := range catalog {
		ha := astro.NormalizeAngle(lst - catalog[i].RA)
		alt, az := astro.AltitudeAzimuth(ha, catalog[i].Dec, loc.Latitude)

		if alt > 15.0 && isDark && alt > bestAlt {
			bestAlt = alt
			obs.Best = &catalog[i]
			obs.Altitude = alt
			obs.Azimuth = az
		}
	}

	return obs
}

func (o Observation) FormatStatus() string {
	if o.Best == nil {
		return ""
	}
	return fmt.Sprintf("🌌 %s (%s) alt %.0f°", o.Best.ID, o.Best.Name, o.Altitude)
}
