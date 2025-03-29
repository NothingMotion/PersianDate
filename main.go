package main

import (
	"fmt"
	"strings"
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

func (p *PersianDate) Jalali() string {
	return p.FORMAT
}

// Detect wheter if given persian year is leap year (kabiseh) or not
func (p *PersianDate) IsLeapYear(year int) bool {
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
		if p.IsLeapYear(jy) {
			return 30
		}
		return 29
	}
	return 31

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

	fmt.Println("-61", pd.IsLeapYear(-61))
	fmt.Println("3177", pd.IsLeapYear(3177))
	fmt.Println("1402", pd.IsLeapYear(1402))
	fmt.Println("1403", pd.IsLeapYear(1403))
	fmt.Println("1404", pd.IsLeapYear(1404))
	fmt.Println("1405", pd.IsLeapYear(1405))
	fmt.Println("1406", pd.IsLeapYear(1406))
	fmt.Println("1407", pd.IsLeapYear(1407))
	fmt.Println("1408", pd.IsLeapYear(1408))
	fmt.Println("1409", pd.IsLeapYear(1409))
	fmt.Println("1410", pd.IsLeapYear(1410))

	fmt.Println(pd.ToPersianNumbers("1402/01/01 salam"))
	fmt.Println(pd.ToLatinNumbers("۱۴۰۲/۰۱/۰۱ سلام"))
}
