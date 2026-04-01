package astro

import (
	"math"
	"time"
)

// SunMeanAnomaly returns the sun's mean anomaly in degrees for Julian century T.
func SunMeanAnomaly(T float64) float64 {
	return NormalizeAngle(357.52911 + 35999.05029*T - 0.0001537*T*T)
}

// SunEclipticLongitude returns the sun's ecliptic longitude in degrees
// for Julian century T.
func SunEclipticLongitude(T float64) float64 {
	M := SunMeanAnomaly(T)
	Mr := Deg2Rad(M)

	// Equation of center
	C := (1.9146-0.004817*T-0.000014*T*T)*math.Sin(Mr) +
		(0.019993-0.000101*T)*math.Sin(2*Mr) +
		0.00029*math.Sin(3*Mr)

	// Sun's true longitude
	omega := 125.04 - 1934.136*T
	lon := NormalizeAngle(M + C + 180.0 + 102.93735)

	// Apparent longitude (nutation correction)
	lon = lon - 0.00569 - 0.00478*math.Sin(Deg2Rad(omega))
	return NormalizeAngle(lon)
}

// SunEquatorial returns the sun's right ascension and declination (degrees)
// for the given time.
func SunEquatorial(t time.Time) (ra, dec float64) {
	jd := JulianDate(t)
	T := JulianCentury(jd)

	lon := SunEclipticLongitude(T)
	obl := MeanObliquity(T)

	// Sun is on the ecliptic, so ecliptic latitude = 0
	return EclipticToEquatorial(lon, 0, obl)
}

// SunPosition returns the sun's altitude and azimuth (degrees) for the
// given time and observer location.
func SunPosition(t time.Time, loc Location) (alt, az float64) {
	ra, dec := SunEquatorial(t)

	jd := JulianDate(t)
	lst := LocalSiderealTime(jd, loc.Longitude)
	ha := NormalizeAngle(lst - ra)

	return AltitudeAzimuth(ha, dec, loc.Latitude)
}
