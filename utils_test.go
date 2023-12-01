package geoloc

import "testing"

/*
	Tests whether coordinates are converted properly. Badly formatted string values
	and those exceeding the assumed limit should trigger an error.
*/
func TestConvertCoordinate(t *testing.T) {
	const limit float64 = 90

	correctCases := map[string]float64{
		"-23.4": -23.4,
		"0"    :   0,
		"46"   :  46,
		"90"   :  90,
	}

	incorrectCases := []string{"123", "abc", "83.461E"}

	for tCase, val := range correctCases {
		out, err := ConvertCoordinate(tCase, limit)
		if err != nil || out != val {
			t.Errorf("Conversion failed. Expected: %.3f - Got: %.3f - Err: %v", val, out, err)
		}
	}

	for _, tCase := range incorrectCases {
		out, err := ConvertCoordinate(tCase, limit)
		if err == nil {
			t.Errorf("Incorrect value did not trigger error. Got: %.3f from %s.", out, tCase)
		}
	}
}
