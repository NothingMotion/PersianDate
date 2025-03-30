package persiandate_test

import (
	"testing"
	"time"

	persiandate "github.com/NothingMotion/PersianDate"
)

func TestJalaliConversion(t *testing.T) {
	pd := persiandate.New("YYYY/MM/DD")

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
		result := pd.FromTimeFull(test.gregorianDate)
		if result.Year != test.expectedYear || result.Month != test.expectedMonth || result.Day != test.expectedDay {
			t.Errorf("JalaliFull(%v) = %d-%02d-%02d, expected %d-%02d-%02d",
				test.gregorianDate, result.Year, result.Month, result.Day,
				test.expectedYear, test.expectedMonth, test.expectedDay)
		}

		jalaliDate := pd.FromTime(test.gregorianDate)
		if jalaliDate.GetYear() != test.expectedYear || jalaliDate.GetMonth() != test.expectedMonth || jalaliDate.GetDay() != test.expectedDay {
			t.Errorf("Jalali(%v) = %d-%02d-%02d, expected %d-%02d-%02d",
				test.gregorianDate, jalaliDate.GetYear(), jalaliDate.GetMonth(), jalaliDate.GetDay(),
				test.expectedYear, test.expectedMonth, test.expectedDay)
		}
	}
}

func TestGregorianConversion(t *testing.T) {
	pd := persiandate.New("YYYY/MM/DD")

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
	pd := persiandate.New("")

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
	pd := persiandate.New("")

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
	pd := persiandate.New("")

	baseDate := pd.FromTime(time.Date(2023, 9, 6, 0, 0, 0, 0, time.UTC)).Date() // 1402-06-15
	t.Logf("baseDate: %v", baseDate)
	// Test adding days
	addedDate := pd.Add(baseDate, 10).Date()
	if addedDate.Date.Year != 1402 || addedDate.Date.Month != 6 || addedDate.Date.Day != 25 {
		t.Errorf("AddDaysToJalali(1402-06-15, 10) = %d-%02d-%02d, expected 1402-06-25",
			addedDate.Date.Year, addedDate.Date.Month, addedDate.Date.Day)
	}

	// Test adding days crossing month boundary
	addedDate = pd.Add(baseDate, 20).Date()
	if addedDate.Date.Year != 1402 || addedDate.Date.Month != 7 || addedDate.Date.Day != 4 {
		t.Errorf("AddDaysToJalali(1402-06-15, 20) = %d-%02d-%02d, expected 1402-07-05",
			addedDate.Date.Year, addedDate.Date.Month, addedDate.Date.Day)
	}

	// Test adding days crossing year boundary
	yearEndDate := pd.FromTime(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)).Date() // 1402-12-25
	addedDate = pd.Add(yearEndDate, 10).Date()
	if addedDate.Date.Year != 1403 || addedDate.Date.Month != 1 || addedDate.Date.Day != 6 {
		t.Errorf("AddDaysToJalali(1402-12-25, 10) = %d-%02d-%02d, expected 1403-01-05",
			addedDate.Date.Year, addedDate.Date.Month, addedDate.Date.Day)
	}

	// Test subtracting days
	subtractedDate := pd.Sub(baseDate, 10).Date()
	if subtractedDate.Date.Year != 1402 || subtractedDate.Date.Month != 6 || subtractedDate.Date.Day != 5 {
		t.Errorf("SubtractDaysFromJalali(1402-06-15, 10) = %d-%02d-%02d, expected 1402-06-05",
			subtractedDate.Date.Year, subtractedDate.Date.Month, subtractedDate.Date.Day)
	}

	// Test days between dates
	startDate := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 6, Day: 1}} // 1402-06-01
	t.Logf("startDate: %v", startDate)
	endDate := persiandate.JalaliDate{Date: persiandate.Date{Year: 1402, Month: 7, Day: 1}} // 1402-07-01
	t.Logf("endDate: %v", endDate)
	daysBetween := pd.Difference(startDate, endDate)
	if daysBetween != 31 {
		t.Errorf("DaysBetweenJalaliDates(1402-06-01, 1402-07-01) = %d, expected 30", daysBetween)
	}
}

func TestDateParsing(t *testing.T) {
	pd := persiandate.New("")

	validDateStr := "1402-06-15"
	date, err := pd.Parse(validDateStr)
	if err != nil {
		t.Errorf("ParseJalaliDateString(%s) returned error: %v", validDateStr, err)
	}
	if date.Date.Year != 1402 || date.Date.Month != 6 || date.Date.Day != 15 {
		t.Errorf("ParseJalaliDateString(%s) = %d-%02d-%02d, expected 1402-06-15",
			validDateStr, date.Date.Year, date.Date.Month, date.Date.Day)
	}

	invalidDateStr := "1402-13-15" // Invalid month
	_, err = pd.Parse(invalidDateStr)
	if err == nil {
		t.Errorf("ParseJalaliDateString(%s) should return error for invalid date", invalidDateStr)
	}

	invalidFormatStr := "1402/06/15" // Wrong separator
	_, err = pd.Parse(invalidFormatStr)
	if err == nil {
		t.Errorf("ParseJalaliDateString(%s) should return error for invalid format", invalidFormatStr)
	}
}

