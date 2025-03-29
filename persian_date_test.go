package persiandate_go_test

import (
	"testing"
	"time"

	persiandate "github.com/NothingMotion/PersianDate-GO"
)

func TestJalaliConversion(t *testing.T) {
	pd := persiandate.NewPersianDate("YYYY/MM/DD")

	tests := []struct {
		gregorianDate time.Time
		expectedYear  int
		expectedMonth int
		expectedDay   int
	}{
		{time.Date(2023, 3, 21, 0, 0, 0, 0, time.UTC), 1402, 1, 1},   // Nowruz
		{time.Date(2023, 9, 23, 0, 0, 0, 0, time.UTC), 1402, 7, 1},   // First day of Mehr
		{time.Date(2023, 12, 22, 0, 0, 0, 0, time.UTC), 1402, 10, 1}, // First day of Dey
		{time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC), 1402, 12, 1},  // First day of Esfand
		{time.Date(2024, 3, 19, 0, 0, 0, 0, time.UTC), 1402, 12, 29}, // Last day of Persian year 1402
		{time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC), 1403, 1, 1},   // Nowruz 1403
		{time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), 1378, 10, 11},  // Y2K
		{time.Date(1979, 2, 11, 0, 0, 0, 0, time.UTC), 1357, 11, 22}, // Islamic Revolution of Iran
	}

	for _, test := range tests {
		result := pd.JalaliFull(test.gregorianDate)
		if result.Year != test.expectedYear || result.Month != test.expectedMonth || result.Day != test.expectedDay {
			t.Errorf("JalaliFull(%v) = %d-%02d-%02d, expected %d-%02d-%02d",
				test.gregorianDate, result.Year, result.Month, result.Day,
				test.expectedYear, test.expectedMonth, test.expectedDay)
		}

		jalaliDate := pd.Jalali(test.gregorianDate)
		if jalaliDate.Year != test.expectedYear || jalaliDate.Month != test.expectedMonth || jalaliDate.Day != test.expectedDay {
			t.Errorf("Jalali(%v) = %d-%02d-%02d, expected %d-%02d-%02d",
				test.gregorianDate, jalaliDate.Year, jalaliDate.Month, jalaliDate.Day,
				test.expectedYear, test.expectedMonth, test.expectedDay)
		}
	}
}

func TestGregorianConversion(t *testing.T) {
	pd := persiandate.NewPersianDate("YYYY/MM/DD")

	tests := []struct {
		jalaliYear    int
		jalaliMonth   int
		jalaliDay     int
		expectedYear  int
		expectedMonth int
		expectedDay   int
	}{
		{1402, 1, 1, 2023, 3, 21},   // Nowruz
		{1402, 7, 1, 2023, 9, 23},   // First day of Mehr
		{1402, 10, 1, 2023, 12, 22}, // First day of Dey
		{1402, 12, 1, 2024, 2, 20},  // First day of Esfand
		{1402, 12, 29, 2024, 3, 19}, // Last day of Persian year 1402
		{1403, 1, 1, 2024, 3, 20},   // Nowruz 1403
		{1378, 10, 11, 2000, 1, 1},  // Y2K
		{1357, 11, 22, 1979, 2, 11}, // Islamic Revolution of Iran
	}

	for _, test := range tests {
		gregorianDate := pd.ToGregorian(test.jalaliYear, test.jalaliMonth, test.jalaliDay)
		if gregorianDate.Year != test.expectedYear || gregorianDate.Date.Month != test.expectedMonth || gregorianDate.Date.Day != test.expectedDay {
			t.Errorf("ToGregorian(%d, %d, %d) = %d-%02d-%02d, expected %d-%02d-%02d",
				test.jalaliYear, test.jalaliMonth, test.jalaliDay,
				gregorianDate.Date.Year, gregorianDate.Date.Month, gregorianDate.Date.Day,
				test.expectedYear, test.expectedMonth, test.expectedDay)
		}
	}
}

