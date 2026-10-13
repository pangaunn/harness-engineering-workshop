package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/booking"
)

type stayRequest struct {
	RoomTypeID string `json:"room_type_id"`
	GuestName  string `json:"guest_name"`
	CheckIn    string `json:"check_in"`
	CheckOut   string `json:"check_out"`
}

func (s stayRequest) dates() (model.Date, model.Date, error) {
	in, err := model.ParseDate(s.CheckIn)
	if err != nil {
		return "", "", fmt.Errorf("%w: check_in: %v", errBadRequest, err)
	}
	out, err := model.ParseDate(s.CheckOut)
	if err != nil {
		return "", "", fmt.Errorf("%w: check_out: %v", errBadRequest, err)
	}
	return in, out, nil
}

func (h *handler) createQuote(w http.ResponseWriter, r *http.Request) {
	var req stayRequest
	if err := decode(r, &req); err != nil {
		writeError(w, fmt.Errorf("%w: %v", errBadRequest, err))
		return
	}
	in, out, err := req.dates()
	if err != nil {
		writeError(w, err)
		return
	}
	q, err := h.Bookings.Quote(r.Context(), req.RoomTypeID, in, out)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, q)
}

func (h *handler) createBooking(w http.ResponseWriter, r *http.Request) {
	var req stayRequest
	if err := decode(r, &req); err != nil {
		writeError(w, fmt.Errorf("%w: %v", errBadRequest, err))
		return
	}
	in, out, err := req.dates()
	if err != nil {
		writeError(w, err)
		return
	}
	b, err := h.Bookings.Reserve(r.Context(), booking.Request{
		RoomTypeID: req.RoomTypeID,
		GuestName:  req.GuestName,
		CheckIn:    in,
		CheckOut:   out,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

type cancelResponse struct {
	model.Booking
	Refunded float64 `json:"refunded"`
}

func (h *handler) cancelBooking(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, fmt.Errorf("%w: invalid booking id", errBadRequest))
		return
	}
	b, refund, err := h.Bookings.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cancelResponse{Booking: b, Refunded: refund.Amount})
}
