package astro

import (
	"math"
	"time"
)

// Location represents an observer's position on Earth.
type Location struct {
	Latitude  float64 // degrees, positive north
	Longitude float64 // degrees, positive east
}

// Deg2Rad converts degrees to radians.
func Deg2Rad(d float64) float64 {
	return d * math.Pi / 180.0
}

// Rad2Deg converts radians to degrees.
func Rad2Deg(r float64) float64 {
	return r * 180.0 / math.Pi
}

// NormalizeAngle clamps an angle in degrees to [0, 360).
func NormalizeAngle(deg float64) float64 {
	deg = math.Mod(deg, 360.0)
	if deg < 0 {
		deg += 360.0
	}
	return deg
}

// JulianDate converts a time.Time to Julian Date.
// Algorithm from Meeus, "Astronomical Algorithms" ch. 7.
func JulianDate(t time.Time) float64 {
	t = t.UTC()
	y := float64(t.Year())
	m := float64(t.Month())
	d := float64(t.Day()) +
		float64(t.Hour())/24.0 +
		float64(t.Minute())/1440.0 +
		float64(t.Second())/86400.0

	if m <= 2 {
		y--
		m += 12
	}

	A := math.Floor(y / 100.0)
	B := 2 - A + math.Floor(A/4.0)

	return math.Floor(365.25*(y+4716)) +
		math.Floor(30.6001*(m+1)) +
		d + B - 1524.5
}

// JulianCentury returns Julian centuries since J2000.0 epoch.
func JulianCentury(jd float64) float64 {
	return (jd - 2451545.0) / 36525.0
}

// MeanObliquity returns the mean obliquity of the ecliptic in degrees
// for the given Julian century T.
func MeanObliquity(T float64) float64 {
	return 23.439291 - 0.0130042*T - 1.64e-7*T*T + 5.04e-7*T*T*T
}

// EclipticToEquatorial converts ecliptic longitude and latitude (degrees)
// to right ascension and declination (degrees), given the obliquity (degrees).
func EclipticToEquatorial(lon, lat, obliquity float64) (ra, dec float64) {
	lonR := Deg2Rad(lon)
	latR := Deg2Rad(lat)
	oblR := Deg2Rad(obliquity)

	sinDec := math.Sin(latR)*math.Cos(oblR) + math.Cos(latR)*math.Sin(oblR)*math.Sin(lonR)
	dec = Rad2Deg(math.Asin(sinDec))

	y := math.Sin(lonR)*math.Cos(oblR) - math.Tan(latR)*math.Sin(oblR)
	x := math.Cos(lonR)
	ra = NormalizeAngle(Rad2Deg(math.Atan2(y, x)))

	return ra, dec
}

// LocalSiderealTime returns the local sidereal time in degrees
// for the given Julian Date and observer longitude (degrees, positive east).
func LocalSiderealTime(jd float64, longitude float64) float64 {
	T := JulianCentury(jd)
	// Greenwich Mean Sidereal Time in degrees (Meeus ch. 12)
	gmst := 280.46061837 +
		360.98564736629*(jd-2451545.0) +
		0.000387933*T*T -
		T*T*T/38710000.0

	return NormalizeAngle(gmst + longitude)
}

// AltitudeAzimuth converts hour angle (degrees), declination (degrees),
// and observer latitude (degrees) to altitude and azimuth (both in degrees).
// Azimuth is measured from north, clockwise.
func AltitudeAzimuth(ha, dec, lat float64) (alt, az float64) {
	haR := Deg2Rad(ha)
	decR := Deg2Rad(dec)
	latR := Deg2Rad(lat)

	sinAlt := math.Sin(decR)*math.Sin(latR) + math.Cos(decR)*math.Cos(latR)*math.Cos(haR)
	alt = Rad2Deg(math.Asin(sinAlt))

	cosAz := (math.Sin(decR) - math.Sin(Deg2Rad(alt))*math.Sin(latR)) /
		(math.Cos(Deg2Rad(alt)) * math.Cos(latR))
	// Clamp for floating point errors
	cosAz = math.Max(-1, math.Min(1, cosAz))
	az = Rad2Deg(math.Acos(cosAz))

	if math.Sin(haR) > 0 {
		az = 360.0 - az
	}

	return alt, az
}
