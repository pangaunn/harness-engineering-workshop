package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

type Bookings struct{ db *sql.DB }

func NewBookings(db *sql.DB) *Bookings { return &Bookings{db: db} }

const bookingColumns = `id, room_id, guest_name, check_in, check_out, total, status, created_at`

func (r *Bookings) Create(ctx context.Context, b model.Booking) (model.Booking, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO bookings (room_id, guest_name, check_in, check_out, total, status, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		b.RoomID, b.GuestName, b.CheckIn, b.CheckOut, b.Total, b.Status, formatTime(b.CreatedAt))
	if err != nil {
		return b, err
	}
	b.ID, err = res.LastInsertId()
	return b, err
}

func (r *Bookings) Get(ctx context.Context, id int64) (model.Booking, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+bookingColumns+` FROM bookings WHERE id = ?`, id)
	b, err := scanBooking(row)
	if errors.Is(err, sql.ErrNoRows) {
		return b, ErrNotFound
	}
	return b, err
}

// ListActiveByRoom returns the confirmed bookings for a room.
func (r *Bookings) ListActiveByRoom(ctx context.Context, roomID int64) ([]model.Booking, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+bookingColumns+` FROM bookings WHERE room_id = ? AND status = ? ORDER BY check_in`,
		roomID, model.StatusConfirmed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Booking
	for rows.Next() {
		b, err := scanBooking(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *Bookings) SetStatus(ctx context.Context, id int64, status model.BookingStatus) error {
	_, err := r.db.ExecContext(ctx, `UPDATE bookings SET status = ? WHERE id = ?`, status, id)
	return err
}

type scanner interface{ Scan(dest ...any) error }

func scanBooking(s scanner) (model.Booking, error) {
	var b model.Booking
	var created string
	err := s.Scan(&b.ID, &b.RoomID, &b.GuestName, &b.CheckIn, &b.CheckOut, &b.Total, &b.Status, &created)
	b.CreatedAt = parseTime(created)
	return b, err
}
