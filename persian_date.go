package persiandate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PersianDate struct {
	FORMAT string

	latinNumbers       []string
	persianNumbers     []string
	persianMonths      []string
	persianShortMonths []string
	persianDays        []string
	persianShortDays   []string
	persianSeasons     []string

	currentDate JalaliDate
	targetDate  JalaliDate
}

type DateResponse struct {
	Year       int
	Month      int
	Day        int
	Hour       int
	Minute     int
	Second     int
	isLeapYear bool
}

func (d DateResponse) String() string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", d.Year, d.Month, d.Day, d.Hour, d.Minute, d.Second)
}

type PersianDateResponse struct {
	DateResponse
}

func (p PersianDateResponse) String() string {
	return p.DateResponse.String()
}

type GregorianDateResponse struct {
	DateResponse
}

func (g GregorianDateResponse) String() string {
	return g.DateResponse.String()
}

type Date struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

func (d Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", d.Year, d.Month, d.Day)
}

type GregorianDate struct {
	Date
}

func (g GregorianDate) String() string {
	return g.Date.String()
}

type JalaliDate struct {
	Date
}

func (j JalaliDate) String() string {
	return j.Date.String()
}

type jalCalReturn struct {
	leap  int
	gy    int
	march int
}

var once sync.Once
var instance *PersianDate

// Instance creates a new PersianDate object which is a singleton
func Instance(format string) *PersianDate {
	once.Do(func() {
		instance = &PersianDate{FORMAT: format, persianNumbers: PersianNumbers, latinNumbers: LatinNumbers, persianMonths: PersianMonths, persianShortMonths: PersianShortMonths, persianDays: PersianDays, persianShortDays: PersianShortDays, persianSeasons: PersianSeasons}
	})
	return instance
}

// NewPersianDate creates a new PersianDate object which is not a singleton
func New(format string) *PersianDate {
	return &PersianDate{FORMAT: format, persianNumbers: PersianNumbers, latinNumbers: LatinNumbers, persianMonths: PersianMonths, persianShortMonths: PersianShortMonths, persianDays: PersianDays, persianShortDays: PersianShortDays, persianSeasons: PersianSeasons}
}

func (p *PersianDate) FromTimeFull(t time.Time) PersianDateResponse {
	year, month, day := t.Date()
	d := p.julianDayToJalali(
		p.gregorianToJulianDay(year,
			int(month), // in case if month is 0, it will be 1
			day,
		),
	)
	response := PersianDateResponse{
		DateResponse: DateResponse{
			Year:       d.Year,
			Month:      d.Month,
			Day:        d.Day,
			Hour:       t.Hour(),
			Minute:     t.Minute(),
			Second:     t.Second(),
			isLeapYear: p.IsLeapYearJalali(d.Year),
		},
	}

	return response
}
func (p *PersianDate) FromTime(t time.Time) *PersianDate {
	year, month, day := t.Date()

	p.currentDate = p.julianDayToJalali(
		p.gregorianToJulianDay(year,
			int(month), // in case if month is 0, it will be 1
			day,
		),
	)
	return p
}

func (p *PersianDate) Now() *PersianDate {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return p.FromTime(time.Now())

	}
	return p.FromTime(time.Now().In(loc))

}

func (p *PersianDate) NowFull() PersianDateResponse {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return p.FromTimeFull(time.Now())
	}
	return p.FromTimeFull(time.Now().In(loc))
}

// Detect wheter if given persian year is leap year (kabiseh) or not
func (p *PersianDate) IsLeapYearJalali(year int) bool {
	if year <= 0 {
		year = year - 1
	}
	yearsInCycle := year % 33
	remaineders := []int{1, 5, 9, 13, 17, 22, 26, 30}
	for _, remaineder := range remaineders {
		if yearsInCycle == remaineder {
			return true
		}
	}
	return false
}

