// Package repository is the SQLite data access layer, one file per table.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("not found")

const schema = `
CREATE TABLE room_types (
	id           TEXT PRIMARY KEY,
	name         TEXT NOT NULL,
	capacity     INTEGER NOT NULL,
	weekday_rate INTEGER NOT NULL,
	weekend_rate INTEGER NOT NULL
);
CREATE TABLE rooms (
	id      INTEGER PRIMARY KEY,
	number  TEXT NOT NULL UNIQUE,
	type_id TEXT NOT NULL REFERENCES room_types(id)
);
CREATE TABLE bookings (
	id         INTEGER PRIMARY KEY,
	room_id    INTEGER NOT NULL REFERENCES rooms(id),
	guest_name TEXT NOT NULL,
	check_in   TEXT NOT NULL,
	check_out  TEXT NOT NULL,
	total      INTEGER NOT NULL,
	status     TEXT NOT NULL,
	created_at TEXT NOT NULL
);
CREATE TABLE payments (
	id         INTEGER PRIMARY KEY,
	booking_id INTEGER NOT NULL REFERENCES bookings(id),
	kind       TEXT NOT NULL,
	amount     INTEGER NOT NULL,
	created_at TEXT NOT NULL
);
CREATE TABLE revenue_reports (
	date     TEXT PRIMARY KEY,
	charges  INTEGER NOT NULL,
	refunds  INTEGER NOT NULL,
	net      INTEGER NOT NULL,
	vat      INTEGER NOT NULL,
	payments INTEGER NOT NULL
);
`

// Open creates an in-memory database with the schema and seed data loaded.
func Open(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	// Every connection to ":memory:" is a new database, so keep exactly one.
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(ctx, schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}
	if err := seed(ctx, db); err != nil {
		db.Close()
		return nil, fmt.Errorf("seed: %w", err)
	}
	return db, nil
}

func formatTime(t time.Time) string { return t.Format(time.RFC3339) }

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
