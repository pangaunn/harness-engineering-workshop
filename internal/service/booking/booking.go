// Package booking is the aggregate for quoting, reserving and cancelling a
// stay. It is the only place that combines the room, availability, pricing
// and payment domains.
package booking

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/availability"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/payment"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/pricing"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/room"
)

const maxNights = 30

var (
	ErrNotFound         = errors.New("booking not found")
	ErrAlreadyCancelled = errors.New("booking is already cancelled")
	ErrNoAvailability   = errors.New("no rooms available for those dates")
	ErrInvalidRequest   = errors.New("invalid booking request")
)

type Service struct {
	bookings     *repository.Bookings
	rooms        *room.Service
	availability *availability.Service
	payments     *payment.Service
	now          func() time.Time
}

func New(bookings *repository.Bookings, rooms *room.Service, avail *availability.Service, payments *payment.Service, now func() time.Time) *Service {
	return &Service{bookings: bookings, rooms: rooms, availability: avail, payments: payments, now: now}
}

type Request struct {
	RoomTypeID string
	GuestName  string
	CheckIn    model.Date
	CheckOut   model.Date
}

func (s *Service) Quote(ctx context.Context, roomTypeID string, checkIn, checkOut model.Date) (model.Quote, error) {
	if err := s.validateStay(checkIn, checkOut); err != nil {
		return model.Quote{}, err
	}
	rt, err := s.rooms.Type(ctx, roomTypeID)
	if err != nil {
		return model.Quote{}, err
	}
	q := pricing.Quote(rt, checkIn, checkOut)
	q.RoomsLeft, err = s.availability.CountFree(ctx, rt.ID, checkIn, checkOut)
	return q, err
}

func (s *Service) Reserve(ctx context.Context, req Request) (model.Booking, error) {
	if strings.TrimSpace(req.GuestName) == "" {
		return model.Booking{}, fmt.Errorf("%w: guest name is required", ErrInvalidRequest)
	}
	if err := s.validateStay(req.CheckIn, req.CheckOut); err != nil {
		return model.Booking{}, err
	}
	rt, err := s.rooms.Type(ctx, req.RoomTypeID)
	if err != nil {
		return model.Booking{}, err
	}
	free, err := s.availability.FreeRooms(ctx, rt.ID, req.CheckIn, req.CheckOut)
	if err != nil {
		return model.Booking{}, err
	}
	if len(free) == 0 {
		return model.Booking{}, ErrNoAvailability
	}
	q := pricing.Quote(rt, req.CheckIn, req.CheckOut)
	b, err := s.bookings.Create(ctx, model.Booking{
		RoomID:    free[0].ID,
		GuestName: req.GuestName,
		CheckIn:   req.CheckIn,
		CheckOut:  req.CheckOut,
		Total:     q.Total,
		Status:    model.StatusConfirmed,
		CreatedAt: s.now(),
	})
	if err != nil {
		return model.Booking{}, err
	}
	if _, err := s.payments.Charge(ctx, b.ID, b.Total); err != nil {
		return model.Booking{}, err
	}
	return b, nil
}

// Cancel cancels a booking and refunds the guest.
func (s *Service) Cancel(ctx context.Context, id int64) (model.Booking, model.Payment, error) {
	b, err := s.bookings.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return b, model.Payment{}, ErrNotFound
	}
	if err != nil {
		return b, model.Payment{}, err
	}
	if b.Status == model.StatusCancelled {
		return b, model.Payment{}, ErrAlreadyCancelled
	}
	if err := s.bookings.SetStatus(ctx, b.ID, model.StatusCancelled); err != nil {
		return b, model.Payment{}, err
	}
	b.Status = model.StatusCancelled
	refund, err := s.payments.Refund(ctx, b.ID)
	return b, refund, err
}

func (s *Service) validateStay(checkIn, checkOut model.Date) error {
	if checkIn == "" || checkOut == "" {
		return fmt.Errorf("%w: check_in and check_out are required", ErrInvalidRequest)
	}
	if checkOut <= checkIn {
		return fmt.Errorf("%w: check_out must be after check_in", ErrInvalidRequest)
	}
	if checkIn < model.DateOf(s.now()) {
		return fmt.Errorf("%w: check_in is in the past", ErrInvalidRequest)
	}
	if model.Nights(checkIn, checkOut) > maxNights {
		return fmt.Errorf("%w: stays are limited to %d nights", ErrInvalidRequest, maxNights)
	}
	return nil
}
