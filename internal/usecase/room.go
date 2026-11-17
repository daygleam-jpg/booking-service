package usecase

import (
	"context"
	"fmt"
	"time"

	"applicationDesignTest/internal/domain"
)

type roomRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Room, error)
	GetByHotelID(ctx context.Context, hotelID int64) ([]domain.Room, error)
	Create(ctx context.Context, room *domain.Room) error
	HasConfirmedBooking(ctx context.Context, roomID int64, date time.Time) (bool, error)
}

type RoomUseCase struct {
	repo roomRepository
}

func NewRoomUseCase(repo roomRepository) *RoomUseCase {
	return &RoomUseCase{repo: repo}
}

func (uc *RoomUseCase) GetByHotelID(ctx context.Context, hotelID int64) ([]domain.Room, error) {
	rooms, err := uc.repo.GetByHotelID(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	today := time.Now()

	for i := range rooms {
		if rooms[i].Status == domain.RoomStatusMaintenance {
			continue
		}

		booked, err := uc.repo.HasConfirmedBooking(ctx, rooms[i].ID, today)
		if err != nil {
			return nil, err
		}

		if booked {
			rooms[i].Status = domain.RoomStatusBooked
		} else {
			rooms[i].Status = domain.RoomStatusAvailable
		}
	}

	return rooms, nil
}

func (uc *RoomUseCase) Create(ctx context.Context, room *domain.Room) error {
	if room.Number == "" {
		return fmt.Errorf("room number is required")
	}

	if room.Status == "" {
		room.Status = domain.RoomStatusAvailable
	}

	if room.Status != domain.RoomStatusAvailable &&
		room.Status != domain.RoomStatusMaintenance {
		return fmt.Errorf("invalid room status")
	}

	return uc.repo.Create(ctx, room)
}
