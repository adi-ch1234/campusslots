package domain

import (
	"time"
)

const (
	WorkshopStatusOpen   = "OPEN"
	WorkshopStatusClosed = "CLOSED"
)

// Workshop represents a scheduled college workshop and its live seat count.
type Workshop struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	Location       string    `json:"location"`
	StartsAt       time.Time `json:"starts_at"`
	Status         string    `json:"status"`
	Capacity       int       `json:"capacity"`
	AvailableSeats int       `json:"available_seats"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateWorkshopInput holds the fields required to create a new workshop.
type CreateWorkshopInput struct {
	Title    string
	Location string
	StartsAt time.Time
	Capacity int
	Status   string
}
