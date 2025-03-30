package persiandate_test

import (
	"testing"

	persiandate "github.com/NothingMotion/PersianDate"
)

func TestMin(t *testing.T) {
	pd := persiandate.New("")

	dates := []persiandate.JalaliDate{
		{Date: persiandate.Date{Year: 1402, Month: 1, Day: 1}},
		{Date: persiandate.Date{Year: 1403, Month: 1, Day: 1}},
		{Date: persiandate.Date{Year: 1404, Month: 1, Day: 1}},
		{Date: persiandate.Date{Year: 1404, Month: 1, Day: 12}},
	}
	oldest := pd.Min(dates)

	t.Logf(oldest.String())

	// Check that we got the oldest date
	if oldest.Year != 1402 || oldest.Month != 1 || oldest.Day != 1 {
		t.Errorf("Expected oldest date to be 1402-01-01, got %s", oldest.String())
	}
}

func TestMax(t *testing.T) {
	pd := persiandate.New("")

	dates := []persiandate.JalaliDate{
		{Date: persiandate.Date{Year: 1402, Month: 1, Day: 1}},
		{Date: persiandate.Date{Year: 1403, Month: 1, Day: 1}},
		{Date: persiandate.Date{Year: 1404, Month: 1, Day: 1}},
		{Date: persiandate.Date{Year: 1404, Month: 1, Day: 12}},
	}
	newest := pd.Max(dates)

	t.Logf(newest.String())

	// Check that we got the newest date
	if newest.Year != 1404 || newest.Month != 1 || newest.Day != 12 {
		t.Errorf("Expected newest date to be 1404-01-12, got %s", newest.String())
	}
}

func TestMinMaxVariadic(t *testing.T) {
	pd := persiandate.New("")

	date1 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 1, Day: 1}}
	date2 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1403, Month: 1, Day: 1}}
	date3 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1404, Month: 1, Day: 1}}
	date4 := persiandate.JalaliDate{Date: persiandate.Date{Year: 1404, Month: 1, Day: 12}}

	// Test Min with variadic arguments
	oldest := pd.Min(date1, date2, date3, date4)
	if oldest.Year != 1402 || oldest.Month != 1 || oldest.Day != 1 {
		t.Errorf("Expected oldest date to be 1402-01-01, got %s", oldest.String())
	}

	// Test Max with variadic arguments
	newest := pd.Max(date1, date2, date3, date4)
	if newest.Year != 1404 || newest.Month != 1 || newest.Day != 12 {
		t.Errorf("Expected newest date to be 1404-01-12, got %s", newest.String())
	}
}
