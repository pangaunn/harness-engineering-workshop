package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

type RevenueReports struct{ db *sql.DB }

func NewRevenueReports(db *sql.DB) *RevenueReports { return &RevenueReports{db: db} }

func (r *RevenueReports) Save(ctx context.Context, rep model.RevenueReport) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO revenue_reports (date, charges, refunds, net, vat, payments) VALUES (?, ?, ?, ?, ?, ?)`,
		rep.Date, rep.Charges, rep.Refunds, rep.Net, rep.VAT, rep.Payments)
	return err
}

func (r *RevenueReports) Get(ctx context.Context, date model.Date) (model.RevenueReport, error) {
	var rep model.RevenueReport
	err := r.db.QueryRowContext(ctx,
		`SELECT date, charges, refunds, net, vat, payments FROM revenue_reports WHERE date = ?`, date).
		Scan(&rep.Date, &rep.Charges, &rep.Refunds, &rep.Net, &rep.VAT, &rep.Payments)
	if errors.Is(err, sql.ErrNoRows) {
		return rep, ErrNotFound
	}
	return rep, err
}
