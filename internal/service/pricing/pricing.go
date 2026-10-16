// Package pricing is the per-night rate and VAT domain.
package pricing

import (
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

// VATPercent is the VAT rate included in every price.
const VATPercent = 7

// Quote prices a stay in a room type. Friday and Saturday nights use the
// weekend rate.
func Quote(rt model.RoomType, checkIn, checkOut model.Date) model.Quote {
	q := model.Quote{RoomTypeID: rt.ID, CheckIn: checkIn, CheckOut: checkOut}
	for d := checkIn; d < checkOut; d = d.AddDays(1) {
		n := model.NightPrice{Date: d, Weekend: isWeekend(d), Rate: rt.WeekdayRate}
		if n.Weekend {
			n.Rate = rt.WeekendRate
		}
		n.VAT = vat(n.Rate)
		q.Nights = append(q.Nights, n)
		q.Subtotal += n.Rate
		q.VAT += n.VAT
	}
	q.Total = q.Subtotal + q.VAT
	return q
}

func isWeekend(d model.Date) bool {
	wd := d.Time().Weekday()
	return wd == time.Friday || wd == time.Saturday
}

// vat is the VAT on one night's rate, rounded half-up to the satang.
func vat(rate model.Money) model.Money {
	return (rate*VATPercent + 50) / 100
}
