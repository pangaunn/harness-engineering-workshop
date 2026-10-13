package model

import "fmt"

// Money is an amount in satang (1 THB = 100 satang).
type Money int64

func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign = "-"
		m = -m
	}
	return fmt.Sprintf("%s%d.%02d", sign, m/100, m%100)
}

// MarshalJSON writes Money as a baht number, e.g. 1234.50.
func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}
