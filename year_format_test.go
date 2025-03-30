package persiandate_test

import (
	"testing"

	persiandate "github.com/NothingMotion/PersianDate"
)

func TestYearFormatting(t *testing.T) {
	// Create a PersianDate instance with an empty format string
	// (we'll specify format in each test)
	pd := persiandate.New("")

	// Set a specific date to test with - use a valid Jalali date
	testDate := pd.ToJalali(2023, 9, 6).Date() // 1402-06-15

	// Test different year format options
	testCases := []struct {
		format   string
		expected string
	}{
		// Year formats
		{"YYYY", "1402"}, // Full year with leading zeros
		{"YYY", "402"},   // 3-digit year
		{"YY", "02"},     // 2-digit year with leading zero
		{"Y", "2"},       // 2-digit year without leading zero
		{"y", "1402"},    // Full year without padding

		// Test with mixed formats
		{"YYYY/MM/DD", "1402/06/15"},  // Standard format
		{"YY-MM-DD", "02-06-15"},      // Short year format
		{"y/M/D", "1402/6/15"},        // Without padding
		{"YYYY MM y", "1402 06 1402"}, // Multiple year formats
	}

	for _, tc := range testCases {
		pd.FORMAT = tc.format
		result := pd.Format(testDate, false)
		if result != tc.expected {
			t.Errorf("Format(%s) = %s, expected %s", tc.format, result, tc.expected)
		}
	}

	// Test hour format options separately
	hourFormatCases := []struct {
		hour     int
		format   string
		expected string
	}{
		{0, "h", "12"},   // 12am
		{0, "hh", "12"},  // 12am with leading zero
		{1, "h", "1"},    // 1am
		{1, "hh", "01"},  // 1am with leading zero
		{12, "h", "12"},  // 12pm
		{12, "hh", "12"}, // 12pm with leading zero
		{13, "h", "1"},   // 1pm
		{13, "hh", "01"}, // 1pm with leading zero
		{23, "h", "11"},  // 11pm
		{23, "hh", "11"}, // 11pm with leading zero
	}

	baseJalaliDate := testDate

	for _, tc := range hourFormatCases {
		// Create a JalaliDate with the specific hour
		testDateWithHour := persiandate.JalaliDate{
			Date: persiandate.Date{
				Year:   baseJalaliDate.Year,
				Month:  baseJalaliDate.Month,
				Day:    baseJalaliDate.Day,
				Hour:   tc.hour,
				Minute: 0,
				Second: 0,
			},
		}

		// Set format and get result
		pd.FORMAT = tc.format
		result := pd.Format(testDateWithHour, false)

		if result != tc.expected {
			t.Errorf("Format(%s) at hour %d = %s, expected %s",
				tc.format, tc.hour, result, tc.expected)
		}
	}
}
