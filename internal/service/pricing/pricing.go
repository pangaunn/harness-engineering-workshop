// Package pricing is the per-night rate and VAT domain.
package pricing

import (
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

// VATPercent is the VAT rate included in every price.
const VATPercent = 7

// SurchargePercent raises a night's rate when at least two thirds of the
// room type's rooms are already booked that night.
const SurchargePercent = 20

// Quote prices a stay in a room type. Friday and Saturday nights use the
// weekend rate. Busy nights (see occ) cost SurchargePercent more.
func Quote(rt model.RoomType, checkIn, checkOut model.Date, occ model.Occupancy) model.Quote {
	q := model.Quote{RoomTypeID: rt.ID, CheckIn: checkIn, CheckOut: checkOut}
	for d := checkIn; d < checkOut; d = d.AddDays(1) {
		n := model.NightPrice{Date: d, Weekend: isWeekend(d), Rate: rt.WeekdayRate}
		if n.Weekend {
			n.Rate = rt.WeekendRate
		}
		if busy(occ, d) {
			n.Rate = (n.Rate*(100+SurchargePercent) + 50) / 100
		}
		n.VAT = vat(n.Rate)
		q.Nights = append(q.Nights, n)
		q.Subtotal += n.Rate
		q.VAT += n.VAT
	}
	q.Total = q.Subtotal + q.VAT
	return q
}

// busy reports whether at least two thirds of the rooms are booked on d.
func busy(occ model.Occupancy, d model.Date) bool {
	return occ.Rooms > 0 && occ.Booked[d]*3 >= occ.Rooms*2
}

func isWeekend(d model.Date) bool {
	wd := d.Time().Weekday()
	return wd == time.Friday || wd == time.Saturday
}

// vat is the VAT on one night's rate, rounded half-up to the satang.
func vat(rate model.Money) model.Money {
	return (rate*VATPercent + 50) / 100
}
