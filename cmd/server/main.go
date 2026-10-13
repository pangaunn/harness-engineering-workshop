package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pangaunn/harness-engineering-workshop/internal/handler"
	"github.com/pangaunn/harness-engineering-workshop/internal/jobs"
	"github.com/pangaunn/harness-engineering-workshop/internal/model"
	"github.com/pangaunn/harness-engineering-workshop/internal/repository"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/availability"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/booking"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/payment"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/report"
	"github.com/pangaunn/harness-engineering-workshop/internal/service/room"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db, err := repository.Open(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bookingRepo := repository.NewBookings(db)
	rooms := room.New(repository.NewRoomTypes(db), repository.NewRooms(db))
	payments := payment.New(repository.NewPayments(db), time.Now)
	reports := report.New(repository.NewRevenueReports(db))
	avail := availability.New(repository.NewRooms(db), bookingRepo)
	bookings := booking.New(bookingRepo, rooms, avail, payments, time.Now)

	go runNightly(ctx, &jobs.RevenueReport{Payments: payments, Reports: reports})

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: handler.New(handler.Services{Rooms: rooms, Bookings: bookings, Reports: reports}),
	}
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()
	log.Printf("listening on %s", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// runNightly builds yesterday's revenue report shortly after midnight.
func runNightly(ctx context.Context, job *jobs.RevenueReport) {
	for {
		now := time.Now().In(model.HotelTime)
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 5, 0, 0, model.HotelTime)
		select {
		case <-ctx.Done():
			return
		case <-time.After(next.Sub(now)):
		}
		yesterday := model.DateOf(time.Now()).AddDays(-1)
		if _, err := job.Run(ctx, yesterday); err != nil {
			log.Printf("revenue report %s: %v", yesterday, err)
		}
	}
}
