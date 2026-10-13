// Package report stores and serves the daily revenue reports.
package report

import (
	"context"
	"errors"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
)

var ErrNotFound = errors.New("no report for that date")

type Service struct{ reports *repository.RevenueReports }

func New(reports *repository.RevenueReports) *Service { return &Service{reports: reports} }

func (s *Service) Save(ctx context.Context, r model.RevenueReport) error {
	return s.reports.Save(ctx, r)
}

func (s *Service) Get(ctx context.Context, date model.Date) (model.RevenueReport, error) {
	r, err := s.reports.Get(ctx, date)
	if errors.Is(err, repository.ErrNotFound) {
		return r, ErrNotFound
	}
	return r, err
}
