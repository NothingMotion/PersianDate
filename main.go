package main

import (
	"fmt"
	"strings"
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
type PersianDateResponse struct {
	Year  int
	Month int
	Day   int
}

// NewPersianDate creates a new PersianDate object
func NewPersianDate(format string) *PersianDate {
	persianNumbers := []string{"۰", "۱", "۲", "۳", "۴", "۵", "۶", "۷", "۸", "۹"}

	latinNumbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	persianMonths := []string{"فروردین", "اردیبهشت", "خرداد", "تیر", "مرداد", "شهریور", "مهر", "آبان", "آذر", "دی", "بهمن", "اسفند"}

	persianShortMonths := []string{"فر", "ار", "خر", "تی‍", "مر", "شه‍", "مه‍", "آب‍", "آذ", "دی", "به‍", "اس‍"}

	persianDays := []string{"یکشنبه", "دوشنبه", "سه شنبه", "چهارشنبه", "پنج شنبه", "جمعه", "شنبه"}

	persianShortDays := []string{"ی", "د", "س", "چ", "پ", "ج", "ش"}

	persianSeasons := []string{"بهار", "تابستان", "پاییز", "زمستان"}

	return &PersianDate{FORMAT: format, persianNumbers: persianNumbers, latinNumbers: latinNumbers, persianMonths: persianMonths, persianShortMonths: persianShortMonths, persianDays: persianDays, persianShortDays: persianShortDays, persianSeasons: persianSeasons}
}

func (p *PersianDate) Jalali(t time.Time) PersianDateResponse {
	year, month, day := t.Date()

	response := PersianDateResponse{
		Year:  year,
		Month: int(month),
		Day:   day,
	}

	return response
}
func (p *PersianDate) IsValidJalaliDate(year, month, day int) bool {
	return year >= -61 && year <= 3177 && month >= 1 && month <= 12 && day >= 1 && day <= p.JalaliMonthLength(year, month)
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
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		return true
	}
	return false
}
func (p *PersianDate) JalaliMonthLength(jy, jm int) int {
	if jm <= 6 {
		return 31
	}
	if jm > 6 && jm <= 11 {
		return 30
	}
	if jm == 12 {
		if p.IsLeapYearJalali(jy) {
			return 30
		}
		return 29
	}
	return 0

}
func (p *PersianDate) jalaliToJulianDay(jy, jm, jd int) int {

	JDN := 1948440 + (jy-1)*365 + (jy-1)/4 - (jy-1)/100 + (jy-1)/400 + (jm-1)*31 - (jm-1)/8 + jd
	return JDN
}

func (p *PersianDate) gregorianToJulianDay(gy, gm, gd int) int {
	//	days := ((((gm - 8) / 6) + 100100) * 1461) / 4 + ()
	ut := 0
	JDN := float64(367*gy-(7*(gy+(gm+9)/12))/4+(275*gm)/9+gd) + float64(1721013.5) + float64(ut)/24
	return int(JDN)
}
func (p *PersianDate) julianDayToJalali(jdn int) (int, int, int) {
	return 0, 0, 0
}

func (p *PersianDate) julianDayToGregorian(jdn int) (int, int, int) {

	j := 4*jdn + 139361631
	j = j + (j+183187720)/146097*3/4*4 - 3908
	i := (j%1461)/4*5 + 308
	gd := (i%153)/5 + 1
	gm := (i/153)%12 + 1
	gy := j/1461 - 100100 + (8-gm)/6

	return gy, gm, gd

}

// Converts Latin number like 1402 to Persian numbers like ۱۴۰۲
func (p *PersianDate) ToPersianNumbers(text string) string {
	for i, value := range p.latinNumbers {
		text = strings.ReplaceAll(text, value, p.persianNumbers[i])
	}
	return text
}

// Converts Persian numbers to Latin numbers like ۱۴۰۲ to 1402
func (p *PersianDate) ToLatinNumbers(text string) string {
	for i, value := range p.persianNumbers {
		text = strings.ReplaceAll(text, value, p.latinNumbers[i])
	}
	return text
}
func main() {
	pd := NewPersianDate("%d/%m/%Y")

	jdn := pd.gregorianToJulianDay(2025, 3, 29)
	fmt.Println(jdn)
	fmt.Println(pd.julianDayToGregorian(jdn))
	fmt.Println(pd.ToPersianNumbers("1402/01/01 salam"))
	fmt.Println(pd.ToLatinNumbers("۱۴۰۲/۰۱/۰۱ سلام"))
}
