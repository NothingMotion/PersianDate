package persiandate_go

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
	Year  int
	Month int
	Day   int
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
func NewPersianDate(format string) *PersianDate {
	return &PersianDate{FORMAT: format, persianNumbers: PersianNumbers, latinNumbers: LatinNumbers, persianMonths: PersianMonths, persianShortMonths: PersianShortMonths, persianDays: PersianDays, persianShortDays: PersianShortDays, persianSeasons: PersianSeasons}
}

func (p *PersianDate) JalaliFull(t time.Time) PersianDateResponse {
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
func (p *PersianDate) Jalali(t time.Time) JalaliDate {
	year, month, day := t.Date()

	return p.julianDayToJalali(
		p.gregorianToJulianDay(year,
			int(month), // in case if month is 0, it will be 1
			day,
		),
	)
}

func (p *PersianDate) JalaliNow() JalaliDate {
	return p.Jalali(time.Now())
}

func (p *PersianDate) JalaliFullNow() PersianDateResponse {
	return p.JalaliFull(time.Now())
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
func (p *PersianDate) ToJalali(gy, gm, gd int) JalaliDate {
	// If gy is a time.Time object, extract the date components
	return p.julianDayToJalali(p.gregorianToJulianDay(gy, gm, gd))
}

// ToGregorian converts a Jalali date to Gregorian
func (p *PersianDate) ToGregorian(jy, jm, jd int) GregorianDate {
	return p.julianDayToGregorian(p.jalaliToJulianDay(jy, jm, jd))
}

// JalaliWeek returns Saturday and Friday day of current week (week starts on Saturday)
func (p *PersianDate) JalaliWeek(jy, jm, jd int) map[string]JalaliDate {
	// Create a time.Time object from Jalali date to get day of week
	gDate := p.ToGregorian(jy, jm, jd)
	t := time.Date(gDate.Year, time.Month(gDate.Month), gDate.Day, 0, 0, 0, 0, time.Local)

	// Get day of week (0 = Sunday, 6 = Saturday)
	dayOfWeek := int(t.Weekday())

	// Calculate difference to Saturday (start of week in Jalali calendar)
	// If it's Saturday (6), difference is 0
	// Otherwise, we need to go back (dayOfWeek + 1) days
	startDayDifference := 0
	if dayOfWeek != 6 {
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
func (p *PersianDate) JalaliToTimeObject(jy, jm, jd, h, m, s, ms int) time.Time {
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
func (p *PersianDate) FormatJalaliDate(jDate JalaliDate) string {
	format := p.FORMAT

	// Replace year
	format = strings.ReplaceAll(format, "YYYY", fmt.Sprintf("%04d", jDate.Year))

	// Replace Month
	format = strings.ReplaceAll(format, "MM", fmt.Sprintf("%02d", jDate.Month))

	// Replace day
	format = strings.ReplaceAll(format, "DD", fmt.Sprintf("%02d", jDate.Day))

	return format
}

// AddDaysToJalali adds days to a Jalali date and returns the new date
func (p *PersianDate) AddDaysToJalali(jDate JalaliDate, days int) JalaliDate {
	timeObject := p.JalaliToTimeObject(jDate.Year, jDate.Month, jDate.Day, 0, 0, 0, 0)
	timeObject = timeObject.AddDate(0, 0, days)
	return p.Jalali(timeObject)
}

// SubtractDaysFromJalali subtracts days from a Jalali date and returns the new date
func (p *PersianDate) SubtractDaysFromJalali(jDate JalaliDate, days int) JalaliDate {
	return p.AddDaysToJalali(jDate, -days)
}

// ParseJalaliDateString parses a string in format YYYY-MM-DD to a Jalali date
func (p *PersianDate) ParseJalaliDateString(dateStr string) (JalaliDate, error) {
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
func (p *PersianDate) DifferenceJalali(start, end JalaliDate) int {
	startJDN := p.jalaliToJulianDay(start.Year, start.Month, start.Day)
	endJDN := p.jalaliToJulianDay(end.Year, end.Month, end.Day)
	return endJDN - startJDN
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

	fmt.Println(pd.JalaliFullNow())

	fmt.Println(pd.GetSeason(5))

}
