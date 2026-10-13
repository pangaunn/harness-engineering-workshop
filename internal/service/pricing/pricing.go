// Package pricing is the per-night rate and VAT domain.
package pricing

import (
	"math"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

const VATRate = 0.07

// Quote prices a stay in a room type. Friday and Saturday nights use the
// weekend rate.
func Quote(rt model.RoomType, checkIn, checkOut model.Date) model.Quote {
	q := model.Quote{RoomTypeID: rt.ID, CheckIn: checkIn, CheckOut: checkOut}
	for d := checkIn; d < checkOut; d = d.AddDays(1) {
		n := model.NightPrice{Date: d, Weekend: isWeekend(d), Rate: rt.WeekdayRate}
		if n.Weekend {
			n.Rate = rt.WeekendRate
		}
		n.VAT = roundSatang(n.Rate * VATRate)
		q.Nights = append(q.Nights, n)
		q.Subtotal += n.Rate
		q.VAT += n.VAT
	}
	q.Subtotal = roundSatang(q.Subtotal)
	q.VAT = roundSatang(q.VAT)
	q.Total = roundSatang(q.Subtotal + q.VAT)
	return q
}

func isWeekend(d model.Date) bool {
	wd := d.Time().Weekday()
	return wd == time.Friday || wd == time.Saturday
}

func roundSatang(baht float64) float64 {
	return math.Round(baht*100) / 100
}
