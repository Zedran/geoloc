package geoloc

import "testing"

/*
	This test aims to check whether the specified string is properly formed into the Location struct by SpecifyLocation function.
	It does not check coordinate conversion.
*/
func TestSpecifyLocation(t *testing.T) {
	correctCases := map[string]Location {
		"London, GB, 51.508, -0.128"     : {"London",  "", "GB",  51.508,  -0.128},
		", ,51.508,-0.128 "              : {"Unknown", "", "N/A", 51.508,  -0.128},
		" Sendai , JP , 38.252, 140.856" : {"Sendai",  "", "JP",  38.252, 140.856},
	}

	incorrectCases := []string{
		"London, GB, 51,508, -0,128",
		"51.508, -0.128",
	}

	for tCase, val := range correctCases {
		out, err := SpecifyLocation(tCase)
		if err != nil || !(out.GetName(true) == val.GetName(true) && out.DistanceTo(&val) == 0) {
			t.Errorf(
				"Parsing failed:\n Got      : %s (%.3f, %.3f)\n Expected : %s (%.3f, %.3f)\n Err: %v",
				out.GetName(true), out.Lat, out.Lon, val.GetName(true), val.Lat, val.Lon, err,
			)
		}
	}

	for _, tCase := range incorrectCases {
		out, err := SpecifyLocation(tCase)
		if err == nil {
			t.Errorf(
				"Incorrect value did not trigger error:\n Got : {%s (%.3f, %.3f)} from \"%s\".",
				out.GetName(true), out.Lat, out.Lon, tCase,
			)
		}
	}
}
