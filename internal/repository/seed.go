package repository

import (
	"context"
	"database/sql"
)

func seed(ctx context.Context, db *sql.DB) error {
	stmts := []string{
		`INSERT INTO room_types VALUES
			('standard', 'Standard Double', 2, 1450.00, 1850.00),
			('deluxe',   'Deluxe River View', 2, 2149.25, 2549.25),
			('suite',    'Riverside Suite', 4, 4350.00, 5250.00)`,
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
