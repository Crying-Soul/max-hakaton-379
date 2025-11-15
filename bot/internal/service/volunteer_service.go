package service

import (
	"context"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

// VolunteerService wraps SQLC queries with model mapping for volunteers.
type VolunteerService interface {
	CreateVolunteer(ctx context.Context, id int64) (model.Volunteer, error)
	UpsertVolunteer(ctx context.Context, volunteer model.Volunteer) (model.Volunteer, error)
	DeleteVolunteer(ctx context.Context, id int64) error
	GetVolunteer(ctx context.Context, id int64) (model.Volunteer, error)
	UpdateVolunteerProfile(ctx context.Context, id int64, cv *string, searchRadius *int32) (model.Volunteer, error)
	UpdateVolunteerCategories(ctx context.Context, id int64, categoryIDs []int32) (model.Volunteer, error)
	ListVolunteers(ctx context.Context, limit, offset int32) ([]model.Volunteer, error)
	ListVolunteersByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]model.Volunteer, error)
	ListVolunteersByIDs(ctx context.Context, ids []int64) ([]model.Volunteer, error)
}

type volunteerService struct {
	q dbsqlc.Querier
}

func NewVolunteerService(q dbsqlc.Querier) VolunteerService {
	return &volunteerService{q: q}
}

func (s *volunteerService) CreateVolunteer(ctx context.Context, id int64) (model.Volunteer, error) {
	params := dbsqlc.CreateVolunteerParams{
		ID:           id,
		Cv:           stringPtrToText(nil),
		SearchRadius: int32PtrToInt4(nil),
		CategoryIds:  nil,
	}
	v, err := s.q.CreateVolunteer(ctx, params)
	if err != nil {
		return model.Volunteer{}, err
	}
	return mapVolunteer(v), nil
}

func (s *volunteerService) UpsertVolunteer(ctx context.Context, volunteer model.Volunteer) (model.Volunteer, error) {
	params := dbsqlc.UpsertVolunteerParams{
		ID:           volunteer.ID,
		Cv:           stringPtrToText(volunteer.CV),
		SearchRadius: int32PtrToInt4(volunteer.SearchRadius),
		CategoryIds:  volunteer.CategoryIDs,
	}
	v, err := s.q.UpsertVolunteer(ctx, params)
	if err != nil {
		return model.Volunteer{}, err
	}
	return mapVolunteer(v), nil
}

func (s *volunteerService) DeleteVolunteer(ctx context.Context, id int64) error {
	return s.q.DeleteVolunteer(ctx, id)
}

func (s *volunteerService) GetVolunteer(ctx context.Context, id int64) (model.Volunteer, error) {
	v, err := s.q.GetVolunteer(ctx, id)
	if err != nil {
		return model.Volunteer{}, err
	}
	return mapVolunteer(v), nil
}

func (s *volunteerService) UpdateVolunteerProfile(ctx context.Context, id int64, cv *string, searchRadius *int32) (model.Volunteer, error) {
	params := dbsqlc.UpdateVolunteerProfileParams{
		ID:           id,
		Cv:           stringPtrToText(cv),
		SearchRadius: int32PtrToInt4(searchRadius),
	}
	v, err := s.q.UpdateVolunteerProfile(ctx, params)
	if err != nil {
		return model.Volunteer{}, err
	}
	return mapVolunteer(v), nil
}

func (s *volunteerService) UpdateVolunteerCategories(ctx context.Context, id int64, categoryIDs []int32) (model.Volunteer, error) {
	params := dbsqlc.UpdateVolunteerCategoriesParams{
		ID:          id,
		CategoryIds: categoryIDs,
	}
	v, err := s.q.UpdateVolunteerCategories(ctx, params)
	if err != nil {
		return model.Volunteer{}, err
	}
	return mapVolunteer(v), nil
}

func (s *volunteerService) ListVolunteers(ctx context.Context, limit, offset int32) ([]model.Volunteer, error) {
	params := dbsqlc.ListVolunteersParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListVolunteers(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapVolunteers(items), nil
}

func (s *volunteerService) ListVolunteersByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]model.Volunteer, error) {
	params := dbsqlc.ListVolunteersByCategoryParams{
		CategoryIds: []int32{categoryID},
		Limit:       limit,
		Offset:      offset,
	}
	items, err := s.q.ListVolunteersByCategory(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapVolunteers(items), nil
}

func (s *volunteerService) ListVolunteersByIDs(ctx context.Context, ids []int64) ([]model.Volunteer, error) {
	items, err := s.q.ListVolunteersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return mapVolunteers(items), nil
}

var _ VolunteerService = (*volunteerService)(nil)
