package pricing

import (
	"testing"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

var (
	standard = model.RoomType{ID: "standard", WeekdayRate: 145000, WeekendRate: 185000}
	deluxe   = model.RoomType{ID: "deluxe", WeekdayRate: 214925, WeekendRate: 254925}
)

func assertMoney(t *testing.T, field string, got, want model.Money) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %s, want %s", field, got, want)
	}
}

func TestQuote_WeekdayStay(t *testing.T) {
	// Monday and Tuesday nights.
	q := Quote(standard, "2026-10-12", "2026-10-14")
	if len(q.Nights) != 2 {
		t.Fatalf("got %d nights, want 2", len(q.Nights))
	}
	assertMoney(t, "subtotal", q.Subtotal, 290000)
	assertMoney(t, "vat", q.VAT, 20300)
	assertMoney(t, "total", q.Total, 310300)
}

func TestQuote_WeekendRate(t *testing.T) {
	q := Quote(standard, "2026-10-16", "2026-10-18")
	for _, n := range q.Nights {
		if !n.Weekend {
			t.Errorf("%s: want weekend night", n.Date)
		}
	}
	assertMoney(t, "subtotal", q.Subtotal, 370000)
}

func TestQuote_VATRoundedPerNight(t *testing.T) {
	// Thursday (weekday), Friday and Saturday (weekend) nights.
	q := Quote(deluxe, "2026-10-15", "2026-10-18")
	assertMoney(t, "night 1 vat", q.Nights[0].VAT, 15045)
	assertMoney(t, "night 2 vat", q.Nights[1].VAT, 17845)
	assertMoney(t, "subtotal", q.Subtotal, 724775)
	assertMoney(t, "vat", q.VAT, 50735)
	assertMoney(t, "total", q.Total, 775510)
}
