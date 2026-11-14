package usecase

import (
	"context"

	"applicationDesignTest/internal/domain"
	"applicationDesignTest/internal/repository"
)

type HotelUseCase struct {
	repo repository.HotelRepository
}

func NewHotelUseCase(repo repository.HotelRepository) *HotelUseCase {
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
func (uc *HotelUseCase) Update(ctx context.Context, hotel *domain.Hotel) error {
	return uc.repo.Update(ctx, hotel)
}
func (uc *HotelUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
