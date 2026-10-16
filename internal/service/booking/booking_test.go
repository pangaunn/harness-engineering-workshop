package booking

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/availability"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/payment"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/room"
)

var today = time.Date(2026, 10, 1, 10, 0, 0, 0, model.HotelTime)

func setup(t *testing.T) (*Service, *repository.Payments) {
	t.Helper()
	db, err := repository.Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	now := func() time.Time { return today }
	bookings := repository.NewBookings(db)
	payments := repository.NewPayments(db)
	svc := New(bookings,
		room.New(repository.NewRoomTypes(db), repository.NewRooms(db)),
		availability.New(repository.NewRooms(db), bookings),
		payment.New(payments, now),
		now)
	return svc, payments
}

func TestReserve_ChargesQuotedTotal(t *testing.T) {
	svc, payments := setup(t)
	ctx := context.Background()

	b, err := svc.Reserve(ctx, Request{RoomTypeID: "deluxe", GuestName: "Somchai", CheckIn: "2026-10-15", CheckOut: "2026-10-18"})
	if err != nil {
		t.Fatal(err)
	}
	if b.Total != 775510 {
		t.Errorf("total: got %s, want 7755.10", b.Total)
	}
	charge, err := payments.ChargeFor(ctx, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if charge.Amount != b.Total {
		t.Errorf("charged %s, want %s", charge.Amount, b.Total)
	}
}

func TestReserve_NoAvailability(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	req := Request{RoomTypeID: "suite", GuestName: "Somchai", CheckIn: "2026-10-20", CheckOut: "2026-10-22"}

	if _, err := svc.Reserve(ctx, req); err != nil {
		t.Fatal(err)
	}
	req.GuestName = "Malee"
	if _, err := svc.Reserve(ctx, req); !errors.Is(err, ErrNoAvailability) {
		t.Errorf("got %v, want ErrNoAvailability", err)
	}
}

func TestReserve_Validation(t *testing.T) {
	svc, _ := setup(t)
	cases := map[string]Request{
		"no guest":   {RoomTypeID: "standard", CheckIn: "2026-10-20", CheckOut: "2026-10-22"},
		"reversed":   {RoomTypeID: "standard", GuestName: "A", CheckIn: "2026-10-22", CheckOut: "2026-10-20"},
		"in past":    {RoomTypeID: "standard", GuestName: "A", CheckIn: "2026-09-20", CheckOut: "2026-09-22"},
		"too long":   {RoomTypeID: "standard", GuestName: "A", CheckIn: "2026-10-20", CheckOut: "2026-12-20"},
		"same dates": {RoomTypeID: "standard", GuestName: "A", CheckIn: "2026-10-20", CheckOut: "2026-10-20"},
	}
	for name, req := range cases {
		if _, err := svc.Reserve(context.Background(), req); !errors.Is(err, ErrInvalidRequest) {
			t.Errorf("%s: got %v, want ErrInvalidRequest", name, err)
		}
	}
}

func TestCancel_Twice(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	b, err := svc.Reserve(ctx, Request{RoomTypeID: "standard", GuestName: "Somchai", CheckIn: "2026-10-20", CheckOut: "2026-10-22"})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.Cancel(ctx, b.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.Cancel(ctx, b.ID); !errors.Is(err, ErrAlreadyCancelled) {
		t.Errorf("got %v, want ErrAlreadyCancelled", err)
	}
}

func TestReserve_SurchargeWhenTwoThirdsBooked(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	// Tuesday night, standard: 1450.00 + 101.50 VAT.
	req := Request{RoomTypeID: "standard", GuestName: "Somchai", CheckIn: "2026-10-20", CheckOut: "2026-10-21"}
	for range 2 {
		b, err := svc.Reserve(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if b.Total != 155150 {
			t.Errorf("before two thirds: total %s, want 1551.50", b.Total)
		}
	}
	// Two of three rooms are booked: 1740.00 + 121.80 VAT.
	q, err := svc.Quote(ctx, "standard", "2026-10-20", "2026-10-21")
	if err != nil {
		t.Fatal(err)
	}
	if q.Total != 186180 {
		t.Errorf("quote: total %s, want 1861.80", q.Total)
	}
	b, err := svc.Reserve(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if b.Total != 186180 {
		t.Errorf("reserve: total %s, want 1861.80", b.Total)
	}
}
