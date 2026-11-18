package domain

import "time"

type RoomStatus string

const (
	RoomStatusAvailable   RoomStatus = "available"
	RoomStatusBooked      RoomStatus = "booked"
	RoomStatusMaintenance RoomStatus = "maintenance"
)

type Room struct {
	ID        int64
	HotelID   int64
	Number    string
	Status    RoomStatus
	CreatedAt time.Time
}
