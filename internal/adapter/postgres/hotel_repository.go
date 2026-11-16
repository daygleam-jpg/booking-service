package postgres

import (
	"context"
	"database/sql"
	"errors"

	"applicationDesignTest/internal/domain"
	"applicationDesignTest/internal/repository"
)

type hotelRepository struct {
	db *sql.DB
}

func NewHotelRepository(db *sql.DB) repository.HotelRepository {
	return &hotelRepository{db: db}
}

func (r *hotelRepository) GetByID(ctx context.Context, id int64) (*domain.Hotel, error) {
	var hotel domain.Hotel

	query := `
    SELECT id, name, location, description, created_at
    FROM hotels
    WHERE id = $1
`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&hotel.ID,
		&hotel.Name,
		&hotel.Location,
		&hotel.Description,
		&hotel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("hotel not found")
		}
		return nil, err
	}

	return &hotel, nil
}

func (r *hotelRepository) GetAll(ctx context.Context) ([]domain.Hotel, error) {
	query := `
		SELECT id, name, location, description, created_at
		FROM hotels
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []domain.Hotel

	for rows.Next() {
		var hotel domain.Hotel

		err := rows.Scan(
			&hotel.ID,
			&hotel.Name,
			&hotel.Location,
			&hotel.Description,
			&hotel.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		hotels = append(hotels, hotel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (r *hotelRepository) Create(ctx context.Context, hotel *domain.Hotel) error {
	query := `
        INSERT INTO hotels (name, location, description)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	err := r.db.QueryRowContext(
		ctx,
		query,
		hotel.Name,
		hotel.Location,
		hotel.Description,
	).Scan(&hotel.ID, &hotel.CreatedAt)
	return err
}

func (r *hotelRepository) Update(ctx context.Context, hotel *domain.Hotel) error {
	query := `
        UPDATE hotels
        SET name = $1,
            location = $2,
            description = $3
        WHERE id = $4
    `

	result, err := r.db.ExecContext(
		ctx,
		query,
		hotel.Name,
		hotel.Location,
		hotel.Description,
		hotel.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("hotel not found")
	}

	return nil
}

func (r *hotelRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM hotels WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("hotel not found")
	}

	return nil
}
