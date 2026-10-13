// Package jobs holds background work that runs outside an HTTP request.
package jobs

import (
	"context"
	"math"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/payment"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/report"
)

// RevenueReport builds the nightly revenue report for one day.
type RevenueReport struct {
	Payments *payment.Service
	Reports  *report.Service
}

func (j *RevenueReport) Run(ctx context.Context, date model.Date) (model.RevenueReport, error) {
	payments, err := j.Payments.On(ctx, date)
	if err != nil {
		return model.RevenueReport{}, err
	}
	var charges, refunds float64
	for _, p := range payments {
		switch p.Kind {
		case model.Charge:
			charges += p.Amount
		case model.Refund:
			refunds += p.Amount
		}
	}
	net := charges - refunds
	r := model.RevenueReport{
		Date:     date,
		Charges:  round2(charges),
		Refunds:  round2(refunds),
		Net:      round2(net),
		VAT:      round2(net * 7 / 107), // prices include 7% VAT
		Payments: len(payments),
	}
	return r, j.Reports.Save(ctx, r)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