func TestJalaliWeek(t *testing.T) {
	pd := persiandate.New("")

	jalaliWeek := pd.WeekYear(1404, 1, 9)
	for key, week := range jalaliWeek {
		if key == "saturday" {
			if week.Date.Year != 1404 || week.Date.Month != 1 || week.Date.Day != 9 {
				t.Errorf("JalaliWeek(1404, 1, 9) = %v %v, expected 1404-01-09", key, week)
			}
		}
		if key == "friday" {
			if week.Date.Year != 1404 || week.Date.Month != 1 || week.Date.Day != 15 {
				t.Errorf("JalaliWeek(1404, 1, 9) = %v %v, expected 1404-01-15", key, week)
			}
		}
		t.Logf("jalaliWeek: %v %v", key, week)
	}
}

func TestYearDay(t *testing.T) {
	pd := persiandate.Instance("")

	d := pd.ToJalali(2025, 3, 29).GetYearDay()
	t.Logf("Day of persian year is: %d", d)

	var yearDays int
	if pd.IsLeapYearJalali(pd.ToJalali(2025, 3, 29).GetYear()) {
		yearDays = 366
	} else {
		yearDays = 365
	}
	percentage := float64(d) / float64(yearDays) * 100
	t.Logf("Overall %f%%  days passed", percentage)
}
func TestWeekDay(t *testing.T) {
	pd := persiandate.New("")
	wd := pd.Now().GetWeekDay()
	formatted := pd.GetDayName(wd)
	t.Logf("Current day of week is: %d, %s", wd, formatted)
}

func TestSince(t *testing.T) {
	pd := persiandate.New("")
	date := pd.ToJalali(2025, 3, 29).Date()
	diff := pd.Since(date, pd.Now().Date())

	t.Logf("Total difference days since : %d", diff)
	if diff != 1 {
		t.Errorf("TestSince(2025,3,29) = %v, expected 0", diff)
	}
}

func TestUntil(t *testing.T) {
	pd := persiandate.New("")

	date := pd.Now().Date()
	remained := pd.Until(pd.ToJalali(2025, 4, 29).Date(), date)

	t.Logf("Total remained days until: %d", remained)
	if remained != 30 {
		t.Errorf("TestUntil(2025,4,29) = %v, expected 30", remained)
	}
}

func TestFormatting(t *testing.T) {
	ts := time.Now()
	formatted := ts.Format("2006 Feb")

	pd := persiandate.New("YY/MM/dd HH:ii:ss a L ff mm rr")
	// pd := persiandate.New("c")
	formatted = pd.Format(pd.Now().Date(), false)
	t.Logf(formatted)

}

func TestNewDate(t *testing.T) {
	pd := persiandate.New("YYYY/MM/DD")
	pd2 := pd.ToJalali(2025, 3, 29)

	start := pd.Date()
	days := pd.Until(pd.ToJalali(2025, 3, 30).Date(), start)

	t.Logf("remaining days: %d", days)

	if pd != pd2 {
		t.Errorf("TestNewDate() Not the same date")
	}
	t.Logf("pd: %v", pd.Date())
	t.Logf("pd2: %v", pd2.Date())
}

func TestDateInstance(t *testing.T) {
	pd := persiandate.Instance("YY/MM/dd")
	formatted := pd.ToJalali(2025, 3, 29).Add(pd.Date(), 1).Format(pd.Date())
	t.Logf(formatted)
}
func TestDateArthematic(t *testing.T) {

	pd := persiandate.New("YYYY/MM/DD")
	target := pd.ToJalali(2025, 4, 29).Date()

	diff := pd.
		ToJalali(2025, 3, 29).
		Add(pd.Date(), 1).
		Sub(pd.Date(), 1).
		Difference(pd.Date(), target)
	if diff != 31 {
		t.Errorf("expected 31 got %v", diff)
	}

	start := pd.ToJalali(2025, 3, 29).Date()
	end := pd.ToJalali(2025, 4, 29).Date()

	// Format the dates for debugging purposes if needed
	_ = pd.Format(start)
	_ = pd.Format(end)

	diff2 := pd.Difference(start, end)
	if diff2 != 31 {
		t.Errorf("expected 31 got %v", diff2)
	}

}