func (p *PersianDate) IsLeapYearGregorian(year int) bool {
	if year <= 0 {
		year = year - 1
	}

	// if year is divisible by 4 and not divisible by 100 or year is divisible by 400, then it is a leap year
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		return true
	}
	return false
}
func (p *PersianDate) JalaliMonthLength(jy, jm int) int {
	// if month is less than or equal to 6, then it is a 31 day month
	if jm <= 6 {
		return 31
	}
	// if month is greater than 6 and less than or equal to 11, then it is a 30 day month
	if jm > 6 && jm <= 11 {
		return 30
	}
	// if month is 12 it will be 30 days in non-leap year and 29 days in leap year
	if jm == 12 {
		if p.IsLeapYearJalali(jy) {
			return 30
		}
		return 29
	}
	return 0

}
func (p *PersianDate) GregorianMonthLength(gy, gm int) int {

	/*
		31-day months: January, March, May, July, August, October, December

		30-day months: April, June, September, November

		February: 28 days (29 in leap years)
	*/
	months := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if p.IsLeapYearGregorian(gy) {
		months[1] = 29
	}
	return months[gm-1]
}
func (p *PersianDate) jalCal(jy int, withoutLeap bool) jalCalReturn {

	breaks := []int{-61, 9, 38, 199, 426, 686, 756, 818, 1111, 1181, 1210,
		1635, 2060, 2097, 2192, 2262, 2324, 2394, 2456, 3178}
	bl := len(breaks)
	gy := jy + 621
	leapJ := -14
	jp := breaks[0]

	if jy < jp || jy >= breaks[bl-1] {
		panic(errors.New("invalid Jalali year " + fmt.Sprintf("%d", jy)))
	}

	// Find the limiting years for the Jalali year jy.
	var jump, jm, i int
	for i = 1; i < bl; i++ {
		jm = breaks[i]
		jump = jm - jp
		if jy < jm {
			break
		}
		leapJ = leapJ + p.div(jump, 33)*8 + p.div(p.mod(jump, 33), 4)
		jp = jm
	}
	n := jy - jp

	// Find the number of leap years from AD 621 to the beginning
	// of the current Jalali year in the Persian calendar.
	leapJ = leapJ + p.div(n, 33)*8 + p.div(p.mod(n, 33)+3, 4)
	if p.mod(jump, 33) == 4 && jump-n == 4 {
		leapJ += 1
	}

	// And the same in the Gregorian calendar (until the year gy).
	leapG := p.div(gy, 4) - p.div((p.div(gy, 100)+1)*3, 4) - 150

	// Determine the Gregorian date of Farvardin the 1st.
	march := 20 + leapJ - leapG

	// Return with gy and march when we don't need leap
	if withoutLeap {
		return jalCalReturn{leap: 0, gy: gy, march: march}
	}

	// Find how many years have passed since the last leap year.
	if jump-n < 6 {
		n = n - jump + p.div(jump+4, 33)*33
	}
	leap := p.mod(p.mod(n+1, 33)-1, 4)
	if leap == -1 {
		leap = 4
	}

	return jalCalReturn{leap: leap, gy: gy, march: march}

}
func (p *PersianDate) jalaliToJulianDay(jy, jm, jd int) int {
	if !p.isValidJalaliDate(JalaliDate{Date: Date{Year: jy, Month: jm, Day: jd}}) {
		panic(errors.New("invalid Jalali date"))
	}
	r := p.jalCal(jy, true)
	return p.gregorianToJulianDay(r.gy, 3, r.march) + (jm-1)*31 - p.div(jm, 7)*(jm-7) + jd - 1
}

func (p *PersianDate) gregorianToJulianDay(gy, gm, gd int) int {
	//	days := ((((gm - 8) / 6) + 100100) * 1461) / 4 + ()

	d := p.div((gy+p.div(gm-8, 6)+100100)*1461, 4) +
		p.div(153*p.mod(gm+9, 12)+2, 5) +
		gd - 34840408
	d = d - p.div(p.div(gy+100100+p.div(gm-8, 6), 100)*3, 4) + 752
	return d

}
func (p *PersianDate) julianDayToJalali(jdn int) JalaliDate {
	gy := p.julianDayToGregorian(jdn).Year // Calculate Gregorian year (gy).
	jy := gy - 621
	r := p.jalCal(jy, false)
	jdn1f := p.gregorianToJulianDay(gy, 3, r.march)
	var jd, jm int
	k := jdn - jdn1f

	if k >= 0 {
		if k <= 185 {
			// The first 6 months.
			jm = 1 + p.div(k, 31)
			jd = p.mod(k, 31) + 1
			return JalaliDate{Date: Date{Year: jy, Month: jm, Day: jd}}
		} else {
			// The remaining months.
			k -= 186
		}
	} else {
		// Previous Jalali year.
		jy -= 1
		k += 179
		if r.leap == 1 {
			k += 1
		}
	}
	jm = 7 + p.div(k, 30)
	jd = p.mod(k, 30) + 1
	return JalaliDate{Date: Date{Year: jy, Month: jm, Day: jd}}
}

