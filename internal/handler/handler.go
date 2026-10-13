// Package handler is the HTTP layer: it decodes JSON, calls a service and
// encodes the result.
package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/pangaunn/harness-engineering-workshop/internal/service/booking"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/payment"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/report"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/room"
)

type Services struct {
	Rooms    *room.Service
	Bookings *booking.Service
	Reports  *report.Service
}

type handler struct{ Services }

func New(s Services) http.Handler {
	h := &handler{s}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rooms", h.listRooms)
	mux.HandleFunc("POST /quotes", h.createQuote)
	mux.HandleFunc("POST /bookings", h.createBooking)
	mux.HandleFunc("POST /bookings/{id}/cancel", h.cancelBooking)
	mux.HandleFunc("GET /reports/revenue", h.revenueReport)
	return mux
}

func decode(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, booking.ErrInvalidRequest), errors.Is(err, errBadRequest):
		status = http.StatusBadRequest
	case errors.Is(err, booking.ErrNotFound), errors.Is(err, room.ErrUnknownType), errors.Is(err, report.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, booking.ErrNoAvailability), errors.Is(err, booking.ErrAlreadyCancelled), errors.Is(err, payment.ErrNoCharge):
		status = http.StatusConflict
	}
	msg := err.Error()
	if status == http.StatusInternalServerError {
		log.Printf("internal error: %v", err)
		msg = "internal error"
	}
	writeJSON(w, status, map[string]string{"error": msg})
}

var errBadRequest = errors.New("bad request")
