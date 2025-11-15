package service

import (
	"context"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

// UserService provides high-level operations over users.
type UserService interface {
	CreateUser(ctx context.Context, id int64, name string) (model.User, error)
	UpsertUser(ctx context.Context, user model.User) (model.User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUserByID(ctx context.Context, id int64) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	ListUsersByRole(ctx context.Context, role string, limit, offset int32) ([]model.User, error)
	ListUsersByState(ctx context.Context, state string, limit, offset int32) ([]model.User, error)
	ListBlockedUsers(ctx context.Context, limit, offset int32) ([]model.User, error)
	SearchUsers(ctx context.Context, searchQuery string, limit, offset int32) ([]model.User, error)
	ListUsersNearLocation(ctx context.Context, lat, lon, radiusKm float64, limit, offset int32) ([]model.User, error)
	ListUsersByIDs(ctx context.Context, ids []int64) ([]model.User, error)
	UpdateUserProfile(ctx context.Context, id int64, username *string, name string) (model.User, error)
	UpdateUserRole(ctx context.Context, id int64, role string) (model.User, error)
	UpdateUserState(ctx context.Context, id int64, state string) (model.User, error)
	UpdateUserLocation(ctx context.Context, id int64, lat, lon *float64) (model.User, error)
	BlockUser(ctx context.Context, id int64) error
	UnblockUser(ctx context.Context, id int64) error
}

type userService struct {
	q dbsqlc.Querier
}

func NewUserService(q dbsqlc.Querier) UserService {
	return &userService{q: q}
}

func (s *userService) CreateUser(ctx context.Context, id int64, name string) (model.User, error) {
	params := dbsqlc.CreateUserParams{
		ID:          id,
		Username:    stringPtrToText(nil),
		Name:        name,
		Role:        "user",
		State:       "new",
		IsBlocked:   false,
		LocationLat: float64PtrToNumeric(nil),
		LocationLon: float64PtrToNumeric(nil),
	}
	u, err := s.q.CreateUser(ctx, params)
	if err != nil {
		return model.User{}, err
	}
	return mapUser(u)
}

func (s *userService) UpsertUser(ctx context.Context, user model.User) (model.User, error) {
	params := dbsqlc.UpsertUserParams{
		ID:          user.ID,
		Username:    stringPtrToText(user.Username),
		Name:        user.Name,
		Role:        user.Role,
		State:       user.State,
		IsBlocked:   user.IsBlocked,
		LocationLat: float64PtrToNumeric(user.LocationLat),
		LocationLon: float64PtrToNumeric(user.LocationLon),
	}
	u, err := s.q.UpsertUser(ctx, params)
	if err != nil {
		return model.User{}, err
	}
	return mapUser(u)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.q.DeleteUser(ctx, id)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (model.User, error) {
	u, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return mapUser(u)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	u, err := s.q.GetUserByUsername(ctx, stringToText(username))
	if err != nil {
		return model.User{}, err
	}
	return mapUser(u)
}

func (s *userService) ListUsersByRole(ctx context.Context, role string, limit, offset int32) ([]model.User, error) {
	params := dbsqlc.ListUsersByRoleParams{
		Role:   role,
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListUsersByRole(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapUsers(items)
}

func (s *userService) ListUsersByState(ctx context.Context, state string, limit, offset int32) ([]model.User, error) {
	params := dbsqlc.ListUsersByStateParams{
		State:  state,
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListUsersByState(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapUsers(items)
}

func (s *userService) ListBlockedUsers(ctx context.Context, limit, offset int32) ([]model.User, error) {
	params := dbsqlc.ListBlockedUsersParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListBlockedUsers(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapUsers(items)
}

func (s *userService) SearchUsers(ctx context.Context, searchQuery string, limit, offset int32) ([]model.User, error) {
	params := dbsqlc.SearchUsersParams{
		Query:  searchQuery,
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.SearchUsers(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapUsers(items)
}

func (s *userService) ListUsersNearLocation(ctx context.Context, lat, lon, radiusKm float64, limit, offset int32) ([]model.User, error) {
	// radiusKm используется для вычисления дельты (примерно 1 градус = 111 км)
	delta := radiusKm / 111.0
	params := dbsqlc.ListUsersNearLocationParams{
		TargetLat: float64ToNumeric(lat),
		LatDelta:  float64ToNumeric(delta),
		TargetLon: float64ToNumeric(lon),
		LonDelta:  float64ToNumeric(delta),
		Limit:     limit,
		Offset:    offset,
	}
	items, err := s.q.ListUsersNearLocation(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapUsers(items)
}

func (s *userService) ListUsersByIDs(ctx context.Context, ids []int64) ([]model.User, error) {
	items, err := s.q.ListUsersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return mapUsers(items)
}

func (s *userService) UpdateUserProfile(ctx context.Context, id int64, username *string, name string) (model.User, error) {
	params := dbsqlc.UpdateUserProfileParams{
		ID:       id,
		Username: stringPtrToText(username),
		Name:     name,
	}
	u, err := s.q.UpdateUserProfile(ctx, params)
	if err != nil {
		return model.User{}, err
	}
	return mapUser(u)
}

func (s *userService) UpdateUserRole(ctx context.Context, id int64, role string) (model.User, error) {
	params := dbsqlc.UpdateUserRoleParams{
		ID:   id,
		Role: role,
	}
	u, err := s.q.UpdateUserRole(ctx, params)
	if err != nil {
		return model.User{}, err
	}
	return mapUser(u)
}

func (s *userService) UpdateUserState(ctx context.Context, id int64, state string) (model.User, error) {
	params := dbsqlc.UpdateUserStateParams{
		ID:    id,
		State: state,
	}
	u, err := s.q.UpdateUserState(ctx, params)
	if err != nil {
		return model.User{}, err
	}
	return mapUser(u)
}

func (s *userService) UpdateUserLocation(ctx context.Context, id int64, lat, lon *float64) (model.User, error) {
	params := dbsqlc.UpdateUserLocationParams{
		ID:          id,
		LocationLat: float64PtrToNumeric(lat),
		LocationLon: float64PtrToNumeric(lon),
	}
	u, err := s.q.UpdateUserLocation(ctx, params)
	if err != nil {
		return model.User{}, err
	}
	return mapUser(u)
}

func (s *userService) BlockUser(ctx context.Context, id int64) error {
	return s.q.BlockUser(ctx, id)
}

func (s *userService) UnblockUser(ctx context.Context, id int64) error {
	return s.q.UnblockUser(ctx, id)
}

var _ UserService = (*userService)(nil)
