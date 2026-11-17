package usecase

import (
	"context"

	"applicationDesignTest/internal/domain"
)

type hotelRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Hotel, error)
	GetAll(ctx context.Context) ([]domain.Hotel, error)
	Create(ctx context.Context, hotel *domain.Hotel) error
	Update(ctx context.Context, hotel *domain.Hotel) error
	Delete(ctx context.Context, id int64) error
}

type HotelUseCase struct {
	repo hotelRepository
}

func NewHotelUseCase(repo hotelRepository) *HotelUseCase {
	return &HotelUseCase{
		repo: repo,
	}
}

func (uc *HotelUseCase) GetByID(ctx context.Context, id int64) (*domain.Hotel, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *HotelUseCase) GetAll(ctx context.Context) ([]domain.Hotel, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *HotelUseCase) Create(ctx context.Context, hotel *domain.Hotel) error {
	return uc.repo.Create(ctx, hotel)
}

func (uc *HotelUseCase) Update(ctx context.Context, hotel *domain.Hotel) (*domain.Hotel, error) {
	if err := uc.repo.Update(ctx, hotel); err != nil {
		return nil, err
	}

	return uc.repo.GetByID(ctx, hotel.ID)
}

func (uc *HotelUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
