package model

import "time"

// Organizer describes organizer-specific details.
// Соответствует таблице organizers.
type Organizer struct {
	ID               int64
	OrganizationName string
	About            *string
	VerifiedAt       *time.Time
	VerifiedBy       *int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// OrganizerVerificationRequest хранит историю заявок на верификацию.
type OrganizerVerificationRequest struct {
	ID               int32
	OrganizerID      int64
	Status           string
	OrganizerComment *string
	AdminComment     *string
	ReviewedBy       *int64
	SubmittedAt      time.Time
	ReviewedAt       *time.Time
}
