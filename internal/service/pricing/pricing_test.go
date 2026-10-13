package pricing

import (
	"math"
	"testing"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

var (
	standard = model.RoomType{ID: "standard", WeekdayRate: 1450.00, WeekendRate: 1850.00}
	deluxe   = model.RoomType{ID: "deluxe", WeekdayRate: 2149.25, WeekendRate: 2549.25}
)

func assertBaht(t *testing.T, field string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.001 {
		t.Errorf("%s: got %.2f, want %.2f", field, got, want)
	}
}

func TestQuote_WeekdayStay(t *testing.T) {
	// Monday and Tuesday nights.
	q := Quote(standard, "2026-10-12", "2026-10-14")
	if len(q.Nights) != 2 {
		t.Fatalf("got %d nights, want 2", len(q.Nights))
	}
	assertBaht(t, "subtotal", q.Subtotal, 2900.00)
	assertBaht(t, "vat", q.VAT, 203.00)
	assertBaht(t, "total", q.Total, 3103.00)
}

func TestQuote_WeekendRate(t *testing.T) {
	q := Quote(standard, "2026-10-16", "2026-10-18")
	for _, n := range q.Nights {
		if !n.Weekend {
			t.Errorf("%s: want weekend night", n.Date)
		}
	}
	assertBaht(t, "subtotal", q.Subtotal, 3700.00)
}

func TestQuote_VATRoundedPerNight(t *testing.T) {
	// Thursday (weekday), Friday and Saturday (weekend) nights.
	q := Quote(deluxe, "2026-10-15", "2026-10-18")
	assertBaht(t, "night 1 vat", q.Nights[0].VAT, 150.45)
	assertBaht(t, "night 2 vat", q.Nights[1].VAT, 178.45)
	assertBaht(t, "subtotal", q.Subtotal, 7247.75)
	assertBaht(t, "vat", q.VAT, 507.35)
	assertBaht(t, "total", q.Total, 7755.10)
}