func (p *PersianDate) julianDayToGregorian(jdn int) GregorianDate {

	j := 4*jdn + 139361631
	j = j + p.div(p.div(4*jdn+183187720, 146097)*3, 4)*4 - 3908
	i := p.div(p.mod(j, 1461), 4)*5 + 308
	gd := p.div(p.mod(i, 153), 5) + 1
	gm := p.mod(p.div(i, 153), 12) + 1
	gy := p.div(j, 1461) - 100100 + p.div(8-gm, 6)

	return GregorianDate{Date: Date{Year: gy, Month: gm, Day: gd}}

}
func (p *PersianDate) isValidJalaliDate(date JalaliDate) bool {
	if p.isDateEmpty(date.Date) {
		return false
	}
	validMonth := date.Month >= 1 && date.Month <= 12
	validDay := date.Day >= 1 && date.Day <= p.JalaliMonthLength(date.Year, date.Month)
	validYear := date.Year >= 0 && date.Year <= 3778
	return validMonth && validDay && validYear
}
func (p *PersianDate) isDateEmpty(date Date) bool {
	return date.Year == 0 && date.Month == 0 && date.Day == 0
}
func (p *PersianDate) div(a, b int) int {
	return (a / b)
}
func (p *PersianDate) mod(a, b int) int {
	return a % b
}

// Converts Latin number like 1402 to Persian numbers like ۱۴۰۲
func ToPersianNumbers(text string) string {
	for i, value := range LatinNumbers {
		text = strings.ReplaceAll(text, value, PersianNumbers[i])
	}
	return text
}

// Converts Persian numbers to Latin numbers like ۱۴۰۲ to 1402
func ToLatinNumbers(text string) string {
	for i, value := range PersianNumbers {
		text = strings.ReplaceAll(text, value, LatinNumbers[i])
	}
	return text
}

// ToJalali converts a Gregorian date to Jalali
func (p *PersianDate) ToJalali(gy, gm, gd int) *PersianDate {

	p.currentDate = p.julianDayToJalali(p.gregorianToJulianDay(gy, gm, gd))
	// If gy is a time.Time object, extract the date components

	return p
}

// ToGregorian converts a Jalali date to Gregorian
func (p *PersianDate) ToGregorian(jy, jm, jd int) GregorianDate {
	return p.julianDayToGregorian(p.jalaliToJulianDay(jy, jm, jd))
}

// JalaliWeek returns Saturday and Friday day of current week (week starts on Saturday)
func (p *PersianDate) WeekYear(jy, jm, jd int) map[string]JalaliDate {
	// Get day of week (0 = Saturday, 6 = Friday) based on jalali date
	p.currentDate = JalaliDate{Date{Year: jy, Month: jm, Day: jd}}
	dayOfWeek := p.GetWeekDay()

	// Calculate difference to Saturday (start of week in Jalali calendar)
	// If it's Saturday (0), difference is 0
	// Otherwise, we need to go back (dayOfWeek + 1) days
	startDayDifference := 0
	if dayOfWeek != 0 {
		startDayDifference = -(dayOfWeek + 1)
	}
	endDayDifference := 6 + startDayDifference

	// Get Julian day number for the current date
	jdn := p.jalaliToJulianDay(jy, jm, jd)

	// Calculate Saturday and Friday of the week
	saturdayDate := p.julianDayToJalali(jdn + startDayDifference)
	fridayDate := p.julianDayToJalali(jdn + endDayDifference)

	return map[string]JalaliDate{
		"saturday": saturdayDate,
		"friday":   fridayDate,
	}
}

