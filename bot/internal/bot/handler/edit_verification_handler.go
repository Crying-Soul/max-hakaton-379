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

// EditVerificationHandler позволяет организатору обновить профиль и отправить новую заявку.
type EditVerificationHandler struct {
	services *di.Services
}

func NewEditVerificationHandler(services *di.Services) *EditVerificationHandler {
	return &EditVerificationHandler{services: services}
}

func (h *EditVerificationHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	organizer, err := h.services.OrganizerService.GetOrganizer(ctx, update.GetUserID())
	if err != nil {
		return h.sendMessage(ctx, update, "Раздел доступен только организаторам. Попросите администратора назначить вам роль организатора.", nil)
	}

	var builder strings.Builder
	builder.WriteString("Текущие данные организации:\n")
	builder.WriteString("Название: ")
	if organizer.OrganizationName == "" {
		builder.WriteString("—")
	} else {
		builder.WriteString(organizer.OrganizationName)
	}
	builder.WriteString("\nКонтакты: ")
	// if organizer.Contacts == nil || *organizer.Contacts == "" {
	// 	builder.WriteString("—")
	// } else {
	// 	builder.WriteString(*organizer.Contacts)
	// }
	builder.WriteString("\n\nОтправьте сообщение в формате:\n")
	builder.WriteString("1 строка — название организации.\n")
	builder.WriteString("Остальной текст — контакты, ссылки и комментарии.\n")
	builder.WriteString("Если заявка уже на проверке, мы просто обновим её данные и передадим админу актуальную версию.")

	keyboard := &maxbot.Keyboard{}
	keyboard.AddRow().AddCallback("← Назад", schemes.NEGATIVE, EncodePayload(fsm.EditVerificationToVerification, nil))

	return h.sendMessage(ctx, update, builder.String(), keyboard)
}

func (h *EditVerificationHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
	case *schemes.MessageCreatedUpdate:
		text := strings.TrimSpace(upd.Message.Body.Text)
		if text == "" {
			return fsm.Error, nil, fmt.Errorf("сообщение не может быть пустым")
		}

		organizer, err := h.services.OrganizerService.GetOrganizer(ctx, update.GetUserID())
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("раздел доступен только организаторам")
		}

		history, err := h.services.OrganizerService.ListOrganizerVerificationHistory(ctx, organizer.ID, 1, 0)
		if err != nil {
			history = nil
		}

		lines := strings.Split(text, "\n")
		newName := strings.TrimSpace(lines[0])
		if newName == "" {
			newName = strings.TrimSpace(organizer.OrganizationName)
		}
		if newName == "" {
			return fsm.Error, nil, fmt.Errorf("укажите название организации в первой строке")
		}

		var contactsPtr *string
		if len(lines) > 1 {
			rest := strings.TrimSpace(strings.Join(lines[1:], "\n"))
			if rest != "" {
				contactsPtr = &rest
			}
		}

		if _, err := h.services.OrganizerService.UpdateOrganizerProfile(ctx, organizer.ID, newName, contactsPtr); err != nil {
			return fsm.Error, nil, fmt.Errorf("не удалось обновить профиль: %w", err)
		}

		if _, err := h.services.OrganizerService.CreateOrganizerVerificationRequest(ctx, organizer.ID, "pending", &text); err != nil {
			return fsm.Error, nil, fmt.Errorf("не удалось создать заявку: %w", err)
		}

		hadPending := len(history) > 0 && strings.EqualFold(history[0].Status, "pending")
		notice := "Заявка отправлена на проверку"
		action := "new"
		if hadPending {
			notice = "Заявка обновлена"
			action = "history"
		}
		return fsm.EditVerificationToVerification, map[string]string{
			"notice": notice,
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
		return fsm.Error, nil, fmt.Errorf("введите текстовое сообщение с данными организации")
	}
}

func (h *EditVerificationHandler) sendMessage(ctx context.Context, update schemes.UpdateInterface, text string, keyboard *maxbot.Keyboard) error {
	msg := maxbot.NewMessage().SetUser(update.GetUserID()).SetText(text)
	if keyboard != nil {
		msg.AddKeyboard(keyboard)
	}
	_, err := h.services.API.Messages.Send(ctx, msg)
	return err
}
