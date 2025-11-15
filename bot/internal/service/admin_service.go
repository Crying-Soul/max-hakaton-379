package service

import (
	"context"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

// AdminService exposes operations for admin records.
type AdminService interface {
	CreateAdmin(ctx context.Context, id int64) (model.Admin, error)
	DeleteAdmin(ctx context.Context, id int64) error
	GetAdmin(ctx context.Context, id int64) (model.Admin, error)
	ListAdmins(ctx context.Context, limit, offset int32) ([]model.Admin, error)
}

type adminService struct {
	q dbsqlc.Querier
}

func NewAdminService(q dbsqlc.Querier) AdminService {
	return &adminService{q: q}
}

func (s *adminService) CreateAdmin(ctx context.Context, id int64) (model.Admin, error) {
	admin, err := s.q.CreateAdmin(ctx, id)
	if err != nil {
		return model.Admin{}, err
	}
	return mapAdmin(admin), nil
}

func (s *adminService) DeleteAdmin(ctx context.Context, id int64) error {
	return s.q.DeleteAdmin(ctx, id)
}

func (s *adminService) GetAdmin(ctx context.Context, id int64) (model.Admin, error) {
	admin, err := s.q.GetAdmin(ctx, id)
	if err != nil {
		return model.Admin{}, err
	}
	return mapAdmin(admin), nil
}

func (s *adminService) ListAdmins(ctx context.Context, limit, offset int32) ([]model.Admin, error) {
	params := dbsqlc.ListAdminsParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListAdmins(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make([]model.Admin, 0, len(items))
	for _, item := range items {
		result = append(result, mapAdmin(item))
	}
	return result, nil
}

var _ AdminService = (*adminService)(nil)
