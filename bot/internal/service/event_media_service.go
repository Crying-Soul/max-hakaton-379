package service

import (
	"context"
	"time"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

// EventMediaService handles event media records.
type EventMediaService interface {
	AddEventMedia(ctx context.Context, eventID *int32, token string, uploadedAt *time.Time, uploadedBy *int64) (model.EventMedia, error)
	DeleteEventMedia(ctx context.Context, id int32) error
	DeleteEventMediaByEvent(ctx context.Context, eventID *int32) error
	GetEventMediaByID(ctx context.Context, id int32) (model.EventMedia, error)
	GetEventMediaByToken(ctx context.Context, token string) (model.EventMedia, error)
	ListEventMedia(ctx context.Context, eventID *int32, limit, offset int32) ([]model.EventMedia, error)
	ListEventMediaByUploader(ctx context.Context, uploadedBy *int64, limit, offset int32) ([]model.EventMedia, error)
}

type eventMediaService struct {
	q dbsqlc.Querier
}

func NewEventMediaService(q dbsqlc.Querier) EventMediaService {
	return &eventMediaService{q: q}
}

func (s *eventMediaService) AddEventMedia(ctx context.Context, eventID *int32, token string, uploadedAt *time.Time, uploadedBy *int64) (model.EventMedia, error) {
	params := dbsqlc.AddEventMediaParams{
		EventID:    int32PtrToInt4(eventID),
		Token:      token,
		UploadedAt: uploadedAt,
		UploadedBy: int64PtrToInt8(uploadedBy),
	}
	m, err := s.q.AddEventMedia(ctx, params)
	if err != nil {
		return model.EventMedia{}, err
	}
	return mapEventMedium(m), nil
}

func (s *eventMediaService) DeleteEventMedia(ctx context.Context, id int32) error {
	return s.q.DeleteEventMedia(ctx, id)
}

func (s *eventMediaService) DeleteEventMediaByEvent(ctx context.Context, eventID *int32) error {
	return s.q.DeleteEventMediaByEvent(ctx, int32PtrToInt4(eventID))
}

func (s *eventMediaService) GetEventMediaByID(ctx context.Context, id int32) (model.EventMedia, error) {
	m, err := s.q.GetEventMediaByID(ctx, id)
	if err != nil {
		return model.EventMedia{}, err
	}
	return mapEventMedium(m), nil
}

func (s *eventMediaService) GetEventMediaByToken(ctx context.Context, token string) (model.EventMedia, error) {
	m, err := s.q.GetEventMediaByToken(ctx, token)
	if err != nil {
		return model.EventMedia{}, err
	}
	return mapEventMedium(m), nil
}

func (s *eventMediaService) ListEventMedia(ctx context.Context, eventID *int32, limit, offset int32) ([]model.EventMedia, error) {
	params := dbsqlc.ListEventMediaParams{
		EventID: int32PtrToInt4(eventID),
		Limit:   limit,
		Offset:  offset,
	}
	items, err := s.q.ListEventMedia(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEventMedia(items), nil
}

func (s *eventMediaService) ListEventMediaByUploader(ctx context.Context, uploadedBy *int64, limit, offset int32) ([]model.EventMedia, error) {
	params := dbsqlc.ListEventMediaByUploaderParams{
		UploadedBy: int64PtrToInt8(uploadedBy),
		Limit:      limit,
		Offset:     offset,
	}
	items, err := s.q.ListEventMediaByUploader(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapEventMedia(items), nil
}

var _ EventMediaService = (*eventMediaService)(nil)
