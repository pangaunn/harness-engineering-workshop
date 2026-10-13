package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

type Payments struct{ db *sql.DB }

func NewPayments(db *sql.DB) *Payments { return &Payments{db: db} }

const paymentColumns = `id, booking_id, kind, amount, created_at`

func (r *Payments) Create(ctx context.Context, p model.Payment) (model.Payment, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO payments (booking_id, kind, amount, created_at) VALUES (?, ?, ?, ?)`,
		p.BookingID, p.Kind, p.Amount, formatTime(p.CreatedAt))
	if err != nil {
		return p, err
	}
	p.ID, err = res.LastInsertId()
	return p, err
}

// ChargeFor returns the charge recorded for a booking.
func (r *Payments) ChargeFor(ctx context.Context, bookingID int64) (model.Payment, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+paymentColumns+` FROM payments WHERE booking_id = ? AND kind = ? ORDER BY id LIMIT 1`,
		bookingID, model.Charge)
	p, err := scanPayment(row)
	if errors.Is(err, sql.ErrNoRows) {
		return p, ErrNotFound
	}
	return p, err
}

// ListOn returns the payments created on the given hotel-local date.
func (r *Payments) ListOn(ctx context.Context, date model.Date) ([]model.Payment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+paymentColumns+` FROM payments WHERE substr(created_at, 1, 10) = ? ORDER BY id`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Payment
	for rows.Next() {
		p, err := scanPayment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func scanPayment(s scanner) (model.Payment, error) {
	var p model.Payment
	var created string
	err := s.Scan(&p.ID, &p.BookingID, &p.Kind, &p.Amount, &created)
	p.CreatedAt = parseTime(created)
	return p, err
}
