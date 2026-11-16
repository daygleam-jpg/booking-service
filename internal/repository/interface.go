package repository

import (
	"applicationDesignTest/internal/domain"
	"context"
	"time"
)

type BookingRepository interface {
	CreateBooking(booking domain.Booking) (int64, error)
	CancelBooking(bookingID int) error
}
type RoomRepository interface {
	GetRoomByID(roomID int) (domain.Room, error)
	IsRoomAvailable(roomID int, checkIn, checkOut time.Time) (bool, error)

	GetByID(ctx context.Context, id int64) (*domain.Room, error)
	GetByHotelID(ctx context.Context, hotelID int64) ([]domain.Room, error)
	Create(ctx context.Context, room *domain.Room) error
	HasConfirmedBooking(ctx context.Context, roomID int64, date time.Time) (bool, error)
}
