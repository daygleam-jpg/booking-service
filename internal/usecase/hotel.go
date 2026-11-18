package usecase

import (
	"context"

	"applicationDesignTest/internal/domain"
)

type CreateHotelInput struct {
	Name        string
	Location    string
	Description string
}

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

func (uc *HotelUseCase) Create(ctx context.Context, input CreateHotelInput) (*domain.Hotel, error) {
	hotel := &domain.Hotel{
		Name:        input.Name,
		Location:    input.Location,
		Description: input.Description,
	}

	if err := uc.repo.Create(ctx, hotel); err != nil {
		return nil, err
	}

	return hotel, nil
}

type UpdateHotelInput struct {
	ID          int64
	Name        string
	Location    string
	Description string
}

func (uc *HotelUseCase) Update(ctx context.Context, input UpdateHotelInput) (*domain.Hotel, error) {
	hotel := &domain.Hotel{
		ID:          input.ID,
		Name:        input.Name,
		Location:    input.Location,
		Description: input.Description,
	}

	if err := uc.repo.Update(ctx, hotel); err != nil {
		return nil, err
	}

	return uc.repo.GetByID(ctx, hotel.ID)
}

func (uc *HotelUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
