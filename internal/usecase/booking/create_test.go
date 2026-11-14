package booking

import (
	"applicationDesignTest/internal/domain"
	"applicationDesignTest/internal/repository"
	"context"
	"errors"
	"testing"
	"time"
)

type FakeRoomRepository struct {
	Room         domain.Room
	Available    bool
	GetRoomError error
}

func (r *FakeRoomRepository) GetRoomByID(roomID int) (domain.Room, error) {
	return r.Room, r.GetRoomError
}

func (r *FakeRoomRepository) IsRoomAvailable(roomID int, checkIn, checkOut time.Time) (bool, error) {
	return r.Available, nil
}

func (r *FakeRoomRepository) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	return &r.Room, nil
}

func (r *FakeRoomRepository) GetByHotelID(ctx context.Context, hotelID int64) ([]domain.Room, error) {
	return []domain.Room{}, nil
}

func (r *FakeRoomRepository) Create(ctx context.Context, room *domain.Room) error {
	return nil
}

type FakeBookingRepository struct {
	Called      bool
	ReturnError error
}

func (r *FakeBookingRepository) CreateBooking(booking domain.Booking) error {
	r.Called = true
	return r.ReturnError
}

func (r *FakeBookingRepository) CancelBooking(bookingID int) error {
	return nil
}

var _ repository.RoomRepository = (*FakeRoomRepository)(nil)
var _ repository.BookingRepository = (*FakeBookingRepository)(nil)

func TestBookingUseCase_CreateBooking_RoomNoFound(t *testing.T) {
	roomRepo := &FakeRoomRepository{
		GetRoomError: errors.New("room not found"),
	}
	bookingRepo := &FakeBookingRepository{}

	uc := NewBookingUseCase(roomRepo, bookingRepo)
	input := CreateBookingInput{}

	err := uc.CreateBooking(input)
	if err == nil {
		t.Fatal("expected error")
	}

	if bookingRepo.Called {
		t.Fatal("CreateBooking should not be called")
	}
}
