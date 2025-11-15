package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	maxbot "github.com/rectid/max-bot-api-client-go"

	dbsqlc "maxBot/internal/db/sqlc"
	"maxBot/internal/model"
)

// OrganizerService provides business-level access to organizer data.
type OrganizerService interface {
	CreateOrganizer(ctx context.Context, id int64, organizationName string) (model.Organizer, error)
	UpsertOrganizer(ctx context.Context, organizer model.Organizer) (model.Organizer, error)
	DeleteOrganizer(ctx context.Context, id int64) error
	GetOrganizer(ctx context.Context, id int64) (model.Organizer, error)
	UpdateOrganizerProfile(ctx context.Context, id int64, organizationName string, about *string) (model.Organizer, error)
	SetOrganizerVerification(ctx context.Context, id int64, verifiedAt *time.Time, verifiedBy *int64) (model.Organizer, error)
	ListOrganizers(ctx context.Context, limit, offset int32) ([]model.Organizer, error)
	ListVerifiedOrganizers(ctx context.Context, limit, offset int32) ([]model.Organizer, error)
	ListUnverifiedOrganizers(ctx context.Context, limit, offset int32) ([]model.Organizer, error)
	ListOrganizerVerificationHistory(ctx context.Context, organizerID int64, limit, offset int32) ([]model.OrganizerVerificationRequest, error)
	CreateOrganizerVerificationRequest(ctx context.Context, organizerID int64, status string, comment *string) (model.OrganizerVerificationRequest, error)
}

type organizerService struct {
	q   dbsqlc.Querier
	api *maxbot.Api
}

func NewOrganizerService(q dbsqlc.Querier, api *maxbot.Api) OrganizerService {
	return &organizerService{q: q, api: api}
}

func (s *organizerService) CreateOrganizer(ctx context.Context, id int64, organizationName string) (model.Organizer, error) {
	params := dbsqlc.CreateOrganizerParams{
		ID:               id,
		OrganizationName: organizationName,
		About:            stringPtrToText(nil),
		VerifiedAt:       timePtrToTimestamp(nil),
		VerifiedBy:       int64PtrToInt8(nil),
	}
	o, err := s.q.CreateOrganizer(ctx, params)
	if err != nil {
		return model.Organizer{}, err
	}
	return mapOrganizer(o)
}

func (s *organizerService) UpsertOrganizer(ctx context.Context, organizer model.Organizer) (model.Organizer, error) {
	params := dbsqlc.UpsertOrganizerParams{
		ID:               organizer.ID,
		OrganizationName: organizer.OrganizationName,
		About:            stringPtrToText(organizer.About),
		VerifiedAt:       timePtrToTimestamp(organizer.VerifiedAt),
		VerifiedBy:       int64PtrToInt8(organizer.VerifiedBy),
	}
	o, err := s.q.UpsertOrganizer(ctx, params)
	if err != nil {
		return model.Organizer{}, err
	}
	return mapOrganizer(o)
}

func (s *organizerService) DeleteOrganizer(ctx context.Context, id int64) error {
	return s.q.DeleteOrganizer(ctx, id)
}

func (s *organizerService) GetOrganizer(ctx context.Context, id int64) (model.Organizer, error) {
	o, err := s.q.GetOrganizer(ctx, id)
	if err != nil {
		return model.Organizer{}, err
	}
	return mapOrganizer(o)
}

func (s *organizerService) UpdateOrganizerProfile(ctx context.Context, id int64, organizationName string, about *string) (model.Organizer, error) {
	params := dbsqlc.UpdateOrganizerProfileParams{
		ID:               id,
		OrganizationName: organizationName,
		About:            stringPtrToText(about),
	}
	o, err := s.q.UpdateOrganizerProfile(ctx, params)
	if err != nil {
		return model.Organizer{}, err
	}
	return mapOrganizer(o)
}

func (s *organizerService) SetOrganizerVerification(ctx context.Context, id int64, verifiedAt *time.Time, verifiedBy *int64) (model.Organizer, error) {
	params := dbsqlc.SetOrganizerVerificationParams{
		ID:         id,
		VerifiedAt: timePtrToTimestamp(verifiedAt),
		VerifiedBy: int64PtrToInt8(verifiedBy),
	}
	o, err := s.q.SetOrganizerVerification(ctx, params)
	if err != nil {
		return model.Organizer{}, err
	}

	// Отправляем уведомление о верификации
	if verifiedAt != nil && s.api != nil {
		msg := maxbot.NewMessage().
			SetUser(id).
			SetText("Ваша организация успешно прошла верификацию! Теперь вам доступен полный функционал организатора.")
		if _, err := s.api.Messages.Send(ctx, msg); err != nil {
			return model.Organizer{}, err
		}
	}
	return mapOrganizer(o)
}

func (s *organizerService) ListOrganizers(ctx context.Context, limit, offset int32) ([]model.Organizer, error) {
	params := dbsqlc.ListOrganizersParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListOrganizers(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapOrganizers(items)
}

func (s *organizerService) ListVerifiedOrganizers(ctx context.Context, limit, offset int32) ([]model.Organizer, error) {
	params := dbsqlc.ListVerifiedOrganizersParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListVerifiedOrganizers(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapOrganizers(items)
}

func (s *organizerService) ListUnverifiedOrganizers(ctx context.Context, limit, offset int32) ([]model.Organizer, error) {
	params := dbsqlc.ListUnverifiedOrganizersParams{
		Limit:  limit,
		Offset: offset,
	}
	items, err := s.q.ListUnverifiedOrganizers(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapOrganizers(items)
}

func (s *organizerService) ListOrganizerVerificationHistory(ctx context.Context, organizerID int64, limit, offset int32) ([]model.OrganizerVerificationRequest, error) {
	params := dbsqlc.ListOrganizerVerificationRequestsParams{
		OrganizerID: organizerID,
		Offset:      offset,
		Limit:       limit,
	}
	items, err := s.q.ListOrganizerVerificationRequests(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapOrganizerVerificationRequests(items), nil
}

func (s *organizerService) CreateOrganizerVerificationRequest(ctx context.Context, organizerID int64, status string, comment *string) (model.OrganizerVerificationRequest, error) {
	if status == "" {
		status = "pending"
	}
	commentText := stringPtrToText(comment)
	if pending, err := s.q.GetLatestPendingOrganizerVerificationRequest(ctx, organizerID); err == nil {
		updated, err := s.q.UpdateOrganizerVerificationRequestComment(ctx, dbsqlc.UpdateOrganizerVerificationRequestCommentParams{
			ID:               pending.ID,
			OrganizerComment: commentText,
		})
		if err != nil {
			return model.OrganizerVerificationRequest{}, err
		}
		return mapOrganizerVerificationRequest(updated), nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return model.OrganizerVerificationRequest{}, err
	}
	params := dbsqlc.CreateOrganizerVerificationRequestParams{
		OrganizerID:      organizerID,
		Status:           status,
		OrganizerComment: commentText,
		AdminComment:     stringPtrToText(nil),
		ReviewedBy:       int64PtrToInt8(nil),
		ReviewedAt:       timePtrToTimestamp(nil),
	}
	item, err := s.q.CreateOrganizerVerificationRequest(ctx, params)
	if err != nil {
		return model.OrganizerVerificationRequest{}, err
	}
	return mapOrganizerVerificationRequest(item), nil
}

var _ OrganizerService = (*organizerService)(nil)
