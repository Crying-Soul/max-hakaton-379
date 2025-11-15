package service

import (
	"context"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

// VolunteerApplicationService manages volunteer application workflow.
type VolunteerApplicationService interface {
	CreateVolunteerApplication(ctx context.Context, eventID int32, volunteerID int64) (model.VolunteerApplication, error)
	UpsertVolunteerApplication(ctx context.Context, application model.VolunteerApplication) (model.VolunteerApplication, error)
	DeleteVolunteerApplication(ctx context.Context, id int32) error
	GetVolunteerApplication(ctx context.Context, eventID *int32, volunteerID *int64) (model.VolunteerApplication, error)
	GetVolunteerApplicationByID(ctx context.Context, id int32) (model.VolunteerApplication, error)
	ListApplicationsByEvent(ctx context.Context, eventID *int32, limit, offset int32) ([]model.VolunteerApplication, error)
	ListApplicationsByStatus(ctx context.Context, status string, limit, offset int32) ([]model.VolunteerApplication, error)
	ListApplicationsByVolunteer(ctx context.Context, volunteerID *int64, limit, offset int32) ([]model.VolunteerApplication, error)
	ListApplicationsForOrganizer(ctx context.Context, organizerID *int64, limit, offset int32) ([]model.VolunteerApplication, error)
	ListPendingApplicationsByEvent(ctx context.Context, eventID *int32, limit, offset int32) ([]model.VolunteerApplication, error)
	UpdateVolunteerApplicationStatus(ctx context.Context, id int32, status string, rejectionReason *string, reviewedBy *int64) (model.VolunteerApplication, error)
	ResetVolunteerApplicationReview(ctx context.Context, id int32) (model.VolunteerApplication, error)
}

type volunteerApplicationService struct {
	q dbsqlc.Querier
}

func NewVolunteerApplicationService(q dbsqlc.Querier) VolunteerApplicationService {
	return &volunteerApplicationService{q: q}
}

func (s *volunteerApplicationService) CreateVolunteerApplication(ctx context.Context, eventID int32, volunteerID int64) (model.VolunteerApplication, error) {
	params := dbsqlc.CreateVolunteerApplicationParams{
		EventID:         int32ToInt4(eventID),
		VolunteerID:     int64ToInt8(volunteerID),
		Status:          nil, // будет "pending" по умолчанию
		RejectionReason: stringPtrToText(nil),
		ReviewedBy:      int64PtrToInt8(nil),
		ReviewedAt:      timePtrToTimestamp(nil),
	}
	a, err := s.q.CreateVolunteerApplication(ctx, params)
	if err != nil {
		return model.VolunteerApplication{}, err
	}
	return mapVolunteerApplication(a), nil
}

func (s *volunteerApplicationService) UpsertVolunteerApplication(ctx context.Context, application model.VolunteerApplication) (model.VolunteerApplication, error) {
	params := dbsqlc.UpsertVolunteerApplicationParams{
		EventID:         int32PtrToInt4(application.EventID),
		VolunteerID:     int64PtrToInt8(application.VolunteerID),
		Status:          application.Status,
		RejectionReason: stringPtrToText(application.RejectionReason),
		ReviewedBy:      int64PtrToInt8(application.ReviewedBy),
		ReviewedAt:      timePtrToTimestamp(application.ReviewedAt),
	}
	a, err := s.q.UpsertVolunteerApplication(ctx, params)
	if err != nil {
		return model.VolunteerApplication{}, err
	}
	return mapVolunteerApplication(a), nil
}

func (s *volunteerApplicationService) DeleteVolunteerApplication(ctx context.Context, id int32) error {
	return s.q.DeleteVolunteerApplication(ctx, id)
}

func (s *volunteerApplicationService) GetVolunteerApplication(ctx context.Context, eventID *int32, volunteerID *int64) (model.VolunteerApplication, error) {
	params := dbsqlc.GetVolunteerApplicationParams{
		EventID:     int32PtrToInt4(eventID),
		VolunteerID: int64PtrToInt8(volunteerID),
	}
	a, err := s.q.GetVolunteerApplication(ctx, params)
	if err != nil {
		return model.VolunteerApplication{}, err
	}
	return mapVolunteerApplication(a), nil
}

func (s *volunteerApplicationService) GetVolunteerApplicationByID(ctx context.Context, id int32) (model.VolunteerApplication, error) {
	a, err := s.q.GetVolunteerApplicationByID(ctx, id)
	if err != nil {
		return model.VolunteerApplication{}, err
	}
	return mapVolunteerApplication(a), nil
}

func (s *volunteerApplicationService) ListApplicationsByEvent(ctx context.Context, eventID *int32, limit, offset int32) ([]model.VolunteerApplication, error) {
	params := dbsqlc.ListApplicationsByEventParams{
		EventID: int32PtrToInt4(eventID),
		Limit:   limit,
		Offset:  offset,
	}
	items, err := s.q.ListApplicationsByEvent(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapVolunteerApplications(items), nil
}

func (s *volunteerApplicationService) ListApplicationsByStatus(ctx context.Context, status string, limit, offset int32) ([]model.VolunteerApplication, error) {
	params := dbsqlc.ListApplicationsByStatusParams{
		Status: stringPtrToText(&status),
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListApplicationsByStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapVolunteerApplications(items), nil
}

func (s *volunteerApplicationService) ListApplicationsByVolunteer(ctx context.Context, volunteerID *int64, limit, offset int32) ([]model.VolunteerApplication, error) {
	params := dbsqlc.ListApplicationsByVolunteerParams{
		VolunteerID: int64PtrToInt8(volunteerID),
		Limit:       limit,
		Offset:      offset,
	}
	items, err := s.q.ListApplicationsByVolunteer(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapVolunteerApplications(items), nil
}

func (s *volunteerApplicationService) ListApplicationsForOrganizer(ctx context.Context, organizerID *int64, limit, offset int32) ([]model.VolunteerApplication, error) {
	params := dbsqlc.ListApplicationsForOrganizerParams{
		OrganizerID: int64PtrToInt8(organizerID),
		Limit:       limit,
		Offset:      offset,
	}
	items, err := s.q.ListApplicationsForOrganizer(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapVolunteerApplications(items), nil
}

func (s *volunteerApplicationService) ListPendingApplicationsByEvent(ctx context.Context, eventID *int32, limit, offset int32) ([]model.VolunteerApplication, error) {
	params := dbsqlc.ListPendingApplicationsByEventParams{
		EventID: int32PtrToInt4(eventID),
		Limit:   limit,
		Offset:  offset,
	}
	items, err := s.q.ListPendingApplicationsByEvent(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapVolunteerApplications(items), nil
}

func (s *volunteerApplicationService) UpdateVolunteerApplicationStatus(ctx context.Context, id int32, status string, rejectionReason *string, reviewedBy *int64) (model.VolunteerApplication, error) {
	params := dbsqlc.UpdateVolunteerApplicationStatusParams{
		ID:              id,
		Status:          stringToText(status),
		RejectionReason: stringPtrToText(rejectionReason),
		ReviewedBy:      int64PtrToInt8(reviewedBy),
	}
	a, err := s.q.UpdateVolunteerApplicationStatus(ctx, params)
	if err != nil {
		return model.VolunteerApplication{}, err
	}
	return mapVolunteerApplication(a), nil
}

func (s *volunteerApplicationService) ResetVolunteerApplicationReview(ctx context.Context, id int32) (model.VolunteerApplication, error) {
	a, err := s.q.ResetVolunteerApplicationReview(ctx, id)
	if err != nil {
		return model.VolunteerApplication{}, err
	}
	return mapVolunteerApplication(a), nil
}

var _ VolunteerApplicationService = (*volunteerApplicationService)(nil)
