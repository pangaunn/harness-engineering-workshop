//go:build integration

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/jobs"
	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

// TestRevenueReport_EndToEnd books and cancels stays over two days, runs the
// nightly job, and reads the report back the way a client would.
func TestRevenueReport_EndToEnd(t *testing.T) {
	s := newTestServer(t)
	job := &jobs.RevenueReport{Payments: s.payments, Reports: s.reports}

	book := func(roomType, checkIn, checkOut string) int64 {
		t.Helper()
		var b struct {
			ID int64 `json:"id"`
		}
		body := map[string]string{"room_type_id": roomType, "guest_name": "Somchai", "check_in": checkIn, "check_out": checkOut}
		if code := s.do(t, "POST", "/bookings", body, &b); code != http.StatusCreated {
			t.Fatalf("book %s %s: status %d", roomType, checkIn, code)
		}
		return b.ID
	}
	cancel := func(id int64) {
		t.Helper()
		if code := s.do(t, "POST", fmt.Sprintf("/bookings/%d/cancel", id), nil, nil); code != http.StatusOK {
			t.Fatalf("cancel %d: status %d", id, code)
		}
	}

	// Day 1: three charges, one refund.
	s.clock.t = time.Date(2026, 10, 14, 9, 30, 0, 0, model.HotelTime)
	standard := book("standard", "2026-10-19", "2026-10-21") // 3103.00
	book("deluxe", "2026-10-19", "2026-10-20")               // 2299.70
	cancel(book("deluxe", "2026-10-15", "2026-10-18"))       // 7755.10, refunded
	// Day 2: refunds exceed charges.
	s.clock.t = time.Date(2026, 10, 15, 11, 0, 0, 0, model.HotelTime)
	book("standard", "2026-10-23", "2026-10-24") // 1979.50
	cancel(standard)                             // 3103.00 refunded

	for _, day := range []struct {
		date                       model.Date
		charges, refunds, net, vat string
		payments                   string
	}{
		{"2026-10-14", "13157.80", "7755.10", "5402.70", "353.45", "4"},
		{"2026-10-15", "1979.50", "3103.00", "-1123.50", "-73.50", "2"},
	} {
		if _, err := job.Run(context.Background(), day.date); err != nil {
			t.Fatal(err)
		}
		got := getReport(t, s, day.date)
		for field, want := range map[string]string{
			"charges": day.charges, "refunds": day.refunds, "net": day.net, "vat": day.vat, "payments": day.payments,
		} {
			if !sameNumber(got[field], want) {
				t.Errorf("%s %s: got %s, want %s", day.date, field, got[field], want)
			}
		}
	}
}

// getReport returns the report's fields as the raw JSON numbers sent on the wire.
func getReport(t *testing.T, s *testServer, date model.Date) map[string]json.Number {
	t.Helper()
	resp, err := http.Get(s.URL + "/reports/revenue?date=" + string(date))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("report %s: status %d", date, resp.StatusCode)
	}
	var raw map[string]any
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(&raw); err != nil {
		t.Fatal(err)
	}
	out := map[string]json.Number{}
	for k, v := range raw {
		if n, ok := v.(json.Number); ok {
			out[k] = n
		}
	}
	return out
}

// sameNumber compares decimals exactly, so 3103 and 3103.00 are equal.
func sameNumber(got json.Number, want string) bool {
	g, ok1 := new(big.Rat).SetString(string(got))
	w, ok2 := new(big.Rat).SetString(want)
	return ok1 && ok2 && g.Cmp(w) == 0
}
