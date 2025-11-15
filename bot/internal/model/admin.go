package model

import "time"

// Admin represents elevated privileges bound to a user record.
// Соответствует таблице admins.
type Admin struct {
	ID        int64
	CreatedAt time.Time
}