// JalaliToTimeObject converts Jalali calendar dates to time.Time object
func (p *PersianDate) ToTime(jy, jm, jd, h, m, s, ms int) time.Time {
	GregorianDate := p.ToGregorian(jy, jm, jd)

	return time.Date(
		GregorianDate.Year,
		time.Month(GregorianDate.Month),
		GregorianDate.Day,
		h, m, s, ms*1000000, // ms to nanoseconds
		time.Local,
	)
}

// FormatJalaliDate formats a Jalali date according to the format string
func (p *PersianDate) Format(jDate JalaliDate, toPersian ...interface{}) string {
	format := p.FORMAT

	// Use the time components directly from JalaliDate
	hour := jDate.Hour
	minute := jDate.Minute
	second := jDate.Second

	var convertNumbers bool

	if len(toPersian) != 0 {
		switch toPersian[0].(type) {
		case bool:
			if toPersian[0] == true {
				convertNumbers = true
			}
		}
	}

	// AM/PM values
	var shortAMPM, longAMPM string
	if hour < 12 {
		shortAMPM = "ق.ظ"
		longAMPM = "قبل از ظهر"
	} else {
		shortAMPM = "ب.ظ"
		longAMPM = "بعد از ظهر"
	}

	// Leap year text
	var leapYearText string
	if p.IsLeapYearJalali(jDate.Year) {
		leapYearText = "بله"
	} else {
		leapYearText = "خیر"
	}

	// Year formatting options
	yearStr := fmt.Sprintf("%d", jDate.Year)
	yearLen := len(yearStr)

	// Create a map of replacements for better performance
	replacements := map[string]string{
		// Year formats - enhanced
		"YYYY": fmt.Sprintf("%04d", jDate.Year),     // Full year with leading zeros (4 digits)
		"YYY":  yearStr[max(0, yearLen-3):],         // 3-digit year
		"YY":   fmt.Sprintf("%02d", jDate.Year%100), // 2-digit year with leading zero
		"Y":    fmt.Sprintf("%d", jDate.Year%100),   // 2-digit year without leading zero
		"y":    yearStr,                             // Full year without padding

		// Month formats
		"MM": fmt.Sprintf("%02d", jDate.Month), // Month number with leading zero
		"M":  fmt.Sprintf("%d", jDate.Month),   // Month number without leading zero
		"mm": p.GetMonthName(jDate.Month),      // Full month name
		"km": p.GetShortMonthName(jDate.Month), // Short month name
		"mb": p.GetMonthSymbol(jDate.Month),    // Month symbol

		// Day formats
		"DD": fmt.Sprintf("%02d", jDate.Day),                              // Day with leading zero
		"D":  fmt.Sprintf("%d", jDate.Day),                                // Day without leading zero
		"dd": fmt.Sprintf("%02d", jDate.Day),                              // Alternative day with leading zero (for compatibility)
		"d":  fmt.Sprintf("%d", jDate.Day),                                // Alternative day without leading zero (for compatibility)
		"rr": PersianMonthDays[min(jDate.Day-1, len(PersianMonthDays)-1)], // Day in Persian words (with bounds check)

		// Weekday formats
		"l":  p.GetDayName(p.GetWeekDay()),      // Full day name
		"rh": p.GetDayName(p.GetWeekDay()),      // Full day name (alias)
		"kh": p.GetShortDayName(p.GetWeekDay()), // Short day name

		// Time formats
		"HH": fmt.Sprintf("%02d", hour),           // 24-hour with leading zero
		"H":  fmt.Sprintf("%d", hour),             // 24-hour without leading zero
		"hh": fmt.Sprintf("%02d", hourTo12(hour)), // 12-hour with leading zero
		"h":  fmt.Sprintf("%d", hourTo12(hour)),   // 12-hour without leading zero
		"ii": fmt.Sprintf("%02d", minute),         // Minutes with leading zero
		"i":  fmt.Sprintf("%d", minute),           // Minutes without leading zero
		"ss": fmt.Sprintf("%02d", second),         // Seconds with leading zero
		"s":  fmt.Sprintf("%d", second),           // Seconds without leading zero

		// AM/PM
		"a": shortAMPM, // Persian AM/PM abbreviated
		"A": longAMPM,  // Persian AM/PM full

		// Other formats
		"L":  leapYearText,                                                // Is leap year
		"b":  fmt.Sprintf("%d", int(float64(jDate.Month)/float64(3.1)+1)), // Season number
		"ff": p.GetSeason(jDate.Month),                                    // Season name
	}

	// Full date-time format in Persian style
	replacements["c"] = fmt.Sprintf("%d/%d/%d ،%d:%d:%d %s",
		jDate.Year, jDate.Month, jDate.Day,
		hour, minute, second,
		p.GetDayName(p.GetWeekDay()))

	// Apply all replacements (using a custom sort to avoid partial replacements)
	orderedPatterns := []string{"YYYY", "YYY", "YY", "Y", "y", "MM", "M", "mm", "km", "mb",
		"DD", "D", "dd", "d", "rr", "l", "rh", "kh", "HH", "H", "hh", "h", "ii", "i", "ss", "s",
		"a", "A", "L", "b", "ff", "c"}

	for _, pattern := range orderedPatterns {
		if replacement, exists := replacements[pattern]; exists {
			format = strings.ReplaceAll(format, pattern, replacement)
		}
	}

	if convertNumbers {
		format = ToPersianNumbers(format)
	}
	return format
}

