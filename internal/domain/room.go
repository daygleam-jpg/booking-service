package domain

import "time"

type RoomStatus string

const (
	RoomStatusAvailable   RoomStatus = "available"
	RoomStatusBooked      RoomStatus = "booked"
	RoomStatusMaintenance RoomStatus = "maintenance"
)

type Room struct {
	ID        int64      `json:"id"`
	HotelID   int64      `json:"hotel_id"`
	Number    string     `json:"number"`
	Status    RoomStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}
