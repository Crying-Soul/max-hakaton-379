package model

import "time"

// EventParticipant represents a volunteer assigned to an event.
// Соответствует таблице event_participants.
type EventParticipant struct {
	ID            int32
	EventID       *int32
	VolunteerID   *int64
	ApplicationID *int32
	JoinedChatAt  time.Time
}
