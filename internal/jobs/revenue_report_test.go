package jobs

import (
	"context"
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

	for i, amount := range []model.Money{310300, 229970, 775510} {
		if _, err := payments.Charge(ctx, int64(i+1), amount); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := payments.Refund(ctx, 3); err != nil {
		t.Fatal(err)
	}
	now = now.AddDate(0, 0, 1)
	if _, err := payments.Charge(ctx, 4, 197950); err != nil {
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
		got, want model.Money
	}{
		{"charges", got.Charges, 1315780},
		{"refunds", got.Refunds, 775510},
		{"net", got.Net, 540270},
		{"vat", got.VAT, 35345},
	} {
		if c.got != c.want {
			t.Errorf("%s: got %s, want %s", c.field, c.got, c.want)
		}
	}
}
