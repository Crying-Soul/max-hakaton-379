package model

import "time"

type Role int

// User represents a bot user with optional geo metadata.
// Соответствует таблице users.
type User struct {
	ID          int64
	Username    *string
	Name        string
	Role        string
	State       string
	IsBlocked   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LocationLat *float64
	LocationLon *float64
}
