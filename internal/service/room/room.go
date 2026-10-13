// Package room is the room types and room inventory domain.
package room

import (
	"context"
	"errors"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
)

// CheckInHour is when rooms are ready, in hotel local time.
const CheckInHour = 14

var ErrUnknownType = errors.New("unknown room type")

type Service struct {
	types *repository.RoomTypes
	rooms *repository.Rooms
}

func New(types *repository.RoomTypes, rooms *repository.Rooms) *Service {
	return &Service{types: types, rooms: rooms}
}

func (s *Service) Types(ctx context.Context) ([]model.RoomType, error) {
	return s.types.List(ctx)
}

func (s *Service) Type(ctx context.Context, id string) (model.RoomType, error) {
	t, err := s.types.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return t, ErrUnknownType
	}
	return t, err
}

func (s *Service) Rooms(ctx context.Context, typeID string) ([]model.Room, error) {
	return s.rooms.ListByType(ctx, typeID)
}
