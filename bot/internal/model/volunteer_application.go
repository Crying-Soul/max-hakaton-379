package model

import "time"

// VolunteerApplication represents a volunteer's application to an event.
type VolunteerApplication struct {
	ID              int32
	EventID         *int32
	VolunteerID     *int64
	Status          *string
	RejectionReason *string
	ReviewedBy      *int64
	AppliedAt       time.Time
	ReviewedAt      *time.Time
}