// Helper function to convert 24-hour format to 12-hour format
func hourTo12(hour int) int {
	if hour == 0 {
		return 12
	}
	if hour > 12 {
		return hour - 12
	}
	return hour
}

// Helper function to return the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to return the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// AddDaysToJalali adds days to a Jalali date and returns the new date
func (p *PersianDate) Add(jDate JalaliDate, days int) *PersianDate {
	timeObject := p.ToTime(jDate.Year, jDate.Month, jDate.Day, 0, 0, 0, 0)
	timeObject = timeObject.AddDate(0, 0, days)
	return p.FromTime(timeObject)

}
func (p *PersianDate) AddDate(jDate JalaliDate, y, m, d int) *PersianDate {
	timeObject := p.ToTime(jDate.Year, jDate.Month, jDate.Day, 0, 0, 0, 0)
	timeObject = timeObject.AddDate(y, m, d)
	return p.FromTime(timeObject)
}

// SubtractDaysFromJalali subtracts days from a Jalali date and returns the new date
func (p *PersianDate) Sub(jDate JalaliDate, days int) *PersianDate {
	p.Add(jDate, -days)
	return p
}

// Returns oldest date between dates
func (p *PersianDate) Min(dates ...interface{}) JalaliDate {
	// Handle a slice passed as single argument
	if len(dates) == 1 {
		switch d := dates[0].(type) {
		case []JalaliDate:
			return p.processMinMax(d, true)
		case []interface{}:
			// Convert to sorted dates and return min
			sortedDates := p.Sort(d...)
			if len(sortedDates) > 0 {
				return sortedDates[0]
			}
			return JalaliDate{}
		}
	}

	// Process regular variadic arguments
	sortedDates := p.Sort(dates...)
	if len(sortedDates) < 1 {
		return JalaliDate{}
	}
	return sortedDates[0]
}

// Return newest date between dates
func (p *PersianDate) Max(dates ...interface{}) JalaliDate {
	// Handle a slice passed as single argument
	if len(dates) == 1 {
		switch d := dates[0].(type) {
		case []JalaliDate:
			return p.processMinMax(d, false)
		case []interface{}:
			// Convert to sorted dates and return max
			sortedDates := p.Sort(d...)
			if len(sortedDates) > 0 {
				return sortedDates[len(sortedDates)-1]
			}
			return JalaliDate{}
		}
	}

	// Process regular variadic arguments
	sortedDates := p.Sort(dates...)
	if len(sortedDates) < 1 {
		return JalaliDate{}
	}
	return sortedDates[len(sortedDates)-1]
}

