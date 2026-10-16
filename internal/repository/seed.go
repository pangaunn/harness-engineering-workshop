package repository

import (
	"context"
	"database/sql"
)

func seed(ctx context.Context, db *sql.DB) error {
	stmts := []string{
		`INSERT INTO room_types VALUES
			('standard', 'Standard Double', 2, 145000, 185000),
			('deluxe',   'Deluxe River View', 2, 214925, 254925),
			('suite',    'Riverside Suite', 4, 435000, 525000)`,
		`INSERT INTO rooms (number, type_id) VALUES
			('101', 'standard'), ('102', 'standard'), ('103', 'standard'),
			('201', 'deluxe'), ('202', 'deluxe'),
			('301', 'suite')`,
	}
	for _, s := range stmts {
		if _, err := db.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
