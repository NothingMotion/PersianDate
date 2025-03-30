package persiandate_test

import (
	"testing"
	"time"

	persiandate "github.com/NothingMotion/PersianDate"
)

func TestSort(t *testing.T) {
	pd := persiandate.New("YYYY/MM/DD")

	// Create test dates
	date1 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 6, Day: 15}}
	date2 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1401, Month: 7, Day: 20}}
	date3 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 12, Day: 1}}
	date4 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 1, Day: 5}}

	// Sort the dates
	sortedDates := pd.Sort(date1, date2, date3, date4)

	// Verify order
	if len(sortedDates) != 4 {
		t.Errorf("Sort() returned %d dates, expected 4", len(sortedDates))
	}

	expectedOrder := []persiandate.JalaliDate{
		date2, // 1401-07-20
		date4, // 1402-01-05
		date1, // 1402-06-15
		date3, // 1402-12-01
	}

	for i, expected := range expectedOrder {
		if !pd.Equal(sortedDates[i], expected) {
			t.Errorf("Sort result[%d] = %v, expected %v", i, sortedDates[i], expected)
		}
	}

	// Test with mixed types
	timeDate := time.Date(2023, 3, 21, 0, 0, 0, 0, time.UTC) // 1402-01-01
	stringDate := "1400-12-29"

	mixedSortedDates := pd.Sort(date1, timeDate, stringDate)

	if len(mixedSortedDates) != 3 {
		t.Errorf("Sort() with mixed types returned %d dates, expected 3", len(mixedSortedDates))
	}

	// Expected order: stringDate (1400-12-29), timeDate (1402-01-01), date1 (1402-06-15)
	parsedStringDate, _ := pd.Parse(stringDate)

	if !pd.Equal(mixedSortedDates[0], parsedStringDate) {
		t.Errorf("First sorted date = %v, expected %v", mixedSortedDates[0], parsedStringDate)
	}

	// Expected 1402-01-01 for timeDate
	expectedTimeDate := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 1, Day: 1}}
	if !pd.Equal(mixedSortedDates[1], expectedTimeDate) {
		t.Errorf("Second sorted date = %v, expected %v", mixedSortedDates[1], expectedTimeDate)
	}

	if !pd.Equal(mixedSortedDates[2], date1) {
		t.Errorf("Third sorted date = %v, expected %v", mixedSortedDates[2], date1)
	}
}

func TestFilter(t *testing.T) {
	pd := persiandate.New("YYYY/MM/DD")

	// Create test dates
	date1 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 6, Day: 15}}
	date2 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1401, Month: 7, Day: 20}}
	date3 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1403, Month: 12, Day: 1}}
	date4 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 1, Day: 5}}

	// Filter for dates in 1402
	filteredDates := pd.Filter(
		func(date persiandate.JalaliDate) bool {
			return date.Year == 1402
		},
		date1, date2, date3, date4,
	)

	// Verify filtered results
	if len(filteredDates) != 2 {
		t.Errorf("Filter() returned %d dates, expected 2", len(filteredDates))
	}

	expectedFiltered := []persiandate.JalaliDate{
		date4, // 1402-01-05
		date1, // 1402-06-15
	}

	for i, expected := range expectedFiltered {
		if !pd.Equal(filteredDates[i], expected) {
			t.Errorf("Filter result[%d] = %v, expected %v", i, filteredDates[i], expected)
		}
	}

	// Test filtering with mixed types
	timeDate := time.Date(2023, 3, 21, 0, 0, 0, 0, time.UTC) // 1402-01-01
	stringDate := "1402-12-29"

	// Filter for dates in month 1
	month1Dates := pd.Filter(
		func(date persiandate.JalaliDate) bool {
			return date.Month == 1
		},
		date1, date2, date3, date4, timeDate, stringDate,
	)

	if len(month1Dates) != 2 {
		t.Errorf("Filter() for month 1 returned %d dates, expected 2", len(month1Dates))
	}

	// Expected dates: 1402-01-05 (date4) and 1402-01-01 (timeDate)
	expectedTimeDate := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 1, Day: 1}}

	foundDate4 := false
	foundTimeDate := false

	for _, date := range month1Dates {
		if pd.Equal(date, date4) {
			foundDate4 = true
		}
		if pd.Equal(date, expectedTimeDate) {
			foundTimeDate = true
		}
	}

	if !foundDate4 || !foundTimeDate {
		t.Errorf("Filter() for month 1 did not return expected dates")
	}
}
