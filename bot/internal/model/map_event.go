package model

import "time"

// MapEvent представляет мероприятие для карты волонтёров с вычисленным расстоянием.
type MapEvent struct {
	ID                int32      `json:"id"`
	Title             string     `json:"title"`
	Description       *string    `json:"description,omitempty"`
	Chat              *int64     `json:"chat,omitempty"`
	Date              time.Time  `json:"date"`
	DurationHours     *int32     `json:"durationHours,omitempty"`
	Location          string     `json:"location"`
	LocationLat       float64    `json:"locationLat"`
	LocationLon       float64    `json:"locationLon"`
	CategoryID        *int32     `json:"categoryId,omitempty"`
	CategoryName      *string    `json:"categoryName,omitempty"`
	OrganizerID       *int64     `json:"organizerId,omitempty"`
	Contacts          *string    `json:"contacts,omitempty"`
	MaxVolunteers     int32      `json:"maxVolunteers"`
	CurrentVolunteers int32      `json:"currentVolunteers"`
	SlotsLeft         int32      `json:"slotsLeft"`
	DistanceKm        float64    `json:"distanceKm"`
	Status            *string    `json:"status,omitempty"`
	CancelledReason   *string    `json:"cancelledReason,omitempty"`
	CompletedAt       *time.Time `json:"completedAt,omitempty"`
	ApplicationStatus *string    `json:"applicationStatus,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