func TestLeapYears(t *testing.T) {
	pd := persiandate.NewPersianDate("")

	jalaliLeapYears := []int{1375, 1379, 1383, 1387, 1391, 1395, 1399, 1403}
	jalaliNonLeapYears := []int{1376, 1377, 1378, 1380, 1381, 1382, 1384}

	for _, year := range jalaliLeapYears {
		if !pd.IsLeapYearJalali(year) {
			t.Errorf("Year %d should be a leap year in Jalali calendar", year)
		}
	}

	for _, year := range jalaliNonLeapYears {
		if pd.IsLeapYearJalali(year) {
			t.Errorf("Year %d should not be a leap year in Jalali calendar", year)
		}
	}

	gregorianLeapYears := []int{2000, 2004, 2008, 2012, 2016, 2020, 2024}
	gregorianNonLeapYears := []int{1900, 2001, 2002, 2003, 2005, 2100}

	for _, year := range gregorianLeapYears {
		if !pd.IsLeapYearGregorian(year) {
			t.Errorf("Year %d should be a leap year in Gregorian calendar", year)
		}
	}

	for _, year := range gregorianNonLeapYears {
		if pd.IsLeapYearGregorian(year) {
			t.Errorf("Year %d should not be a leap year in Gregorian calendar", year)
		}
	}
}

func TestMonthLength(t *testing.T) {
	pd := persiandate.NewPersianDate("")

	// Test Jalali month lengths
	// First 6 months should have 31 days
	for month := 1; month <= 6; month++ {
		if pd.JalaliMonthLength(1402, month) != 31 {
			t.Errorf("Month %d in year 1402 should have 31 days", month)
		}
	}

	// Next 5 months should have 30 days
	for month := 7; month <= 11; month++ {
		if pd.JalaliMonthLength(1402, month) != 30 {
			t.Errorf("Month %d in year 1402 should have 30 days", month)
		}
	}

	// Last month should have 29 days in non-leap years
	if pd.JalaliMonthLength(1402, 12) != 29 {
		t.Errorf("Month 12 in year 1402 should have 29 days")
	}

	// Last month should have 30 days in leap years
	if pd.JalaliMonthLength(1403, 12) != 30 {
		t.Errorf("Month 12 in leap year 1403 should have 30 days")
	}

	// Test Gregorian month lengths
	expectedLengths := map[int]int{
		1: 31, 2: 28, 3: 31, 4: 30, 5: 31, 6: 30,
		7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31,
	}

	for month, days := range expectedLengths {
		if pd.GregorianMonthLength(2023, month) != days {
			t.Errorf("Month %d in year 2023 should have %d days", month, days)
		}
	}

	// February in leap year
	if pd.GregorianMonthLength(2024, 2) != 29 {
		t.Errorf("February in leap year 2024 should have 29 days")
	}
}

func TestNumberConversion(t *testing.T) {
	latinNumbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	persianNumbers := []string{"۰", "۱", "۲", "۳", "۴", "۵", "۶", "۷", "۸", "۹"}

	// Test Latin to Persian
	for i, latin := range latinNumbers {
		persian := persiandate.ToPersianNumbers(latin)
		if persian != persianNumbers[i] {
			t.Errorf("ToPersianNumbers(%s) = %s, expected %s", latin, persian, persianNumbers[i])
		}
	}

	// Test Persian to Latin
	for i, persian := range persianNumbers {
		latin := persiandate.ToLatinNumbers(persian)
		if latin != latinNumbers[i] {
			t.Errorf("ToLatinNumbers(%s) = %s, expected %s", persian, latin, latinNumbers[i])
		}
	}

	// Test mixed string conversion
	mixedString := "Year 1402 month 6 day 31"
	persianConverted := persiandate.ToPersianNumbers(mixedString)
	expectedPersian := "Year ۱۴۰۲ month ۶ day ۳۱"
	if persianConverted != expectedPersian {
		t.Errorf("ToPersianNumbers(%s) = %s, expected %s", mixedString, persianConverted, expectedPersian)
	}
}

