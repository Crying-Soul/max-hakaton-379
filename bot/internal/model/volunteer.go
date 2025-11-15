package model

// Volunteer stores volunteer-specific profile attributes.
// Соответствует таблице volunteers.
type Volunteer struct {
	ID           int64
	About        *string
	SearchRadius *int32
	CategoryIDs  []int32
}
