package simple

import (
	"fmt"
	"time"

	persiandate "github.com/NothingMotion/PersianDate"
)

func main() {
	// Initialize Persian date
	pd := persiandate.New("YYYY/MM/DD")

	////////////////////////////////////////////

	// Gets current date to jalali
	now := pd.Now().Date()         // returns jalali date object
	nowFormatted := pd.Format(now) // formats it with provided format

	fmt.Printf("Current date: %s %s", now, nowFormatted)

	////////////////////////////////////////////

	// Gets jalali date from specific date
	from := pd.FromTime(time.Date(2025, 3, 29, 0, 0, 0, 0, nil)).Date()
	fromFormatted := pd.Format(from) // 1404-01-10

	fmt.Printf("From Time formatted: %s %s", from, fromFormatted)

	////////////////////////////////////////////

	// Gets jalali date from specific date (with timestamps)
	ts := pd.FromTime(time.Unix(1234, 0)).Date()
	fmt.Printf("From Timestamp: %s", ts)

	////////////////////////////////////////////

	// Converting jalali date object to time object
	pd.Now().ToTime(pd.GetYear(), pd.GetMonth(), pd.GetDay(), 0, 0, 0, 0)

	////////////////////////////////////////////

	// Gets day index of week from 0 (Saturday) to 6 (Friday)
	wd := pd.GetWeekDay()
	// Gets name of the day
	pd.GetDayName(wd)
	////////////////////////////////////////////

	// Gets day of year
	pd.GetYearDay()

	////////////////////////////////////////////
}
