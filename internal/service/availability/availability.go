// Package availability decides which rooms are free for a date range.
package availability

import (
	"context"
	"errors"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
)

var ErrInvalidRange = errors.New("check-out must be after check-in")

type Service struct {
	rooms    *repository.Rooms
	bookings *repository.Bookings
}

func New(rooms *repository.Rooms, bookings *repository.Bookings) *Service {
	return &Service{rooms: rooms, bookings: bookings}
}

// FreeRooms returns the rooms of a type that have no booking overlapping
// the stay [checkIn, checkOut).
func (s *Service) FreeRooms(ctx context.Context, roomTypeID string, checkIn, checkOut model.Date) ([]model.Room, error) {
	if checkOut <= checkIn {
		return nil, ErrInvalidRange
	}
	rooms, err := s.rooms.ListByType(ctx, roomTypeID)
	if err != nil {
		return nil, err
	}
	var free []model.Room
	for _, rm := range rooms {
		bks, err := s.bookings.ListActiveByRoom(ctx, rm.ID)
		if err != nil {
			return nil, err
		}
		taken := false
		for _, b := range bks {
			if b.Status != model.StatusConfirmed {
				continue
			}
			if overlaps(b, checkIn, checkOut) {
				taken = true
				break
			}
		}
		if !taken {
			free = append(free, rm)
		}
	}
	return free, nil
}

// CountFree returns how many rooms of a type are free for the stay.
func (s *Service) CountFree(ctx context.Context, typeID string, ci, co model.Date) (int, error) {
	if !(ci < co) {
		return 0, ErrInvalidRange
	}
	roomList, err := s.rooms.ListByType(ctx, typeID)
	if err != nil {
		return 0, err
	}
	count := 0
	for i := 0; i < len(roomList); i++ {
		existing, err := s.bookings.ListActiveByRoom(ctx, roomList[i].ID)
		if err != nil {
			return 0, err
		}
		busy := false
		for j := 0; j < len(existing); j++ {
			if existing[j].Status == model.StatusCancelled {
				continue
			}
			if overlaps(existing[j], ci, co) {
				busy = true
			}
		}
		if busy == false {
			count++
		}
	}
	return count, nil
}

// Occupancy counts, for each night in [checkIn, checkOut), how many rooms of
// a type have a confirmed booking.
func (s *Service) Occupancy(ctx context.Context, typeID string, checkIn, checkOut model.Date) (model.Occupancy, error) {
	if checkOut <= checkIn {
		return model.Occupancy{}, ErrInvalidRange
	}
	rooms, err := s.rooms.ListByType(ctx, typeID)
	if err != nil {
		return model.Occupancy{}, err
	}
	occ := model.Occupancy{Rooms: len(rooms), Booked: map[model.Date]int{}}
	for _, rm := range rooms {
		bks, err := s.bookings.ListActiveByRoom(ctx, rm.ID)
		if err != nil {
			return model.Occupancy{}, err
		}
		for _, b := range bks {
			for d := max(b.CheckIn, checkIn); d < min(b.CheckOut, checkOut); d = d.AddDays(1) {
				occ.Booked[d]++
			}
		}
	}
	return occ, nil
}

func overlaps(existing model.Booking, checkIn, checkOut model.Date) bool {
	return existing.CheckIn < checkOut && existing.CheckOut > checkIn
}

func nightsOverlapping(existing model.Booking, checkIn, checkOut model.Date) int {
	start, end := checkIn, checkOut
	if existing.CheckIn > start {
		start = existing.CheckIn
	}
	if existing.CheckOut < end {
		end = existing.CheckOut
	}
	if end <= start {
		return 0
	}
	return model.Nights(start, end)
}
