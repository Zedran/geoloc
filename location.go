package geoloc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
)

const (
	// Maximum number of matches returned from the Geocoding API
	MAX_RESP_LOCS     int    = 10

	// A template URL for querying the Geocoding API
	OPEN_WEATHER_URL  string = "https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=%d&appid=%s"
)

var (
    // This error is returned by GetLocation if the Geocoding API does not provide any result for the query
    errLocationNotFound error = errors.New("location not found")

	// The error raised if the user-specified location (-l) is improperly structured
	errLocStringInvalid error = errors.New("improper structure of the specified location")
)

/* This struct represents a location, typically a city. */
type Location struct {
	// Name of a location
	City    string  `json:"name"`
	State   string  `json:"state"`
	Country string  `json:"country"`

	// Coordinates (S and W are negative)
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

/* Estimates distance [km] to another location. */
func (loc *Location) DistanceTo(loc2 *Location) float64 {
	const R = 6371 // Earth's mean radius [km]

	deltaLat := Rad(loc.Lat - loc2.Lat)
	deltaLon := Rad(loc.Lon - loc2.Lon)
	meanLat  := Rad((loc.Lat + loc2.Lat) / 2)

	return R * math.Sqrt(math.Pow(deltaLat, 2) + math.Pow(math.Cos(meanLat) * deltaLon, 2))
}

/* Returns a full name of the location. Optionally, a state can be included if available. */
func (loc *Location) GetName(includeState bool) string {
	if includeState && len(loc.State) > 0 {
		return fmt.Sprintf("%s, %s, %s", loc.City, loc.State, loc.Country)
	}

	return fmt.Sprintf("%s, %s", loc.City, loc.Country)
}

/* Checks whether the two locations overlap. */
func (loc *Location) Overlaps(loc2 *Location) bool {
	// Distance [km] at which the two locations are considered as overlapping
	const OVERLAPPING_D float64 = 10

	return loc.City == loc2.City && loc.DistanceTo(loc2) <= OVERLAPPING_D
}

/* Queries the OpenWeather Geocoding API for the name specified by locName and returns a slice containing matching location names. */
func FindLocation(client *http.Client, keyOW, locName string) ([]Location, error) {
	resp, err := client.Get(fmt.Sprintf(OPEN_WEATHER_URL, locName, MAX_RESP_LOCS, keyOW))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var matches []Location

	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil {
		stream, _ := io.ReadAll(resp.Body)
		fmt.Print(string(stream))
		return nil, err
	} else if len(matches) == 0 {
		return nil, errLocationNotFound
	}

	// Geocoding API has a tendency to return duplicated matches, which need to be filtered out
	if len(matches) > 1 {
		return RemoveOverlappingLocations(matches), nil
	}

	return matches, nil
}

/* Returns a copy of a Location slice without duplicated entries. */
func RemoveOverlappingLocations(matches []Location) []Location {
	var locations []Location

	for i := range matches { // matches returned from API
		unique := true

		for j := range locations { // unique locations
			loc := &locations[j]
			if matches[i].Overlaps(loc) {
				unique = false
				break
			}
		}

		if unique {
			locations = append(locations, matches[i])
		}
	}

	return locations
}

/* Converts the string specified with the -l flag into the Location struct. */
func SpecifyLocation(locString string) (*Location, error) {
	separated := strings.Split(locString, ",")

	if len(separated) != 4 {
		return nil, errLocStringInvalid
	}

	for i := range separated {
		separated[i] = strings.TrimSpace(separated[i])
	}

	var (
		err error
		loc Location
	)

	if len(separated[0]) == 0 {
		loc.City = "Unknown"
	} else {
		loc.City = separated[0]
	}

	if len(separated[1]) == 0 {
		loc.Country = "N/A"
	} else {
		loc.Country = separated[1]
	}
	
	loc.Lat, err = ConvertCoordinate(separated[2], 90)
	if err != nil {
		return nil, err
	}

	loc.Lon, err = ConvertCoordinate(separated[3], 180)
	if err != nil {
		return nil, err
	}

	return &loc, nil
}
