package service

import (
	"context"
	"time"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

// EventService aggregates event-related database operations.
type EventService interface {
	CreateEvent(ctx context.Context, title string, date time.Time, location string, locationLat, locationLon float64, maxVolunteers int32, organizerID int64) (model.Event, error)
	UpdateEvent(ctx context.Context, id int32, title string, description *string, chat *int64, date time.Time, durationHours *int32, location string, locationLat, locationLon float64, categoryID *int32, contacts *string, maxVolunteers int32) (model.Event, error)
	UpdateEventStatus(ctx context.Context, id int32, status string) (model.Event, error)
	CancelEvent(ctx context.Context, id int32, reason *string) (model.Event, error)
	CompleteEvent(ctx context.Context, id int32) (model.Event, error)
	DeleteEvent(ctx context.Context, id int32) error
	GetEventByID(ctx context.Context, id int32) (model.Event, error)
	IncrementEventVolunteers(ctx context.Context, id int32, delta int32) (int32, error)
	CountAvailableEventsForVolunteer(ctx context.Context, volunteerID int64) (int64, error)
	ListAvailableEventsForVolunteer(ctx context.Context, volunteerID int64, limit, offset int32) ([]model.Event, error)
	ListEvents(ctx context.Context, limit, offset int32) ([]model.Event, error)
	ListEventsByOrganizer(ctx context.Context, organizerID int64, limit, offset int32) ([]model.Event, error)
	ListEventsByStatus(ctx context.Context, status string, limit, offset int32) ([]model.Event, error)
	ListUpcomingEvents(ctx context.Context, limit, offset int32) ([]model.Event, error)
	ListEventsByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]model.Event, error)
	ListEventsNearLocation(ctx context.Context, lat, lon, radiusKm float64, limit, offset int32) ([]model.Event, error)
	ListEventsForVolunteer(ctx context.Context, volunteerID int64, limit, offset int32) ([]model.Event, error)
	ListEventsWithPendingApplications(ctx context.Context, limit, offset int32) ([]model.Event, error)
	ListEventsForMap(ctx context.Context, params ListMapEventsParams) ([]model.MapEvent, error)
	ListEventsForMapByVolunteer(ctx context.Context, params ListMapEventsForVolunteerParams) ([]model.MapEvent, error)
}

type eventService struct {
	q dbsqlc.Querier
}

// ListMapEventsParams описывает фильтры для REST API карты волонтёров.
type ListMapEventsParams struct {
	Offset      int32
	Limit       int32
	Lat         float64
	Lon         float64
	RadiusKm    float64
	CategoryIDs []int32
}

type ListMapEventsForVolunteerParams struct {
	ListMapEventsParams
	VolunteerID int64
}

func NewEventService(q dbsqlc.Querier) EventService {
	return &eventService{q: q}
}

func (s *eventService) CreateEvent(ctx context.Context, title string, date time.Time, location string, locationLat, locationLon float64, maxVolunteers int32, organizerID int64) (model.Event, error) {
	params := dbsqlc.CreateEventParams{
		Title:             title,
		Description:       stringPtrToText(nil),
		Chat:              int64PtrToInt8(nil),
		Date:              timePtrToTimestamp(&date),
		DurationHours:     int32PtrToInt4(nil),
		Location:          location,
		LocationLat:       float64ToNumeric(locationLat),
		LocationLon:       float64ToNumeric(locationLon),
		CategoryID:        int32PtrToInt4(nil),
		OrganizerID:       int64ToInt8(organizerID),
		Contacts:          stringPtrToText(nil),
		MaxVolunteers:     maxVolunteers,
		CurrentVolunteers: nil,
		Status:            nil, // будет "planned" по умолчанию
		CancelledReason:   stringPtrToText(nil),
		CompletedAt:       timePtrToTimestamp(nil),
	}
	e, err := s.q.CreateEvent(ctx, params)
	if err != nil {
		return model.Event{}, err
	}
	return mapEvent(e)
}

func (s *eventService) UpdateEvent(ctx context.Context, id int32, title string, description *string, chat *int64, date time.Time, durationHours *int32, location string, locationLat, locationLon float64, categoryID *int32, contacts *string, maxVolunteers int32) (model.Event, error) {
	params := dbsqlc.UpdateEventParams{
		ID:            id,
		Title:         title,
		Description:   stringPtrToText(description),
		Chat:          int64PtrToInt8(chat),
		Date:          timePtrToTimestamp(&date),
		DurationHours: int32PtrToInt4(durationHours),
		Location:      location,
		LocationLat:   float64ToNumeric(locationLat),
		LocationLon:   float64ToNumeric(locationLon),
		CategoryID:    int32PtrToInt4(categoryID),
		Contacts:      stringPtrToText(contacts),
		MaxVolunteers: maxVolunteers,
	}
	e, err := s.q.UpdateEvent(ctx, params)
	if err != nil {
		return model.Event{}, err
	}
	return mapEvent(e)
}

func (s *eventService) UpdateEventStatus(ctx context.Context, id int32, status string) (model.Event, error) {
	params := dbsqlc.UpdateEventStatusParams{
		ID:     id,
		Status: stringToText(status),
	}
	e, err := s.q.UpdateEventStatus(ctx, params)
	if err != nil {
		return model.Event{}, err
	}
	return mapEvent(e)
}

func (s *eventService) CancelEvent(ctx context.Context, id int32, reason *string) (model.Event, error) {
	params := dbsqlc.CancelEventParams{
		ID:     id,
		Reason: stringPtrToText(reason),
	}
	e, err := s.q.CancelEvent(ctx, params)
	if err != nil {
		return model.Event{}, err
	}
	return mapEvent(e)
}

