package domain

import "time"

type Hotel struct {
	ID          int64
	Name        string
	Location    string
	Description string
	CreatedAt   time.Time
}
