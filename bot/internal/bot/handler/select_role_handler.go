package handler

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"maxBot/internal/di"
	"maxBot/internal/fsm"

	"github.com/jackc/pgx/v5"
	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

// NewUserHandler обрабатывает состояние NewUser
type SelectRoleHandler struct {
	services *di.Services
}

// NewNewUserHandler создаёт новый хендлер для состояния NewUser
func NewSelectRoleHandler(services *di.Services) *SelectRoleHandler {
	return &SelectRoleHandler{services: services}
}

// EnterState отправляет приветственное сообщение новому пользователю
func (h *SelectRoleHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	keyboard := h.services.API.Messages.NewKeyboardBuilder()
	volunteerPayload := EncodePayload(fsm.SelectRoleToMainMenu, map[string]string{"role": "volunteer"})
	organizerPayload := EncodePayload(fsm.SelectRoleToMainMenu, map[string]string{"role": "organizer"})
	_, err := h.services.AdminService.GetAdmin(ctx, update.GetUserID())
	if err == nil {
		adminPayload := EncodePayload(fsm.SelectRoleToMainMenu, map[string]string{"role": "admin"})
		keyboard.AddRow().AddCallback("Администратор", schemes.DEFAULT, adminPayload)
	}

	keyboard.AddRow().AddCallback("Волонтер", schemes.DEFAULT, volunteerPayload)
	keyboard.AddRow().AddCallback("Организатор", schemes.DEFAULT, organizerPayload)

	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText("Выберите вашу роль:").
		AddKeyboard(keyboard)

	err = h.services.API.Messages.EditMessage(ctx, update.(*schemes.MessageCallbackUpdate).Message.Body.Mid, msg)
	return err
}

// LeaveState проверяет апдейт и возвращает событие для выхода из состояния NewUser
func (h *SelectRoleHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
	case *schemes.MessageCallbackUpdate:
		event, params, err := DecodePayload(upd.Callback.Payload)
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("неверный callback")
		}
		if !slices.Contains(availableTransitions, event.String()) {
			return fsm.Error, nil, fmt.Errorf("неверный ответ, воспользуйтесь кнопками")
		}
		role := params["role"]
		if role == "" {
			return fsm.Error, nil, fmt.Errorf("роль не указана в параметрах")
		}
		_, err = h.services.UserService.UpdateUserRole(ctx, update.GetUserID(), role)
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("failed to update user role: %w", err)
		}
		if role == "organizer" {
			if err := h.ensureOrganizer(ctx, update.GetUserID()); err != nil {
				return fsm.Error, nil, err
			}
		}
		return event, params, nil
	}
	return fsm.Error, nil, fmt.Errorf("неверный ответ")
}

func (h *SelectRoleHandler) ensureOrganizer(ctx context.Context, userID int64) error {
	_, err := h.services.OrganizerService.GetOrganizer(ctx, userID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("не удалось проверить статус организатора: %w", err)
	}
	user, err := h.services.UserService.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("не удалось получить данные пользователя: %w", err)
	}
	name := strings.TrimSpace(user.Name)
	if name == "" {
		name = fmt.Sprintf("Организация %d", userID)
	}
	if _, err := h.services.OrganizerService.CreateOrganizer(ctx, userID, name); err != nil {
		return fmt.Errorf("не удалось создать запись организатора: %w", err)
	}
	return nil
}
