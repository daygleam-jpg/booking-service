package booking

import (
	"applicationDesignTest/internal/domain"
	"applicationDesignTest/internal/repository"
	"errors"
	"time"
)

type BookingUseCase struct {
	roomRepo    repository.RoomRepository
	bookingRepo repository.BookingRepository
}

func NewBookingUseCase(roomRepo repository.RoomRepository, bookingRepo repository.BookingRepository) *BookingUseCase {
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

func (uc *BookingUseCase) CreateBooking(input CreateBookingInput) error {
	room, err := uc.roomRepo.GetRoomByID(input.RoomID)
	if err != nil {
		return err
	}
	if int64(input.HotelID) != room.HotelID {
		return errors.New("room does not belong to the specified hotel")
	}

	available, err := uc.roomRepo.IsRoomAvailable(input.RoomID, input.CheckIn, input.CheckOut)
	if err != nil {
		return err
	}
	if !available {
		return errors.New("room is not available")
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

	err = booking.Validate()
	if err != nil {
		return err
	}

	err = uc.bookingRepo.CreateBooking(booking)
	if err != nil {
		return err
	}
	return nil
}

func (uc *BookingUseCase) CancelBooking(bookingID int) error {
	return uc.bookingRepo.CancelBooking(bookingID)
}
