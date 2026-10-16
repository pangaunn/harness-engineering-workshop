package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/availability"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/booking"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/payment"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/report"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/room"
)

type clock struct{ t time.Time }

func (c *clock) Now() time.Time { return c.t }

type testServer struct {
	*httptest.Server
	clock    *clock
	payments *payment.Service
	reports  *report.Service
}

// newTestServer wires the whole app against a fresh in-memory database.
func newTestServer(t *testing.T) *testServer {
	t.Helper()
	db, err := repository.Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	clk := &clock{t: time.Date(2026, 10, 1, 10, 0, 0, 0, model.HotelTime)}
	bookingRepo := repository.NewBookings(db)
	rooms := room.New(repository.NewRoomTypes(db), repository.NewRooms(db))
	payments := payment.New(repository.NewPayments(db), clk.Now)
	bookings := booking.New(bookingRepo, rooms, availability.New(repository.NewRooms(db), bookingRepo), payments, clk.Now)
	reports := report.New(repository.NewRevenueReports(db))

	srv := httptest.NewServer(New(Services{Rooms: rooms, Bookings: bookings, Reports: reports}))
	t.Cleanup(func() {
		srv.Close()
		db.Close()
	})
	return &testServer{Server: srv, clock: clk, payments: payments, reports: reports}
}

// do sends a request and decodes the JSON response into out (if non-nil).
func (s *testServer) do(t *testing.T, method, path string, body, out any) int {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	req, err := http.NewRequest(method, s.URL+path, &buf)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			t.Fatalf("%s %s: decode: %v", method, path, err)
		}
	}
	return resp.StatusCode
}
