package repository

import (
	"context"

	"applicationDesignTest/internal/domain"
)

type HotelRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Hotel, error)
	GetAll(ctx context.Context) ([]domain.Hotel, error)
	Create(ctx context.Context, hotel *domain.Hotel) error
	Update(ctx context.Context, hotel *domain.Hotel) error
	Delete(ctx context.Context, id int64) error
}
