package http

import (
	"time"

	"applicationDesignTest/internal/domain"
)

type CreateRoomRequest struct {
	Number string            `json:"number"`
	Status domain.RoomStatus `json:"status"`
}

type RoomResponse struct {
	ID        int64             `json:"id"`
	HotelID   int64             `json:"hotel_id"`
	Number    string            `json:"number"`
	Status    domain.RoomStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}
