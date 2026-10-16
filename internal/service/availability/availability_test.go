package availability

import (
	"context"
	"testing"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
)

func setup(t *testing.T) (*Service, *repository.Bookings) {
	t.Helper()
	db, err := repository.Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	bookings := repository.NewBookings(db)
	return New(repository.NewRooms(db), bookings), bookings
}

func book(t *testing.T, repo *repository.Bookings, roomID int64, in, out model.Date, status model.BookingStatus) {
	t.Helper()
	_, err := repo.Create(context.Background(), model.Booking{
		RoomID: roomID, GuestName: "Existing Guest", CheckIn: in, CheckOut: out,
		Total: 1000, Status: status, CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
}

// The suite type has exactly one room (id 6).
const suiteRoom = 6

func TestFreeRooms_OverlappingStay(t *testing.T) {
	svc, repo := setup(t)
	book(t, repo, suiteRoom, "2026-11-10", "2026-11-14", model.StatusConfirmed)

	cases := []struct {
		name    string
		in, out model.Date
		free    bool
	}{
		{"inside", "2026-11-11", "2026-11-12", false},
		{"covers", "2026-11-09", "2026-11-15", false},
		{"starts during", "2026-11-13", "2026-11-16", false},
		{"ends during", "2026-11-08", "2026-11-11", false},
		{"before", "2026-11-05", "2026-11-08", true},
		{"after", "2026-11-16", "2026-11-18", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			free, err := svc.FreeRooms(context.Background(), "suite", c.in, c.out)
			if err != nil {
				t.Fatal(err)
			}
			if got := len(free) == 1; got != c.free {
				t.Errorf("free = %v, want %v", got, c.free)
			}
		})
	}
}

func TestFreeRooms_IgnoresCancelled(t *testing.T) {
	svc, repo := setup(t)
	book(t, repo, suiteRoom, "2026-11-10", "2026-11-14", model.StatusCancelled)

	free, err := svc.FreeRooms(context.Background(), "suite", "2026-11-11", "2026-11-12")
	if err != nil {
		t.Fatal(err)
	}
	if len(free) != 1 {
		t.Errorf("got %d free rooms, want 1", len(free))
	}
}

func TestCountFree(t *testing.T) {
	svc, repo := setup(t)
	book(t, repo, 1, "2026-11-10", "2026-11-12", model.StatusConfirmed)

	n, err := svc.CountFree(context.Background(), "standard", "2026-11-10", "2026-11-11")
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("got %d, want 2", n)
	}
}

func TestFreeRooms_InvalidRange(t *testing.T) {
	svc, _ := setup(t)
	if _, err := svc.FreeRooms(context.Background(), "suite", "2026-11-12", "2026-11-12"); err != ErrInvalidRange {
		t.Errorf("got %v, want ErrInvalidRange", err)
	}
}

func TestOccupancy_PerNight(t *testing.T) {
	svc, repo := setup(t)
	book(t, repo, 1, "2026-11-09", "2026-11-14", model.StatusConfirmed)
	book(t, repo, 2, "2026-11-10", "2026-11-14", model.StatusConfirmed)
	book(t, repo, 3, "2026-11-10", "2026-11-12", model.StatusCancelled)
	book(t, repo, 3, "2026-11-08", "2026-11-10", model.StatusConfirmed) // checks out on 11-10

	occ, err := svc.Occupancy(context.Background(), "standard", "2026-11-10", "2026-11-12")
	if err != nil {
		t.Fatal(err)
	}
	if occ.Rooms != 3 {
		t.Errorf("rooms: got %d, want 3", occ.Rooms)
	}
	for night, want := range map[model.Date]int{"2026-11-10": 2, "2026-11-11": 2} {
		if got := occ.Booked[night]; got != want {
			t.Errorf("%s: got %d booked, want %d", night, got, want)
		}
	}
	if len(occ.Booked) != 2 {
		t.Errorf("booked nights outside the stay: %v", occ.Booked)
	}
}
