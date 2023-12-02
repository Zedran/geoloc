package geoloc

import "math"

/* Converts degrees to radians. */
func Rad(deg float64) float64 {
	return deg * math.Pi / 180
}
