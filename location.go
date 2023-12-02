package geoloc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
)

// Default maximum number of matches requested from the Geocoding API
const DEFAULT_MAX_RESP_LOCS int = 10

// This error is returned by GetLocation if the Geocoding API does not provide any result for the query
var errLocationNotFound error = errors.New("location not found")

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

/*
	Queries the OpenWeather Geocoding API for the name specified by locName and returns a slice containing matching location names.
	maxRespLocs limits the number of matches requested from API. If maxRespLocs is equal or less than 0, it is set to DEFAULT_MAX_RESP_LOCS.
*/
func FindLocation(client *http.Client, keyOW, locName string, maxRespLocs int) ([]Location, error) {
	// A template URL for querying the Geocoding API
	const OPEN_WEATHER_URL string = "https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=%d&appid=%s"

	if maxRespLocs <= 0 {
		maxRespLocs = DEFAULT_MAX_RESP_LOCS
	}

	resp, err := client.Get(fmt.Sprintf(OPEN_WEATHER_URL, locName, maxRespLocs, keyOW))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var matches []Location

	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil {
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
