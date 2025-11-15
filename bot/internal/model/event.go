package model

import "time"

// Event describes a volunteer event.
// Соответствует таблице events.
type Event struct {
	ID                int32
	Title             string
	Description       *string
	Chat              *int64
	Date              time.Time
	DurationHours     *int32
	Location          string
	LocationLat       float64
	LocationLon       float64
	CategoryID        *int32
	OrganizerID       *int64
	Contacts          *string
	MaxVolunteers     int32
	CurrentVolunteers *int32
	Status            *string
	CancelledReason   *string
	CompletedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
