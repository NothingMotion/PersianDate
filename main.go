package main

import (
	"fmt"
	"strings"
)

type PersianDate struct {
	FORMAT             string
	numbersMap         map[string]string
	persianMonths      []string
	persianShortMonths []string
	persianDays        []string
	persianShortDays   []string
	persianSeasons     []string
}

func NewPersianDate(format string) *PersianDate {
	numbersMap := map[string]string{
		"0": "۰",
		"1": "۱",
		"2": "۲",
		"3": "۳",
		"4": "۴",
		"5": "۵",
		"6": "۶",
		"7": "۷",
		"8": "۸",
		"9": "۹",
		".": ".",
	}

	persianMonths := []string{"فروردین", "اردیبهشت", "خرداد", "تیر", "مرداد", "شهریور", "مهر", "آبان", "آذر", "دی", "بهمن", "اسفند"}

	persianShortMonths := []string{"فر", "ار", "خر", "تی‍", "مر", "شه‍", "مه‍", "آب‍", "آذ", "دی", "به‍", "اس‍"}

	persianDays := []string{"یکشنبه", "دوشنبه", "سه شنبه", "چهارشنبه", "پنج شنبه", "جمعه", "شنبه"}

	persianShortDays := []string{"ی", "د", "س", "چ", "پ", "ج", "ش"}

	persianSeasons := []string{"بهار", "تابستان", "پاییز", "زمستان"}

	return &PersianDate{FORMAT: format, numbersMap: numbersMap, persianMonths: persianMonths, persianShortMonths: persianShortMonths, persianDays: persianDays, persianShortDays: persianShortDays, persianSeasons: persianSeasons}
}
func (p *PersianDate) Jalali() string {
	return p.FORMAT
}
func (p *PersianDate) IsLeapYear(year int) bool {
	yearsInCycle := year % 33
	remaineders := []int{1, 5, 9, 13, 17, 21, 25, 29, 0}
	for _, remaineder := range remaineders {
		if yearsInCycle == remaineder {
			return true
		}
	}
	return false
}
func (p *PersianDate) ReplaceNumbers(text string) string {
	for key, value := range p.numbersMap {
		text = strings.ReplaceAll(text, key, value)
	}
	return text
}
func main() {
	persianDate := NewPersianDate("%d/%m/%Y")

	fmt.Println("1402", persianDate.IsLeapYear(1402))
	fmt.Println("1403", persianDate.IsLeapYear(1403))
	fmt.Println("1404", persianDate.IsLeapYear(1404))
	fmt.Println("1405", persianDate.IsLeapYear(1405))
	fmt.Println("1406", persianDate.IsLeapYear(1406))
	fmt.Println("1407", persianDate.IsLeapYear(1407))
	fmt.Println("1408", persianDate.IsLeapYear(1408))
	fmt.Println("1409", persianDate.IsLeapYear(1409))
	fmt.Println("1410", persianDate.IsLeapYear(1410))

	fmt.Println(persianDate.ReplaceNumbers("1402/01/01"))
}
