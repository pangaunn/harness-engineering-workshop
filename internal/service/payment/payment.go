// Package payment is the charges and refunds domain. The gateway is fake:
// every payment succeeds and is recorded in the database.
package payment

import (
	"context"
	"errors"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
)

var ErrNoCharge = errors.New("booking has no charge to refund")

type Service struct {
	payments *repository.Payments
	now      func() time.Time
}

func New(payments *repository.Payments, now func() time.Time) *Service {
	return &Service{payments: payments, now: now}
}

func (s *Service) Charge(ctx context.Context, bookingID int64, amount float64) (model.Payment, error) {
	return s.record(ctx, bookingID, model.Charge, amount)
}

func (s *Service) Refund(ctx context.Context, bookingID int64) (model.Payment, error) {
	charge, err := s.payments.ChargeFor(ctx, bookingID)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Payment{}, ErrNoCharge
	}
	if err != nil {
		return model.Payment{}, err
	}
	return s.record(ctx, bookingID, model.Refund, charge.Amount)
}

// On returns every payment made on a hotel-local date.
func (s *Service) On(ctx context.Context, date model.Date) ([]model.Payment, error) {
	return s.payments.ListOn(ctx, date)
}

func (s *Service) record(ctx context.Context, bookingID int64, kind model.PaymentKind, amount float64) (model.Payment, error) {
	return s.payments.Create(ctx, model.Payment{
		BookingID: bookingID,
		Kind:      kind,
		Amount:    amount,
		CreatedAt: s.now().In(model.HotelTime),
	})
}
