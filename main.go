package main

import "fmt"

type PersianDate struct {
	FORMAT string
}

func NewPersianDate(format string) *PersianDate {

	return &PersianDate{FORMAT: format}
}
func (p *PersianDate) Jalali() string {
	return p.FORMAT
}
func (p *PersianDate) IsLeapYear(year int) bool {
	yearsInCycle := year % 33
	if yearsInCycle == 1 || yearsInCycle == 5 || yearsInCycle == 9 || yearsInCycle == 13 || yearsInCycle == 17 || yearsInCycle == 21 || yearsInCycle == 25 || yearsInCycle == 29 || yearsInCycle == 0 {
		return true
	}
	return false
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
}
