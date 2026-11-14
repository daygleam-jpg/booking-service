package room

import (
	"applicationDesignTest/internal/domain"
	"applicationDesignTest/internal/repository"
	"context"
	"fmt"
)

type RoomUseCase struct {
	repo repository.RoomRepository
}

func NewRoomUseCase(repo repository.RoomRepository) *RoomUseCase {
	return &RoomUseCase{repo: repo}
}

func (uc *RoomUseCase) GetByHotelID(ctx context.Context, hotelID int64) ([]domain.Room, error) {
	return uc.repo.GetByHotelID(ctx, hotelID)
}

func (uc *RoomUseCase) Create(ctx context.Context, room *domain.Room) error {
	if room.Number == "" {
		return fmt.Errorf("room number is required")
	}

	if room.Status == "" {
		room.Status = "available"
	}

	return uc.repo.Create(ctx, room)
}
