package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

type RoomTypes struct{ db *sql.DB }

func NewRoomTypes(db *sql.DB) *RoomTypes { return &RoomTypes{db: db} }

const roomTypeColumns = `id, name, capacity, weekday_rate, weekend_rate`

func (r *RoomTypes) List(ctx context.Context) ([]model.RoomType, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+roomTypeColumns+` FROM room_types ORDER BY weekday_rate`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.RoomType
	for rows.Next() {
		var t model.RoomType
		if err := rows.Scan(&t.ID, &t.Name, &t.Capacity, &t.WeekdayRate, &t.WeekendRate); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *RoomTypes) Get(ctx context.Context, id string) (model.RoomType, error) {
	var t model.RoomType
	err := r.db.QueryRowContext(ctx, `SELECT `+roomTypeColumns+` FROM room_types WHERE id = ?`, id).
		Scan(&t.ID, &t.Name, &t.Capacity, &t.WeekdayRate, &t.WeekendRate)
	if errors.Is(err, sql.ErrNoRows) {
		return t, ErrNotFound
	}
	return t, err
}
