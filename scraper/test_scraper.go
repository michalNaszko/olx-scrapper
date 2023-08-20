package scraper

import "testing"

type extractDateTest struct {
	placeAndDate, expected string
}

var extractDateTests = []extractDateTest{
	{"Wrocław, Stare Miasto - Dzisiaj o 14:21", "2023-08-20 14:21"},
	{"Kiełczów - Dzisiaj o 13:16", "2023-08-20 13:16"},
	{"Olsztyn - 14 sierpnia 2023", "2023-08-14 00:00"},
	{"Toruń - 19 lutego 2023", "2023-02-19 00:00"},
}

func TestExtractDate(t *testing.T) {
	for _, test := range extractDateTests {
		if output := extractDate(test.placeAndDate); output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}
