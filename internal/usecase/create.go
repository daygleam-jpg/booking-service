package usecase

import (
	"errors"
	"time"

	"applicationDesignTest/internal/domain"
)

type bookingRoomRepository interface {
	GetRoomByID(roomID int) (domain.Room, error)
	IsRoomAvailable(roomID int, checkIn, checkOut time.Time) (bool, error)
}

type bookingRepository interface {
	CreateBooking(booking domain.Booking) (int64, error)
	CancelBooking(bookingID int) error
}

type BookingUseCase struct {
	roomRepo    bookingRoomRepository
	bookingRepo bookingRepository
}

func NewBookingUseCase(roomRepo bookingRoomRepository, bookingRepo bookingRepository) *BookingUseCase {
	return &BookingUseCase{
		roomRepo:    roomRepo,
		bookingRepo: bookingRepo,
	}
}

type CreateBookingInput struct {
	UserID   int
	HotelID  int
	RoomID   int
	CheckIn  time.Time
	CheckOut time.Time
}

func (uc *BookingUseCase) CreateBooking(input CreateBookingInput) (int64, error) {
	room, err := uc.roomRepo.GetRoomByID(input.RoomID)
	if err != nil {
		return 0, err
	}

	if int64(input.HotelID) != room.HotelID {
		return 0, errors.New("room does not belong to the specified hotel")
	}

	if room.Status == domain.RoomStatusMaintenance {
		return 0, errors.New("room is under maintenance")
	}

	available, err := uc.roomRepo.IsRoomAvailable(input.RoomID, input.CheckIn, input.CheckOut)
	if err != nil {
		return 0, err
	}

	if !available {
		return 0, errors.New("room is not available")
	}

	booking := domain.Booking{
		UserID:    input.UserID,
		HotelID:   input.HotelID,
		RoomID:    input.RoomID,
		CheckIn:   input.CheckIn,
		CheckOut:  input.CheckOut,
		CreatedAt: time.Now(),
		Status:    domain.StatusConfirmed,
	}

	if err := booking.Validate(); err != nil {
		return 0, err
	}

	bookingID, err := uc.bookingRepo.CreateBooking(booking)
	if err != nil {
		return 0, err
	}

	return bookingID, nil
}

func (uc *BookingUseCase) CancelBooking(bookingID int) error {
	return uc.bookingRepo.CancelBooking(bookingID)
}