// Helper function to process min/max from a slice of JalaliDate
func (p *PersianDate) processMinMax(dates []JalaliDate, isMin bool) JalaliDate {
	if len(dates) == 0 {
		return JalaliDate{}
	}

	if len(dates) == 1 {
		return dates[0]
	}

	index := 0
	if !isMin {
		index = len(dates) - 1
	}

	// Sort by date
	sortedDates := make([]JalaliDate, len(dates))
	copy(sortedDates, dates)

	// Sort the dates by converting to Julian days
	for i := 0; i < len(sortedDates); i++ {
		for j := i + 1; j < len(sortedDates); j++ {
			jdni := p.jalaliToJulianDay(sortedDates[i].Year, sortedDates[i].Month, sortedDates[i].Day)
			jdnj := p.jalaliToJulianDay(sortedDates[j].Year, sortedDates[j].Month, sortedDates[j].Day)

			if jdni > jdnj {
				sortedDates[i], sortedDates[j] = sortedDates[j], sortedDates[i]
			}
		}
	}

	return sortedDates[index]
}

// ParseJalaliDateString parses a string in format YYYY-MM-DD to a Jalali date
func (p *PersianDate) Parse(dateStr string) (JalaliDate, error) {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return JalaliDate{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return JalaliDate{}, errors.New("invalid year format")
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return JalaliDate{}, errors.New("invalid month format")
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return JalaliDate{}, errors.New("invalid day format")
	}

	if !p.isValidJalaliDate(JalaliDate{Date: Date{Year: year, Month: month, Day: day}}) {
		return JalaliDate{}, errors.New("invalid date values")
	}

	return JalaliDate{Date: Date{Year: year, Month: month, Day: day}}, nil
}

// DaysBetweenJalaliDates calculates the number of days between two Jalali dates
func (p *PersianDate) Difference(start, end JalaliDate) int {
	startJDN := p.jalaliToJulianDay(start.Year, start.Month, start.Day)
	endJDN := p.jalaliToJulianDay(end.Year, end.Month, end.Day)
	return endJDN - startJDN
}

// Until calculates days until the end date
// If no date is provided, it uses the current date as the start date
func (p *PersianDate) Until(end JalaliDate, startOpt ...JalaliDate) int {
	var start JalaliDate
	if len(startOpt) > 0 {
		start = startOpt[0]
	} else {
		start = p.currentDate
	}
	return p.Difference(start, end)
}

// Since calculates days since the start date
// If no date is provided, it uses the current date as the end date
func (p *PersianDate) Since(start JalaliDate, endOpt ...JalaliDate) int {
	var end JalaliDate
	if len(endOpt) > 0 {
		end = endOpt[0]
	} else {
		end = p.currentDate
	}
	return p.Difference(start, end)
}

func (p *PersianDate) Equal(a, b JalaliDate) bool {
	return a.Year == b.Year && a.Month == b.Month && a.Day == b.Day
}

func (p *PersianDate) Sort(dates ...interface{}) []JalaliDate {
	if len(dates) < 1 {
		fmt.Println("dates to sort are less than 1")
		return []JalaliDate{}
	}
	var jalaliDates []JalaliDate

	// Process each date in the input
	for _, date := range dates {
		switch d := date.(type) {
		case JalaliDate:
			jalaliDates = append(jalaliDates, d)
		case Date:
			jalaliDates = append(jalaliDates, JalaliDate{Date: d})
		case string:
			parsed, err := p.Parse(d)
			if err == nil {
				jalaliDates = append(jalaliDates, parsed)
			}
		case time.Time:
			response := p.FromTimeFull(d)
			jalaliDate := JalaliDate{
				Date: Date{
					Year:   response.Year,
					Month:  response.Month,
					Day:    response.Day,
					Hour:   response.Hour,
					Minute: response.Minute,
					Second: response.Second,
				},
			}
			jalaliDates = append(jalaliDates, jalaliDate)
		}
	}

	// Sort the dates by converting to Julian days
	for i := 0; i < len(jalaliDates); i++ {
		for j := i + 1; j < len(jalaliDates); j++ {
			jdni := p.jalaliToJulianDay(jalaliDates[i].Year, jalaliDates[i].Month, jalaliDates[i].Day)
			jdnj := p.jalaliToJulianDay(jalaliDates[j].Year, jalaliDates[j].Month, jalaliDates[j].Day)

			if jdni > jdnj {
				jalaliDates[i], jalaliDates[j] = jalaliDates[j], jalaliDates[i]
			}
		}
	}

	return jalaliDates
}

