package model

import (
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

// Date is a calendar date in YYYY-MM-DD form. Stays are half-open
// intervals [CheckIn, CheckOut), so plain string comparison orders dates.
type Date string

func ParseDate(s string) (Date, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return "", fmt.Errorf("invalid date %q, want YYYY-MM-DD", s)
	}
	return Date(t.Format(dateLayout)), nil
}

func DateOf(t time.Time) Date {
	return Date(t.In(HotelTime).Format(dateLayout))
}

func (d Date) Time() time.Time {
	t, _ := time.ParseInLocation(dateLayout, string(d), HotelTime)
	return t
}

func (d Date) AddDays(n int) Date {
	return Date(d.Time().AddDate(0, 0, n).Format(dateLayout))
}

func Nights(checkIn, checkOut Date) int {
	return int(checkOut.Time().Sub(checkIn.Time()).Hours() / 24)
}
