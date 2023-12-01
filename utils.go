package geoloc

import (
	"errors"
	"math"
	"strconv"
)

var errCoordinateOverflow    error  = errors.New("coordinate value out of range")

/*
	Converts geographical coordinate from string to float. Errors are returned in case of parsing problems 
	or if the coordinate's value limit was exceeded. Value limits should correspond to the coordinate type (lat == 90, lon == 180).
*/
func ConvertCoordinate(coordString string, limit float64) (float64, error) {
	coordNum, err := strconv.ParseFloat(coordString, 64)
	if err != nil {
		return 0, err
	}

	if math.Abs(coordNum) > limit {
		return 0, errCoordinateOverflow
	}

	return coordNum, nil
}

/* Converts degrees to radians. */
func Rad(deg float64) float64 {
	return deg * math.Pi / 180
}
