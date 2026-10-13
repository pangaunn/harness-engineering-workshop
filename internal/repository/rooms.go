package repository

import (
	"context"
	"database/sql"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

type Rooms struct{ db *sql.DB }

func NewRooms(db *sql.DB) *Rooms { return &Rooms{db: db} }

func (r *Rooms) ListByType(ctx context.Context, typeID string) ([]model.Room, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, number, type_id FROM rooms WHERE type_id = ? ORDER BY number`, typeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Room
	for rows.Next() {
		var rm model.Room
		if err := rows.Scan(&rm.ID, &rm.Number, &rm.TypeID); err != nil {
			return nil, err
		}
		out = append(out, rm)
	}
	return out, rows.Err()
}
