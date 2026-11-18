package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"applicationDesignTest/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
)

type roomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *roomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *domain.Room) error {
	query := `
        INSERT INTO rooms (hotel_id, number, status)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		room.HotelID,
		room.Number,
		room.Status,
	).Scan(&room.ID, &room.CreatedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				return fmt.Errorf("room number already exists")
			case "23503":
				return fmt.Errorf("hotel not found")
			}
		}
		return err
	}

	return nil
}

func (r *roomRepository) GetByHotelID(ctx context.Context, hotelID int64) ([]domain.Room, error) {
	query := `
        SELECT id, hotel_id, number, status, created_at
        FROM rooms
        WHERE hotel_id = $1
        ORDER BY number
    `
	rows, err := r.db.QueryContext(ctx, query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room

	for rows.Next() {
		var room domain.Room

		err := rows.Scan(
			&room.ID,
			&room.HotelID,
			&room.Number,
			&room.Status,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	query := `
        SELECT id, hotel_id, number, status, created_at
        FROM rooms
        WHERE id = $1
    `

	var room domain.Room

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.HotelID,
		&room.Number,
		&room.Status,
		&room.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &room, nil
}
func (r *roomRepository) GetRoomByID(roomID int) (domain.Room, error) {
	var room domain.Room

	query := `
        SELECT id, hotel_id, number, status, created_at
        FROM rooms
        WHERE id = $1
    `

	err := r.db.QueryRow(query, roomID).Scan(
		&room.ID,
		&room.HotelID,
		&room.Number,
		&room.Status,
		&room.CreatedAt,
	)

	return room, err
}

func (r *roomRepository) IsRoomAvailable(roomID int, checkIn, checkOut time.Time) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM rooms
			WHERE id = $1
			  AND status = 'available'
		)
		AND NOT EXISTS (
			SELECT 1
			FROM bookings
			WHERE room_id = $1
			  AND status = 'confirmed'
			  AND check_in < $3
			  AND check_out > $2
		)
	`

	var available bool

	err := r.db.QueryRow(
		query,
		roomID,
		checkIn,
		checkOut,
	).Scan(&available)

	return available, err
}

func (r *roomRepository) HasConfirmedBooking(ctx context.Context, roomID int64, date time.Time) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM bookings
			WHERE room_id = $1
			  AND status = 'confirmed'
			  AND check_in <= $2
			  AND check_out > $2
		)
	`

	var booked bool

	err := r.db.QueryRowContext(ctx, query, roomID, date).Scan(&booked)

	return booked, err
}
