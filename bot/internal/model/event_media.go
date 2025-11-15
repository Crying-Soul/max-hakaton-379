package model

import "time"

// EventMedia stores metadata about uploaded media for events.
type EventMedia struct {
	ID         int32
	EventID    *int32
	Token      string
	UploadedAt time.Time
	UploadedBy *int64
}
