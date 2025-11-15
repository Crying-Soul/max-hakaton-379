package handler

import (
	"context"
	"fmt"
	"strings"

	"maxBot/internal/di"
	"maxBot/internal/fsm"

	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

// ReplyVerificationHandler принимает сообщения организатора для ответа администратору.
type ReplyVerificationHandler struct {
	services *di.Services
}

func NewReplyVerificationHandler(services *di.Services) *ReplyVerificationHandler {
	return &ReplyVerificationHandler{services: services}
}

func (h *ReplyVerificationHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	text := "Напишите сообщение для администратора. Если заявка ещё на проверке, мы приложим комментарий к ней."
	keyboard := &maxbot.Keyboard{}
	payloadBack := EncodePayload(fsm.ReplyVerificationToVerification, nil)
	keyboard.AddRow().AddCallback("← Назад", schemes.NEGATIVE, payloadBack)

	return h.sendMessage(ctx, update, text, keyboard)
}

func (h *ReplyVerificationHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
	case *schemes.MessageCreatedUpdate:
		text := strings.TrimSpace(upd.Message.Body.Text)
		if text == "" {
			return fsm.Error, nil, fmt.Errorf("сообщение не может быть пустым")
		}
		organizerID := update.GetUserID()
		if _, err := h.services.OrganizerService.GetOrganizer(ctx, organizerID); err != nil {
			return fsm.Error, nil, fmt.Errorf("раздел доступен только организаторам")
		}
		history, err := h.services.OrganizerService.ListOrganizerVerificationHistory(ctx, organizerID, 1, 0)
		if err != nil {
			history = nil
		}
		if _, err := h.services.OrganizerService.CreateOrganizerVerificationRequest(ctx, organizerID, "pending", &text); err != nil {
			return fsm.Error, nil, fmt.Errorf("не удалось сохранить сообщение. Попробуйте позже")
		}
		hadPending := len(history) > 0 && strings.EqualFold(history[0].Status, "pending")
		action := "history"
		if !hadPending {
			action = "new"
		}
		return fsm.ReplyVerificationToVerification, map[string]string{
			"notice": "Сообщение отправлено администратору",
			"action": action,
		}, nil
	case *schemes.MessageCallbackUpdate:
		event, params, err := DecodePayload(upd.Callback.Payload)
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("неверный callback")
		}
		if !containsTransition(availableTransitions, event.String()) {
			return fsm.Error, nil, fmt.Errorf("действие недоступно")
		}
		return event, params, nil
	default:
		return fsm.Error, nil, fmt.Errorf("введите текстовое сообщение")
	}
}

func (h *ReplyVerificationHandler) sendMessage(ctx context.Context, update schemes.UpdateInterface, text string, keyboard *maxbot.Keyboard) error {
	msg := maxbot.NewMessage().SetUser(update.GetUserID()).SetText(text)
	if keyboard != nil {
		msg.AddKeyboard(keyboard)
	}
	_, err := h.services.API.Messages.Send(ctx, msg)
	return err
}
