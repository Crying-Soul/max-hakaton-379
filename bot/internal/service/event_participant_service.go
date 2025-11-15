package service

import (
	"context"
	"time"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

// EventParticipantService manages participant records.
type EventParticipantService interface {
	AddEventParticipant(ctx context.Context, eventID *int32, volunteerID *int64, applicationID *int32, joinedChatAt *time.Time) (model.EventParticipant, error)
	RemoveEventParticipant(ctx context.Context, eventID *int32, volunteerID *int64) error
	DeleteParticipantsByEvent(ctx context.Context, eventID *int32) error
	GetEventParticipant(ctx context.Context, eventID *int32, volunteerID *int64) (model.EventParticipant, error)
	ListEventParticipants(ctx context.Context, eventID *int32, limit, offset int32) ([]model.EventParticipant, error)
	ListParticipantEvents(ctx context.Context, volunteerID *int64, limit, offset int32) ([]model.EventParticipant, error)
	CountParticipantsForEvent(ctx context.Context, eventID *int32) (int64, error)
}

type eventParticipantService struct {
	q dbsqlc.Querier
}

func NewEventParticipantService(q dbsqlc.Querier) EventParticipantService {
	return &eventParticipantService{q: q}
}

func (s *eventParticipantService) AddEventParticipant(ctx context.Context, eventID *int32, volunteerID *int64, applicationID *int32, joinedChatAt *time.Time) (model.EventParticipant, error) {
	params := dbsqlc.AddEventParticipantParams{
		EventID:       int32PtrToInt4(eventID),
		VolunteerID:   int64PtrToInt8(volunteerID),
		ApplicationID: int32PtrToInt4(applicationID),
		JoinedChatAt:  joinedChatAt,
	}
	ep, err := s.q.AddEventParticipant(ctx, params)
	if err != nil {
		return model.EventParticipant{}, err
	}
	return mapEventParticipant(ep), nil
}

func (s *eventParticipantService) RemoveEventParticipant(ctx context.Context, eventID *int32, volunteerID *int64) error {
	params := dbsqlc.RemoveEventParticipantParams{
		EventID:     int32PtrToInt4(eventID),
		VolunteerID: int64PtrToInt8(volunteerID),
	}
	return s.q.RemoveEventParticipant(ctx, params)
}

func (s *eventParticipantService) DeleteParticipantsByEvent(ctx context.Context, eventID *int32) error {
	return s.q.DeleteParticipantsByEvent(ctx, int32PtrToInt4(eventID))
}

func (s *eventParticipantService) GetEventParticipant(ctx context.Context, eventID *int32, volunteerID *int64) (model.EventParticipant, error) {
	params := dbsqlc.GetEventParticipantParams{
		EventID:     int32PtrToInt4(eventID),
		VolunteerID: int64PtrToInt8(volunteerID),
	}
	ep, err := s.q.GetEventParticipant(ctx, params)
	if err != nil {
		return model.EventParticipant{}, err
	}
	return mapEventParticipant(ep), nil
}

func (s *eventParticipantService) ListEventParticipants(ctx context.Context, eventID *int32, limit, offset int32) ([]model.EventParticipant, error) {
	params := dbsqlc.ListEventParticipantsParams{
		EventID: int32PtrToInt4(eventID),
		Limit:   limit,
		Offset:  offset,
	}
	items, err := s.q.ListEventParticipants(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEventParticipants(items), nil
}

func (s *eventParticipantService) ListParticipantEvents(ctx context.Context, volunteerID *int64, limit, offset int32) ([]model.EventParticipant, error) {
	params := dbsqlc.ListParticipantEventsParams{
		VolunteerID: int64PtrToInt8(volunteerID),
		Limit:       limit,
		Offset:      offset,
	}
	items, err := s.q.ListParticipantEvents(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEventParticipants(items), nil
}

func (s *eventParticipantService) CountParticipantsForEvent(ctx context.Context, eventID *int32) (int64, error) {
	return s.q.CountParticipantsForEvent(ctx, int32PtrToInt4(eventID))
}

var _ EventParticipantService = (*eventParticipantService)(nil)