func (p *PersianDate) Filter(comparator func(JalaliDate) bool, dates ...interface{}) []JalaliDate {
	if len(dates) < 1 {
		return []JalaliDate{}
	}

	// Sort the dates
	sortedDates := p.Sort(dates...)

	// Filter the dates based on the comparator function
	var filteredDates []JalaliDate
	for _, date := range sortedDates {
		if comparator(date) {
			filteredDates = append(filteredDates, date)
		}
	}

	return filteredDates
}

func (p *PersianDate) GetWeekDay() int {
	jDate := p.currentDate
	t := p.ToTime(jDate.Year, jDate.Month, jDate.Day, 0, 0, 0, 0)
	return int((t.Weekday() + 1) % 7) // conversion to jalali days (saturday from 6 to 0 , and friday to 6)
}
func (p *PersianDate) GetYearDay() int {
	jDate := p.currentDate
	year := jDate.Year
	month := int(jDate.Month)
	day := jDate.Day

	dayOfYear := day

	if month < 1 || month > 12 {

		return 0
	}
	for i := 1; i < month; i++ {
		if i <= 6 {
			dayOfYear += 31
		}
		if i > 6 && i <= 11 {
			dayOfYear += 30
		}
		if i == 12 {
			if p.IsLeapYearJalali(year) {
				dayOfYear += 30
			}
			dayOfYear += 29
		}
	}
	return dayOfYear
}
func (p *PersianDate) GetYear() int {

	return p.Date().Year
}

func (p *PersianDate) GetMonth() int {

	return p.Date().Month
}

func (p *PersianDate) GetDay() int {

	return p.Date().Day
}

func (p *PersianDate) GetHour() int {
	return p.Date().Hour
}

func (p *PersianDate) GetMinute() int {
	return p.Date().Minute
}

func (p *PersianDate) GetSecond() int {
	return p.Date().Second
}

func (p *PersianDate) Clock() (int, int, int) {

	return p.GetHour(), p.GetMinute(), p.GetSecond()
}
func (p *PersianDate) Date() JalaliDate {

	return p.currentDate
}

func (p *PersianDate) Time() (int, int, int) {
	return p.GetHour(), p.GetMinute(), p.GetSecond()
}

// GetMonthName returns the Persian name of the Month
func (p *PersianDate) GetMonthName(month int) string {
	if month < 1 || month > 12 {
		return ""
	}
	return p.persianMonths[month-1]
}

// GetShortMonthName returns the short Persian name of the month
func (p *PersianDate) GetShortMonthName(month int) string {
	if month < 1 || month > 12 {
		return ""
	}
	return p.persianShortMonths[month-1]
}

func (p *PersianDate) GetMonthSymbol(month int) string {
	if month < 1 || month > 12 {
		return ""
	}
	return ""
}

// GetDayName returns the Persian name of the day of week
func (p *PersianDate) GetDayName(dayOfWeek int) string {
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return ""
	}
	return p.persianDays[dayOfWeek]
}

// GetShortDayName returns the short Persian name of the day of week
func (p *PersianDate) GetShortDayName(dayOfWeek int) string {
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return ""
	}
	return p.persianShortDays[dayOfWeek]
}

func (p *PersianDate) GetSeason(month int) string {
	if month < 1 || month > 12 {
		return ""
	}
	// Spring
	if month >= 1 && month <= 3 {
		return p.persianSeasons[0]
	}
	// Summer
	if month >= 4 && month <= 6 {
		return p.persianSeasons[1]
	}
	// Autumn
	if month >= 7 && month <= 9 {
		return p.persianSeasons[2]
	}
	// Winter
	if month >= 10 && month <= 12 {
		return p.persianSeasons[3]
	}
	return p.persianSeasons[3]
}

// Example function showing how to use the package
func main() {
	pd := Instance("YYYY/MM/DD")

	fmt.Println(pd.NowFull())

	fmt.Println(pd.GetSeason(5))

}
