package domain

import (
	"errors"
	"time"
)

type Booking struct {
	UserID    int
	HotelID   int
	RoomID    int
	CheckIn   time.Time
	CheckOut  time.Time
	CreatedAt time.Time
	Status    BookingStatus
}

type BookingStatus string

const (
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

func (b Booking) Validate() error {
	if b.UserID <= 0 || b.HotelID <= 0 || b.RoomID <= 0 {
		return errors.New("invalid booking identifiers")
	}

	if !startOfDay(b.CheckOut).After(startOfDay(b.CheckIn)) {
		return errors.New("departure date must be after arrival date")
	}
	if startOfDay(b.CheckIn).Before(startOfDay(time.Now())) {
		return errors.New("arrival date cannot be in the past")
	}
	if b.Status != StatusConfirmed && b.Status != StatusCancelled {
		return errors.New("invalid booking status")
	}
	return nil
}

func startOfDay(t time.Time) time.Time {
	year, month, day := t.Date()

	return time.Date(
		year,
		month,
		day,
		0, 0, 0, 0,
		t.Location(),
	)
}

func (b *Booking) Cancel() error {
	if b.Status == StatusCancelled {
		return errors.New("booking already cancelled")
	}
	b.Status = StatusCancelled
	return nil
}
