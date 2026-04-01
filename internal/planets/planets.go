package planets

import (
	"math"
	"strings"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
)

// PlanetName identifies a planet.
type PlanetName string

const (
	Mercury PlanetName = "Mercury"
	Venus   PlanetName = "Venus"
	Earth   PlanetName = "Earth"
	Mars    PlanetName = "Mars"
	Jupiter PlanetName = "Jupiter"
	Saturn  PlanetName = "Saturn"
)

var visiblePlanets = []PlanetName{Mercury, Venus, Mars, Jupiter, Saturn}

// PlanetStatus holds a planet's position and visibility.
type PlanetStatus struct {
	Name     PlanetName
	Altitude float64
	Azimuth  float64
	Visible  bool
}

// Visibility holds the visibility status of all naked-eye planets.
type Visibility struct {
	Planets   []PlanetStatus
	Timestamp time.Time
}

// Current returns planet visibility for the current time.
func Current(loc astro.Location) Visibility {
	return ForTime(time.Now(), loc)
}

// ForTime returns planet visibility for the given time and location.
func ForTime(t time.Time, loc astro.Location) Visibility {
	jd := astro.JulianDate(t)
	T := astro.JulianCentury(jd)

	sunAlt, _ := astro.SunPosition(t, loc)
	isDark := sunAlt < -6.0

	earthLon, earthLat, earthR := heliocentricPosition(Earth, T)

	var statuses []PlanetStatus
	for _, name := range visiblePlanets {
		pLon, pLat, pR := heliocentricPosition(name, T)

		geoLon, geoLat := toGeocentric(pLon, pLat, pR, earthLon, earthLat, earthR)

		obl := astro.MeanObliquity(T)
		ra, dec := astro.EclipticToEquatorial(geoLon, geoLat, obl)

		lst := astro.LocalSiderealTime(jd, loc.Longitude)
		ha := astro.NormalizeAngle(lst - ra)
		alt, az := astro.AltitudeAzimuth(ha, dec, loc.Latitude)

		statuses = append(statuses, PlanetStatus{
			Name:     name,
			Altitude: alt,
			Azimuth:  az,
			Visible:  alt > 5.0 && isDark,
		})
	}

	return Visibility{Planets: statuses, Timestamp: t}
}

// FormatStatus returns a compact status string listing visible planets.
func (v Visibility) FormatStatus() string {
	var names []string
	for _, p := range v.Planets {
		if p.Visible {
			names = append(names, string(p.Name))
		}
	}
	if len(names) == 0 {
		return "🔭 —"
	}
	return "🔭 " + strings.Join(names, " ")
}

func heliocentricPosition(name PlanetName, T float64) (lon, lat, r float64) {
	el := elements[name]

	L := astro.NormalizeAngle(el.L0 + el.L1*T)
	e := el.e0 + el.e1*T
	I := el.I0 + el.I1*T
	w := el.w0 + el.w1*T
	O := el.O0 + el.O1*T

	// Mean anomaly
	M := astro.NormalizeAngle(L - w)
	Mr := astro.Deg2Rad(M)

	// Solve Kepler's equation: E - e*sin(E) = M (iterate)
	E := Mr
	for range 10 {
		E = Mr + e*math.Sin(E)
	}

	// True anomaly
	cosV := (math.Cos(E) - e) / (1 - e*math.Cos(E))
	sinV := math.Sqrt(1-e*e) * math.Sin(E) / (1 - e*math.Cos(E))
	v := math.Atan2(sinV, cosV)

	// Heliocentric distance
	r = el.a * (1 - e*math.Cos(E))

	// Heliocentric ecliptic coordinates
	wR := astro.Deg2Rad(w - O)
	OR := astro.Deg2Rad(O)
	IR := astro.Deg2Rad(I)

	cosO := math.Cos(OR)
	sinO := math.Sin(OR)
	cosI := math.Cos(IR)
	sinI := math.Sin(IR)
	cosWV := math.Cos(wR + v)
	sinWV := math.Sin(wR + v)

	x := r * (cosO*cosWV - sinO*sinWV*cosI)
	y := r * (sinO*cosWV + cosO*sinWV*cosI)
	z := r * sinWV * sinI

	lon = astro.NormalizeAngle(astro.Rad2Deg(math.Atan2(y, x)))
	lat = astro.Rad2Deg(math.Asin(z / r))
	return lon, lat, r
}

func toGeocentric(pLon, pLat, pR, eLon, eLat, eR float64) (geoLon, geoLat float64) {
	pLonR := astro.Deg2Rad(pLon)
	pLatR := astro.Deg2Rad(pLat)
	eLonR := astro.Deg2Rad(eLon)
	eLatR := astro.Deg2Rad(eLat)

	px := pR*math.Cos(pLatR)*math.Cos(pLonR) - eR*math.Cos(eLatR)*math.Cos(eLonR)
	py := pR*math.Cos(pLatR)*math.Sin(pLonR) - eR*math.Cos(eLatR)*math.Sin(eLonR)
	pz := pR*math.Sin(pLatR) - eR*math.Sin(eLatR)

	geoLon = astro.NormalizeAngle(astro.Rad2Deg(math.Atan2(py, px)))
	dist := math.Sqrt(px*px + py*py + pz*pz)
	geoLat = astro.Rad2Deg(math.Asin(pz / dist))
	return geoLon, geoLat
}
