package model

// Volunteer stores volunteer-specific profile attributes.
// Соответствует таблице volunteers.
type Volunteer struct {
	ID           int64
	CV           *string
	SearchRadius *int32
	CategoryIDs  []int32
}
