package domain

import (
	"time"
)

type Hotel struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
