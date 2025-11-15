package service

import (
	"context"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

// CategoryService handles category-related operations.
type CategoryService interface {
	CreateCategory(ctx context.Context, name string) (model.Category, error)
	GetCategory(ctx context.Context, id int32) (model.Category, error)
	GetCategoryByName(ctx context.Context, name string) (model.Category, error)
	UpdateCategory(ctx context.Context, id int32, name string, description *string) (model.Category, error)
	SetCategoryActive(ctx context.Context, id int32, isActive bool) (model.Category, error)
	ListCategories(ctx context.Context, limit, offset int32) ([]model.Category, error)
	ListActiveCategories(ctx context.Context, limit, offset int32) ([]model.Category, error)
	SearchCategories(ctx context.Context, searchQuery string, limit, offset int32) ([]model.Category, error)
}

type categoryService struct {
	q dbsqlc.Querier
}

func NewCategoryService(q dbsqlc.Querier) CategoryService {
	return &categoryService{q: q}
}

func (s *categoryService) CreateCategory(ctx context.Context, name string) (model.Category, error) {
	params := dbsqlc.CreateCategoryParams{
		Name:        name,
		Description: stringPtrToText(nil),
	}
	c, err := s.q.CreateCategory(ctx, params)
	if err != nil {
		return model.Category{}, err
	}
	return mapCategory(c), nil
}

func (s *categoryService) GetCategory(ctx context.Context, id int32) (model.Category, error) {
	c, err := s.q.GetCategory(ctx, id)
	if err != nil {
		return model.Category{}, err
	}
	return mapCategory(c), nil
}

func (s *categoryService) GetCategoryByName(ctx context.Context, name string) (model.Category, error) {
	c, err := s.q.GetCategoryByName(ctx, name)
	if err != nil {
		return model.Category{}, err
	}
	return mapCategory(c), nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id int32, name string, description *string) (model.Category, error) {
	params := dbsqlc.UpdateCategoryParams{
		ID:          id,
		Name:        name,
		Description: stringPtrToText(description),
	}
	c, err := s.q.UpdateCategory(ctx, params)
	if err != nil {
		return model.Category{}, err
	}
	return mapCategory(c), nil
}

func (s *categoryService) SetCategoryActive(ctx context.Context, id int32, isActive bool) (model.Category, error) {
	params := dbsqlc.SetCategoryActiveParams{
		ID:       id,
		IsActive: boolToBool(isActive),
	}
	c, err := s.q.SetCategoryActive(ctx, params)
	if err != nil {
		return model.Category{}, err
	}
	return mapCategory(c), nil
}

func (s *categoryService) ListCategories(ctx context.Context, limit, offset int32) ([]model.Category, error) {
	params := dbsqlc.ListCategoriesParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListCategories(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapCategories(items), nil
}

func (s *categoryService) ListActiveCategories(ctx context.Context, limit, offset int32) ([]model.Category, error) {
	params := dbsqlc.ListActiveCategoriesParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListActiveCategories(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapCategories(items), nil
}

func (s *categoryService) SearchCategories(ctx context.Context, searchQuery string, limit, offset int32) ([]model.Category, error) {
	params := dbsqlc.SearchCategoriesParams{
		Query:  searchQuery,
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.SearchCategories(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapCategories(items), nil
}

var _ CategoryService = (*categoryService)(nil)
