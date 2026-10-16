// Package model holds the plain data types shared by every layer.
package model

import "time"

// HotelTime is the hotel's local time zone.
var HotelTime = time.FixedZone("ICT", 7*60*60)

type RoomType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Capacity    int    `json:"capacity"`
	WeekdayRate Money  `json:"weekday_rate"`
	WeekendRate Money  `json:"weekend_rate"`
	// FromPrice is the lowest nightly price including VAT, for the rooms
	// page. Only RoomTypes.List fills it.
	FromPrice Money `json:"from_price"`
}

type Room struct {
	ID     int64  `json:"id"`
	Number string `json:"number"`
	TypeID string `json:"type_id"`
}

type BookingStatus string

const (
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID        int64         `json:"id"`
	RoomID    int64         `json:"room_id"`
	GuestName string        `json:"guest_name"`
	CheckIn   Date          `json:"check_in"`
	CheckOut  Date          `json:"check_out"`
	Total     Money         `json:"total"`
	Status    BookingStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

type PaymentKind string

const (
	Charge PaymentKind = "charge"
	Refund PaymentKind = "refund"
)

type Payment struct {
	ID        int64       `json:"id"`
	BookingID int64       `json:"booking_id"`
	Kind      PaymentKind `json:"kind"`
	Amount    Money       `json:"amount"`
	CreatedAt time.Time   `json:"created_at"`
}

type NightPrice struct {
	Date    Date  `json:"date"`
	Weekend bool  `json:"weekend"`
	Rate    Money `json:"rate"`
	VAT     Money `json:"vat"`
}

type Quote struct {
	RoomTypeID string       `json:"room_type_id"`
	CheckIn    Date         `json:"check_in"`
	CheckOut   Date         `json:"check_out"`
	Nights     []NightPrice `json:"nights"`
	Subtotal   Money        `json:"subtotal"`
	VAT        Money        `json:"vat"`
	Total      Money        `json:"total"`
	RoomsLeft  int          `json:"rooms_left"`
}

type RevenueReport struct {
	Date     Date  `json:"date"`
	Charges  Money `json:"charges"`
	Refunds  Money `json:"refunds"`
	Net      Money `json:"net"`
	VAT      Money `json:"vat"`
	Payments int   `json:"payments"`
}
