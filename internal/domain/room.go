package domain

import "time"

type Room struct {
	ID        int64     `json:"id"`
	HotelID   int64     `json:"hotel_id"`
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
