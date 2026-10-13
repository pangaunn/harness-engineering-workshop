package handler

import (
	"fmt"
	"net/http"
	"testing"
)

type bookingJSON struct {
	ID       int64   `json:"id"`
	RoomID   int64   `json:"room_id"`
	Status   string  `json:"status"`
	Total    float64 `json:"total"`
	Refunded float64 `json:"refunded"`
}

func TestRooms(t *testing.T) {
	s := newTestServer(t)
	var out struct {
		CheckInFrom string `json:"check_in_from"`
		RoomTypes   []struct {
			ID string `json:"id"`
		} `json:"room_types"`
	}
	if code := s.do(t, "GET", "/rooms", nil, &out); code != http.StatusOK {
		t.Fatalf("status %d", code)
	}
	if len(out.RoomTypes) != 3 || out.CheckInFrom != "14:00" {
		t.Errorf("unexpected response: %+v", out)
	}
}

func TestQuote(t *testing.T) {
	s := newTestServer(t)
	var q struct {
		Total     float64 `json:"total"`
		RoomsLeft int     `json:"rooms_left"`
		Nights    []any   `json:"nights"`
	}
	body := map[string]string{"room_type_id": "deluxe", "check_in": "2026-10-15", "check_out": "2026-10-18"}
	if code := s.do(t, "POST", "/quotes", body, &q); code != http.StatusOK {
		t.Fatalf("status %d", code)
	}
	if q.Total != 7755.10 || q.RoomsLeft != 2 || len(q.Nights) != 3 {
		t.Errorf("unexpected quote: %+v", q)
	}
}

func TestBookAndCancelEarly(t *testing.T) {
	s := newTestServer(t)
	var b bookingJSON
	body := map[string]string{"room_type_id": "standard", "guest_name": "Somchai", "check_in": "2026-10-12", "check_out": "2026-10-14"}
	if code := s.do(t, "POST", "/bookings", body, &b); code != http.StatusCreated {
		t.Fatalf("book: status %d", code)
	}
	if b.Total != 3103.00 || b.Status != "confirmed" {
		t.Errorf("unexpected booking: %+v", b)
	}

	var c bookingJSON
	if code := s.do(t, "POST", fmt.Sprintf("/bookings/%d/cancel", b.ID), nil, &c); code != http.StatusOK {
		t.Fatalf("cancel: status %d", code)
	}
	if c.Status != "cancelled" || c.Refunded != 3103.00 {
		t.Errorf("unexpected cancel: %+v", c)
	}
	if code := s.do(t, "POST", fmt.Sprintf("/bookings/%d/cancel", b.ID), nil, nil); code != http.StatusConflict {
		t.Errorf("second cancel: status %d, want 409", code)
	}
}

func TestBookingErrors(t *testing.T) {
	s := newTestServer(t)
	cases := []struct {
		name string
		body map[string]string
		want int
	}{
		{"bad date", map[string]string{"room_type_id": "standard", "guest_name": "A", "check_in": "12/10/2026", "check_out": "2026-10-14"}, http.StatusBadRequest},
		{"unknown type", map[string]string{"room_type_id": "penthouse", "guest_name": "A", "check_in": "2026-10-12", "check_out": "2026-10-14"}, http.StatusNotFound},
		{"no guest", map[string]string{"room_type_id": "standard", "check_in": "2026-10-12", "check_out": "2026-10-14"}, http.StatusBadRequest},
	}
	for _, c := range cases {
		if code := s.do(t, "POST", "/bookings", c.body, nil); code != c.want {
			t.Errorf("%s: status %d, want %d", c.name, code, c.want)
		}
	}
	if code := s.do(t, "POST", "/bookings/999/cancel", nil, nil); code != http.StatusNotFound {
		t.Errorf("cancel unknown: status %d, want 404", code)
	}
}

func TestRevenueReport_NotFound(t *testing.T) {
	s := newTestServer(t)
	if code := s.do(t, "GET", "/reports/revenue?date=2026-10-01", nil, nil); code != http.StatusNotFound {
		t.Errorf("status %d, want 404", code)
	}
	if code := s.do(t, "GET", "/reports/revenue?date=yesterday", nil, nil); code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", code)
	}
}
