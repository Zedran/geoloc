package main

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	SETTINGS_DIR  string = "settings"
	SETTINGS_FILE string = SETTINGS_DIR + "/uv.json"
)

/* This struct represents the application's settings. */
type Settings struct {
	// OpenWeather API key for its geocoding service
	OpenWeatherKey  string    `json:"open_weather_key"`
	
	// OpenUV API key
	OpenUVKey       string    `json:"open_uv_key"`

	// Default location for which the report is generated if no other location is specified
	DefaultLocation *Location `json:"default_location"`
}

/* Loads settings from file. If the file is not present, it calls SaveSettings to generate
 * an empty one.
 */
func LoadSettings() (*Settings, error) {
	stream, err := os.ReadFile(SETTINGS_FILE)
	if errors.Is(err, os.ErrNotExist) {
		return nil, SaveSettings(nil)
	} else if err != nil {
		return nil, err
	}

	var s Settings

	if err = json.Unmarshal(stream, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

/* Saves settings to a file. If s is nil, the empty file is generated. */
func SaveSettings(s *Settings) error {
	if s == nil {
		s = &Settings{DefaultLocation: nil}
	}

	if _, err := os.Stat(SETTINGS_DIR); os.IsNotExist(err) {
		err := os.Mkdir(SETTINGS_DIR, 0644)
		if err != nil {
			return err
		}
	}

	stream, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(SETTINGS_FILE, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(stream)

	return nil
}