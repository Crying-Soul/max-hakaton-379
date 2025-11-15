package service

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

func mapUser(u dbsqlc.User) (model.User, error) {
	lat, err := numericToFloat64Ptr(u.LocationLat)
	if err != nil {
		return model.User{}, fmt.Errorf("map user location_lat: %w", err)
	}
	lon, err := numericToFloat64Ptr(u.LocationLon)
	if err != nil {
		return model.User{}, fmt.Errorf("map user location_lon: %w", err)
	}
	return model.User{
		ID:          u.ID,
		Username:    textToPtr(u.Username),
		Name:        u.Name,
		Role:        u.Role,
		State:       u.State,
		IsBlocked:   boolValue(u.IsBlocked),
		CreatedAt:   timestampToTime(u.CreatedAt),
		UpdatedAt:   timestampToTime(u.UpdatedAt),
		LocationLat: lat,
		LocationLon: lon,
	}, nil
}

func mapUsers(items []dbsqlc.User) ([]model.User, error) {
	result := make([]model.User, 0, len(items))
	for _, item := range items {
		mapped, err := mapUser(item)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func mapAdmin(a dbsqlc.Admin) model.Admin {
	return model.Admin{
		ID:        a.ID,
		CreatedAt: timestampToTime(a.CreatedAt),
	}
}

func mapVolunteer(v dbsqlc.Volunteer) model.Volunteer {
	return model.Volunteer{
		ID:           v.ID,
		About:        textToPtr(v.About),
		SearchRadius: int4ToPtr(v.SearchRadius),
		CategoryIDs:  copyInt32s(v.CategoryIds),
	}
}

func mapVolunteers(items []dbsqlc.Volunteer) []model.Volunteer {
	result := make([]model.Volunteer, 0, len(items))
	for _, v := range items {
		result = append(result, mapVolunteer(v))
	}
	return result
}

// NOTE: ранее здесь были mapVolunteerWithUser* функции, возвращающие составные модели.
// По новому требованию модели теперь только по одной структуре на файл/таблицу,
// поэтому маппинг volunteer+user следует делать там, где он действительно нужен,
// напрямую через sqlc-строки.

func mapOrganizer(o dbsqlc.Organizer) (model.Organizer, error) {
	verifiedAt := timestampToPtr(o.VerifiedAt)
	verifiedBy := int8ToPtr(o.VerifiedBy)
	return model.Organizer{
		ID:               o.ID,
		OrganizationName: o.OrganizationName,
		About:            textToPtr(o.About),
		VerifiedAt:       verifiedAt,
		VerifiedBy:       verifiedBy,
		CreatedAt:        timestampToTime(o.CreatedAt),
		UpdatedAt:        timestampToTime(o.UpdatedAt),
	}, nil
}

func mapOrganizers(items []dbsqlc.Organizer) ([]model.Organizer, error) {
	result := make([]model.Organizer, 0, len(items))
	for _, item := range items {
		mapped, err := mapOrganizer(item)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

// Аналогично, совместные структуры OrganizerWithUser больше не мапим здесь.

func mapOrganizerVerificationRequest(item dbsqlc.OrganizerVerificationRequest) model.OrganizerVerificationRequest {
	return model.OrganizerVerificationRequest{
		ID:               item.ID,
		OrganizerID:      item.OrganizerID,
		Status:           item.Status,
		OrganizerComment: textToPtr(item.OrganizerComment),
		AdminComment:     textToPtr(item.AdminComment),
		ReviewedBy:       int8ToPtr(item.ReviewedBy),
		SubmittedAt:      timestampToTime(item.SubmittedAt),
		ReviewedAt:       timestampToPtr(item.ReviewedAt),
	}
}

func mapOrganizerVerificationRequests(items []dbsqlc.OrganizerVerificationRequest) []model.OrganizerVerificationRequest {
	result := make([]model.OrganizerVerificationRequest, 0, len(items))
	for _, item := range items {
		result = append(result, mapOrganizerVerificationRequest(item))
	}
	return result
}

func mapCategory(c dbsqlc.Category) model.Category {
	return model.Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: textToPtr(c.Description),
		IsActive:    boolPtr(c.IsActive),
		CreatedAt:   timestampToTime(c.CreatedAt),
	}
}

func mapCategories(items []dbsqlc.Category) []model.Category {
	result := make([]model.Category, 0, len(items))
	for _, item := range items {
		result = append(result, mapCategory(item))
	}
	return result
}

func mapEvent(e dbsqlc.Event) (model.Event, error) {
	locationLat, err := numericToFloat64(e.LocationLat)
	if err != nil {
		return model.Event{}, fmt.Errorf("map event location_lat: %w", err)
	}
	locationLon, err := numericToFloat64(e.LocationLon)
	if err != nil {
		return model.Event{}, fmt.Errorf("map event location_lon: %w", err)
	}
	return model.Event{
		ID:                e.ID,
		Title:             e.Title,
		Description:       textToPtr(e.Description),
		Chat:              int8ToPtr(e.Chat),
		Date:              timestampToTime(e.Date),
		DurationHours:     int4ToPtr(e.DurationHours),
		Location:          e.Location,
		LocationLat:       locationLat,
		LocationLon:       locationLon,
		CategoryID:        int4ToPtr(e.CategoryID),
		OrganizerID:       int8ToPtr(e.OrganizerID),
		Contacts:          textToPtr(e.Contacts),
		MaxVolunteers:     e.MaxVolunteers,
		CurrentVolunteers: int4ToPtr(e.CurrentVolunteers),
		Status:            textToPtr(e.Status),
		CancelledReason:   textToPtr(e.CancelledReason),
		CompletedAt:       timestampToPtr(e.CompletedAt),
		CreatedAt:         timestampToTime(e.CreatedAt),
		UpdatedAt:         timestampToTime(e.UpdatedAt),
	}, nil
}

func mapEvents(items []dbsqlc.Event) ([]model.Event, error) {
	result := make([]model.Event, 0, len(items))
	for _, event := range items {
		mapped, err := mapEvent(event)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func mapEventParticipant(ep dbsqlc.EventParticipant) model.EventParticipant {
	return model.EventParticipant{
		ID:            ep.ID,
		EventID:       int4ToPtr(ep.EventID),
		VolunteerID:   int8ToPtr(ep.VolunteerID),
		ApplicationID: int4ToPtr(ep.ApplicationID),
		JoinedChatAt:  timestampToTime(ep.JoinedChatAt),
	}
}

func mapEventParticipants(items []dbsqlc.EventParticipant) []model.EventParticipant {
	result := make([]model.EventParticipant, 0, len(items))
	for _, ep := range items {
		result = append(result, mapEventParticipant(ep))
	}
	return result
}

// EventParticipantWithUser больше не используется, так как доменная модель только EventParticipant.

func mapEventMedium(m dbsqlc.EventMedium) model.EventMedia {
	return model.EventMedia{
		ID:         m.ID,
		EventID:    int4ToPtr(m.EventID),
		Token:      m.Token,
		UploadedAt: timestampToTime(m.UploadedAt),
		UploadedBy: int8ToPtr(m.UploadedBy),
	}
}

func mapEventMedia(items []dbsqlc.EventMedium) []model.EventMedia {
	result := make([]model.EventMedia, 0, len(items))
	for _, item := range items {
		result = append(result, mapEventMedium(item))
	}
	return result
}

func mapMapEvent(e dbsqlc.ListEventsForMapRow) (model.MapEvent, error) {
	locationLat, err := numericToFloat64(e.LocationLat)
	if err != nil {
		return model.MapEvent{}, fmt.Errorf("map map event location_lat: %w", err)
	}
	locationLon, err := numericToFloat64(e.LocationLon)
	if err != nil {
		return model.MapEvent{}, fmt.Errorf("map map event location_lon: %w", err)
	}
	return model.MapEvent{
		ID:                e.ID,
		Title:             e.Title,
		Description:       textToPtr(e.Description),
		Chat:              int8ToPtr(e.Chat),
		Date:              timestampToTime(e.Date),
		DurationHours:     int4ToPtr(e.DurationHours),
		Location:          e.Location,
		LocationLat:       locationLat,
		LocationLon:       locationLon,
		CategoryID:        int4ToPtr(e.CategoryID),
		CategoryName:      textToPtr(e.CategoryName),
		OrganizerID:       int8ToPtr(e.OrganizerID),
		Contacts:          textToPtr(e.Contacts),
		MaxVolunteers:     e.MaxVolunteers,
		CurrentVolunteers: e.CurrentVolunteers,
		SlotsLeft:         e.SlotsLeft,
		DistanceKm:        e.DistanceKm,
		Status:            textToPtr(e.Status),
		CancelledReason:   textToPtr(e.CancelledReason),
		CompletedAt:       timestampToPtr(e.CompletedAt),
		CreatedAt:         timestampToTime(e.CreatedAt),
		UpdatedAt:         timestampToTime(e.UpdatedAt),
	}, nil
}

func mapMapEvents(items []dbsqlc.ListEventsForMapRow) ([]model.MapEvent, error) {
	result := make([]model.MapEvent, 0, len(items))
	for _, item := range items {
		mapped, err := mapMapEvent(item)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func mapMapEventWithStatus(e dbsqlc.ListEventsForMapByVolunteerRow) (model.MapEvent, error) {
	base := dbsqlc.ListEventsForMapRow{
		ID:                e.ID,
		Title:             e.Title,
		Description:       e.Description,
		Date:              e.Date,
		DurationHours:     e.DurationHours,
		Location:          e.Location,
		LocationLat:       e.LocationLat,
		LocationLon:       e.LocationLon,
		CategoryID:        e.CategoryID,
		OrganizerID:       e.OrganizerID,
		Contacts:          e.Contacts,
		Chat:              e.Chat,
		MaxVolunteers:     e.MaxVolunteers,
		CurrentVolunteers: e.CurrentVolunteers,
		Status:            e.Status,
		CancelledReason:   e.CancelledReason,
		CompletedAt:       e.CompletedAt,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
		CategoryName:      e.CategoryName,
		SlotsLeft:         e.SlotsLeft,
		DistanceKm:        e.DistanceKm,
	}
	mapped, err := mapMapEvent(base)
	if err != nil {
		return model.MapEvent{}, err
	}
	mapped.ApplicationStatus = textToPtr(e.ApplicationStatus)
	return mapped, nil
}

func mapMapEventsWithStatus(items []dbsqlc.ListEventsForMapByVolunteerRow) ([]model.MapEvent, error) {
	result := make([]model.MapEvent, 0, len(items))
	for _, item := range items {
		mapped, err := mapMapEventWithStatus(item)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func mapVolunteerApplication(a dbsqlc.VolunteerApplication) model.VolunteerApplication {
	status := textToPtr(a.Status)
	rejectionReason := textToPtr(a.RejectionReason)
	return model.VolunteerApplication{
		ID:              a.ID,
		EventID:         int4ToPtr(a.EventID),
		VolunteerID:     int8ToPtr(a.VolunteerID),
		Status:          status,
		RejectionReason: rejectionReason,
		ReviewedBy:      int8ToPtr(a.ReviewedBy),
		AppliedAt:       timestampToTime(a.AppliedAt),
		ReviewedAt:      timestampToPtr(a.ReviewedAt),
	}
}

func mapVolunteerApplications(items []dbsqlc.VolunteerApplication) []model.VolunteerApplication {
	result := make([]model.VolunteerApplication, 0, len(items))
	for _, item := range items {
		result = append(result, mapVolunteerApplication(item))
	}
	return result
}

func textToPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

func boolPtr(b pgtype.Bool) *bool {
	if !b.Valid {
		return nil
	}
	v := b.Bool
	return &v
}

func boolValue(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

func timestampToTime(ts pgtype.Timestamp) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time
}

func timestampToPtr(ts pgtype.Timestamp) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}

func int4ToPtr(v pgtype.Int4) *int32 {
	if !v.Valid {
		return nil
	}
	val := v.Int32
	return &val
}

func int8ToPtr(v pgtype.Int8) *int64 {
	if !v.Valid {
		return nil
	}
	val := v.Int64
	return &val
}

func int4Value(v pgtype.Int4) int32 {
	if !v.Valid {
		return 0
	}
	return v.Int32
}

func int64FromInt8(v pgtype.Int8) int64 {
	if !v.Valid {
		return 0
	}
	return v.Int64
}

func numericToFloat64Ptr(n pgtype.Numeric) (*float64, error) {
	if !n.Valid {
		return nil, nil
	}
	val, err := n.Float64Value()
	if err != nil {
		return nil, err
	}
	if !val.Valid {
		return nil, nil
	}
	out := val.Float64
	return &out, nil
}

func numericToFloat64(n pgtype.Numeric) (float64, error) {
	if !n.Valid {
		return 0, fmt.Errorf("numeric value is null")
	}
	val, err := n.Float64Value()
	if err != nil {
		return 0, err
	}
	if !val.Valid {
		return 0, fmt.Errorf("numeric value is invalid")
	}
	return val.Float64, nil
}

func copyInt32s(src []int32) []int32 {
	if len(src) == 0 {
		return nil
	}
	dst := make([]int32, len(src))
	copy(dst, src)
	return dst
}

// Reverse mapping helpers: from model types to pgtype

func stringPtrToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func stringToText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func float64PtrToNumeric(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	return float64ToNumeric(*f)
}

func float64ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(f)
	return n
}

func boolToBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func int32ToInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func int32PtrToInt4(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *i, Valid: true}
}

func int64ToInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}

func int64PtrToInt8(i *int64) pgtype.Int8 {
	if i == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *i, Valid: true}
}

func timePtrToTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: *t, Valid: true}
}