func TestDateArithmetic(t *testing.T) {
	pd := persiandate.NewPersianDate("")

	baseDate := pd.Jalali(time.Date(2023, 9, 6, 0, 0, 0, 0, time.UTC)) // 1402-06-15
	t.Logf("baseDate: %v", baseDate)
	// Test adding days
	addedDate := pd.AddDaysToJalaali(baseDate, 10)
	if addedDate.Date.Year != 1402 || addedDate.Date.Month != 6 || addedDate.Date.Day != 25 {
		t.Errorf("AddDaysToJalaali(1402-06-15, 10) = %d-%02d-%02d, expected 1402-06-25",
			addedDate.Date.Year, addedDate.Date.Month, addedDate.Date.Day)
	}

	// Test adding days crossing month boundary
	addedDate = pd.AddDaysToJalaali(baseDate, 20)
	if addedDate.Date.Year != 1402 || addedDate.Date.Month != 7 || addedDate.Date.Day != 5 {
		t.Errorf("AddDaysToJalaali(1402-06-15, 20) = %d-%02d-%02d, expected 1402-07-05",
			addedDate.Date.Year, addedDate.Date.Month, addedDate.Date.Day)
	}

	// Test adding days crossing year boundary
	yearEndDate := pd.Jalali(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)) // 1402-12-25
	addedDate = pd.AddDaysToJalaali(yearEndDate, 10)
	if addedDate.Date.Year != 1403 || addedDate.Date.Month != 1 || addedDate.Date.Day != 5 {
		t.Errorf("AddDaysToJalaali(1402-12-25, 10) = %d-%02d-%02d, expected 1403-01-05",
			addedDate.Date.Year, addedDate.Date.Month, addedDate.Date.Day)
	}

	// Test subtracting days
	subtractedDate := pd.SubtractDaysFromJalaali(baseDate, 10)
	if subtractedDate.Date.Year != 1402 || subtractedDate.Date.Month != 6 || subtractedDate.Date.Day != 5 {
		t.Errorf("SubtractDaysFromJalaali(1402-06-15, 10) = %d-%02d-%02d, expected 1402-06-05",
			subtractedDate.Date.Year, subtractedDate.Date.Month, subtractedDate.Date.Day)
	}

	// Test days between dates
	startDate := pd.Jalali(time.Date(2023, 8, 23, 0, 0, 0, 0, time.UTC)) // 1402-06-01
	endDate := pd.Jalali(time.Date(2023, 9, 23, 0, 0, 0, 0, time.UTC))   // 1402-07-01
	daysBetween := pd.DifferenceJalali(startDate, endDate)
	if daysBetween != 30 {
		t.Errorf("DaysBetweenJalaaliDates(1402-06-01, 1402-07-01) = %d, expected 30", daysBetween)
	}
}

func TestDateParsing(t *testing.T) {
	pd := persiandate.NewPersianDate("")

	validDateStr := "1402-06-15"
	date, err := pd.ParseJalaaliDateString(validDateStr)
	if err != nil {
		t.Errorf("ParseJalaaliDateString(%s) returned error: %v", validDateStr, err)
	}
	if date.Date.Year != 1402 || date.Date.Month != 6 || date.Date.Day != 15 {
		t.Errorf("ParseJalaaliDateString(%s) = %d-%02d-%02d, expected 1402-06-15",
			validDateStr, date.Date.Year, date.Date.Month, date.Date.Day)
	}

	invalidDateStr := "1402-13-15" // Invalid month
	_, err = pd.ParseJalaaliDateString(invalidDateStr)
	if err == nil {
		t.Errorf("ParseJalaaliDateString(%s) should return error for invalid date", invalidDateStr)
	}

	invalidFormatStr := "1402/06/15" // Wrong separator
	_, err = pd.ParseJalaaliDateString(invalidFormatStr)
	if err == nil {
		t.Errorf("ParseJalaaliDateString(%s) should return error for invalid format", invalidFormatStr)
	}
}
