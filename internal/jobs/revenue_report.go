// Package jobs holds background work that runs outside an HTTP request.
package jobs

import (
	"context"

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
	var charges, refunds model.Money
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
		Charges:  charges,
		Refunds:  refunds,
		Net:      net,
		VAT:      includedVAT(net),
		Payments: len(payments),
	}
	return r, j.Reports.Save(ctx, r)
}

// includedVAT is the 7% VAT contained in a VAT-inclusive amount, rounded
// half away from zero to the satang.
func includedVAT(m model.Money) model.Money {
	if m < 0 {
		return -includedVAT(-m)
	}
	return (m*7*2 + 107) / (107 * 2)
}
