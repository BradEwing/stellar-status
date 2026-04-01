package planets

// orbitalElements holds mean orbital elements and their rates per Julian century.
// Source: Standish (1992), JPL "Approximate Positions of the Planets".
// Valid for 1800-2050 with ~1 degree accuracy.
type orbitalElements struct {
	a      float64 // semi-major axis (AU)
	e0, e1 float64 // eccentricity and rate per century
	I0, I1 float64 // inclination (degrees) and rate
	L0, L1 float64 // mean longitude (degrees) and rate
	w0, w1 float64 // longitude of perihelion (degrees) and rate
	O0, O1 float64 // longitude of ascending node (degrees) and rate
}

var elements = map[PlanetName]orbitalElements{
	Mercury: {
		a: 0.38709927, e0: 0.20563593, e1: 0.00001906,
		I0: 7.00497902, I1: -0.00594749,
		L0: 252.25032350, L1: 149472.67411175,
		w0: 77.45779628, w1: 0.16047689,
		O0: 48.33076593, O1: -0.12534081,
	},
	Venus: {
		a: 0.72333566, e0: 0.00677672, e1: -0.00004107,
		I0: 3.39467605, I1: -0.00078890,
		L0: 181.97909950, L1: 58517.81538729,
		w0: 131.60246718, w1: 0.00268329,
		O0: 76.67984255, O1: -0.27769418,
	},
	Earth: {
		a: 1.00000261, e0: 0.01671123, e1: -0.00004392,
		I0: -0.00001531, I1: -0.01294668,
		L0: 100.46457166, L1: 35999.37244981,
		w0: 102.93768193, w1: 0.32327364,
		O0: 0.0, O1: 0.0,
	},
	Mars: {
		a: 1.52371034, e0: 0.09339410, e1: 0.00007882,
		I0: 1.84969142, I1: -0.00813131,
		L0: -4.55343205, L1: 19140.30268499,
		w0: -23.94362959, w1: 0.44441088,
		O0: 49.55953891, O1: -0.29257343,
	},
	Jupiter: {
		a: 5.20288700, e0: 0.04838624, e1: -0.00013253,
		I0: 1.30439695, I1: -0.00183714,
		L0: 34.39644051, L1: 3034.74612775,
		w0: 14.72847983, w1: 0.21252668,
		O0: 100.47390909, O1: 0.20469106,
	},
	Saturn: {
		a: 9.53667594, e0: 0.05386179, e1: -0.00050991,
		I0: 2.48599187, I1: 0.00193609,
		L0: 49.95424423, L1: 1222.49362201,
		w0: 92.59887831, w1: -0.41897216,
		O0: 113.66242448, O1: -0.28867794,
	},
}