func (s *eventService) CompleteEvent(ctx context.Context, id int32) (model.Event, error) {
	e, err := s.q.CompleteEvent(ctx, id)
	if err != nil {
		return model.Event{}, err
	}
	return mapEvent(e)
}

func (s *eventService) DeleteEvent(ctx context.Context, id int32) error {
	return s.q.DeleteEvent(ctx, id)
}

func (s *eventService) GetEventByID(ctx context.Context, id int32) (model.Event, error) {
	e, err := s.q.GetEventByID(ctx, id)
	if err != nil {
		return model.Event{}, err
	}
	return mapEvent(e)
}

func (s *eventService) IncrementEventVolunteers(ctx context.Context, id int32, delta int32) (int32, error) {
	params := dbsqlc.IncrementEventVolunteersParams{
		ID:    id,
		Delta: int32ToInt4(delta),
	}
	val, err := s.q.IncrementEventVolunteers(ctx, params)
	if err != nil {
		return 0, err
	}
	return int4Value(val), nil
}

func (s *eventService) CountAvailableEventsForVolunteer(ctx context.Context, volunteerID int64) (int64, error) {
	return s.q.CountAvailableEventsForVolunteer(ctx, int64ToInt8(volunteerID))
}

func (s *eventService) ListAvailableEventsForVolunteer(ctx context.Context, volunteerID int64, limit, offset int32) ([]model.Event, error) {
	params := dbsqlc.ListAvailableEventsForVolunteerParams{
		VolunteerID: int64ToInt8(volunteerID),
		Limit:       limit,
		Offset:      offset,
	}
	items, err := s.q.ListAvailableEventsForVolunteer(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEvents(items)
}

func (s *eventService) ListEvents(ctx context.Context, limit, offset int32) ([]model.Event, error) {
	params := dbsqlc.ListEventsParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListEvents(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEvents(items)
}

func (s *eventService) ListEventsByOrganizer(ctx context.Context, organizerID int64, limit, offset int32) ([]model.Event, error) {
	params := dbsqlc.ListEventsByOrganizerParams{
		OrganizerID: int64ToInt8(organizerID),
		Limit:       limit,
		Offset:      offset,
	}
	items, err := s.q.ListEventsByOrganizer(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEvents(items)
}

func (s *eventService) ListEventsByStatus(ctx context.Context, status string, limit, offset int32) ([]model.Event, error) {
	params := dbsqlc.ListEventsByStatusParams{
		Status: stringPtrToText(&status),
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListEventsByStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEvents(items)
}

func (s *eventService) ListUpcomingEvents(ctx context.Context, limit, offset int32) ([]model.Event, error) {
	params := dbsqlc.ListUpcomingEventsParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListUpcomingEvents(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEvents(items)
}

func (s *eventService) ListEventsByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]model.Event, error) {
	params := dbsqlc.ListEventsByCategoryParams{
		CategoryID: int32PtrToInt4(&categoryID),
		Limit:      limit,
		Offset:     offset,
	}
	items, err := s.q.ListEventsByCategory(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEvents(items)
}

func (s *eventService) ListEventsNearLocation(ctx context.Context, lat, lon, radiusKm float64, limit, offset int32) ([]model.Event, error) {
	delta := radiusKm / 111.0
	params := dbsqlc.ListEventsNearLocationParams{
		TargetLat: float64ToNumeric(lat),
		LatDelta:  float64ToNumeric(delta),
		TargetLon: float64ToNumeric(lon),
		LonDelta:  float64ToNumeric(delta),
		Limit:     limit,
		Offset:    offset,
	}
	items, err := s.q.ListEventsNearLocation(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEvents(items)
}

func (s *eventService) ListEventsForVolunteer(ctx context.Context, volunteerID int64, limit, offset int32) ([]model.Event, error) {
	params := dbsqlc.ListEventsForVolunteerParams{
		VolunteerID: int64ToInt8(volunteerID),
		Limit:       limit,
		Offset:      offset,
	}
	items, err := s.q.ListEventsForVolunteer(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEvents(items)
}

func (s *eventService) ListEventsWithPendingApplications(ctx context.Context, limit, offset int32) ([]model.Event, error) {
	params := dbsqlc.ListEventsWithPendingApplicationsParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListEventsWithPendingApplications(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEvents(items)
}

func (s *eventService) ListEventsForMap(ctx context.Context, params ListMapEventsParams) ([]model.MapEvent, error) {
	rows, err := s.q.ListEventsForMap(ctx, dbsqlc.ListEventsForMapParams{
		Offset:      params.Offset,
		Limit:       params.Limit,
		Lat:         params.Lat,
		Lon:         params.Lon,
		CategoryIds: params.CategoryIDs,
		RadiusKm:    params.RadiusKm,
	})
	if err != nil {
		return nil, err
	}
	return mapMapEvents(rows)
}

func (s *eventService) ListEventsForMapByVolunteer(ctx context.Context, params ListMapEventsForVolunteerParams) ([]model.MapEvent, error) {
	rows, err := s.q.ListEventsForMapByVolunteer(ctx, dbsqlc.ListEventsForMapByVolunteerParams{
		Offset:      params.Offset,
		Limit:       params.Limit,
		VolunteerID: int64ToInt8(params.VolunteerID),
		Lat:         params.Lat,
		Lon:         params.Lon,
		CategoryIds: params.CategoryIDs,
		RadiusKm:    params.RadiusKm,
	})
	if err != nil {
		return nil, err
	}
	return mapMapEventsWithStatus(rows)
}

var _ EventService = (*eventService)(nil)
