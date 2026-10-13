package jobs

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/payment"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/report"
)

func TestRevenueReport_Run(t *testing.T) {
	ctx := context.Background()
	db, err := repository.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	now := time.Date(2026, 10, 14, 9, 30, 0, 0, model.HotelTime)
	payments := payment.New(repository.NewPayments(db), func() time.Time { return now })
	reports := report.New(repository.NewRevenueReports(db))

	for i, amount := range []float64{3103.00, 2299.70, 7755.10} {
		if _, err := payments.Charge(ctx, int64(i+1), amount); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := payments.Refund(ctx, 3); err != nil {
		t.Fatal(err)
	}
	now = now.AddDate(0, 0, 1)
	if _, err := payments.Charge(ctx, 4, 1979.50); err != nil {
		t.Fatal(err)
	}

	job := &RevenueReport{Payments: payments, Reports: reports}
	if _, err := job.Run(ctx, "2026-10-14"); err != nil {
		t.Fatal(err)
	}

	got, err := reports.Get(ctx, "2026-10-14")
	if err != nil {
		t.Fatal(err)
	}
	if got.Payments != 4 {
		t.Errorf("payments: got %d, want 4", got.Payments)
	}
	for _, c := range []struct {
		field     string
		got, want float64
	}{
		{"charges", got.Charges, 13157.80},
		{"refunds", got.Refunds, 7755.10},
		{"net", got.Net, 5402.70},
		{"vat", got.VAT, 353.45},
	} {
		if math.Abs(c.got-c.want) > 0.001 {
			t.Errorf("%s: got %.2f, want %.2f", c.field, c.got, c.want)
		}
	}
}
