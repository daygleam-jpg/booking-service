package postgres

import (
	"database/sql"

	"applicationDesignTest/internal/domain"
	"applicationDesignTest/internal/repository"
)

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) repository.BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CreateBooking(booking domain.Booking) (int64, error) {
	query := `
		INSERT INTO bookings (
			user_id, hotel_id, room_id,
			check_in, check_out, status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var bookingID int64

	err := r.db.QueryRow(
		query,
		booking.UserID,
		booking.HotelID,
		booking.RoomID,
		booking.CheckIn,
		booking.CheckOut,
		booking.Status,
	).Scan(&bookingID)

	return bookingID, err
}

func (r *bookingRepository) CancelBooking(bookingID int) error {
	query := `
		UPDATE bookings
		SET status = 'cancelled'
		WHERE id = $1 AND status = 'confirmed'
	`

	result, err := r.db.Exec(query, bookingID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
