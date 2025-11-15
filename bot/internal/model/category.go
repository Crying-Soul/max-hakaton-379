package model

import "time"

// Category represents an event category descriptor.
type Category struct {
	ID          int32
	Name        string
	Description *string
	IsActive    *bool
	CreatedAt   time.Time
}
