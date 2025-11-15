package handler

import (
	"context"
	"fmt"
	"strings"

	"maxBot/internal/di"
	"maxBot/internal/fsm"
	"maxBot/internal/model"

	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

// VerificationsHandler показывает пользователю историю заявок на верификацию.
type VerificationsHandler struct {
	services *di.Services
}

func NewVerificationsHandler(services *di.Services) *VerificationsHandler {
	return &VerificationsHandler{services: services}
}

func (h *VerificationsHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	userID := update.GetUserID()

	organizer, err := h.services.OrganizerService.GetOrganizer(ctx, userID)
	if err != nil {
		return h.sendMessage(ctx, update, "Раздел доступен только организаторам. Попросите администратора назначить вам роль организатора.", nil)
	}

	history, err := h.services.OrganizerService.ListOrganizerVerificationHistory(ctx, organizer.ID, 5, 0)
	if err != nil {
		return h.sendMessage(ctx, update, "Не удалось загрузить историю заявок. Попробуйте позже.", nil)
	}

	//status := translateVerificationStatusPtr(organizer.VerificationStatus)

	var builder strings.Builder
	builder.WriteString("Статус организации: ")
	//builder.WriteString(status)
	builder.WriteString("\n\n")
	// if organizer.RejectionReason != nil && *organizer.RejectionReason != "" {
	// 	builder.WriteString("Причина отклонения: ")
	// 	builder.WriteString(*organizer.RejectionReason)
	// 	builder.WriteString("\n\n")
	// }
	builder.WriteString(formatVerificationHistory(history))

	hasPending := len(history) > 0 && strings.EqualFold(history[0].Status, "pending")
	keyboard := &maxbot.Keyboard{}
	payloadHistory := EncodePayload(fsm.VerificationsToVerification, map[string]string{"action": "history"})
	keyboard.AddRow().AddCallback("Последняя заявка", schemes.DEFAULT, payloadHistory)
	if hasPending {
		payloadEdit := EncodePayload(fsm.VerificationsToVerification, map[string]string{"action": "edit"})
		keyboard.AddRow().AddCallback("Изменить заявку", schemes.POSITIVE, payloadEdit)
	} else {
		payloadNew := EncodePayload(fsm.VerificationsToVerification, map[string]string{"action": "new"})
		keyboard.AddRow().AddCallback("Новая заявка", schemes.POSITIVE, payloadNew)
	}

	return h.sendMessage(ctx, update, builder.String(), keyboard)
}

func (h *VerificationsHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
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
		return fsm.Error, nil, fmt.Errorf("используйте кнопки ниже")
	}
}

func (h *VerificationsHandler) sendMessage(ctx context.Context, update schemes.UpdateInterface, text string, keyboard *maxbot.Keyboard) error {
	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText(text)
	if keyboard != nil {
		msg.AddKeyboard(keyboard)
	}
	_, err := h.services.API.Messages.Send(ctx, msg)
	return err
}

func formatVerificationHistory(history []model.OrganizerVerificationRequest) string {
	if len(history) == 0 {
		return "Нет заявок"
	}
	var builder strings.Builder
	builder.WriteString("Последние заявки:\n")
	for i, item := range history {
		builder.WriteString(fmt.Sprintf("%d. %s — %s\n", i+1, item.SubmittedAt.Format("02.01.2006 15:04"), translateVerificationStatus(item.Status)))
		if item.OrganizerComment != nil && *item.OrganizerComment != "" {
			builder.WriteString("   ✏️ ")
			builder.WriteString(*item.OrganizerComment)
			builder.WriteString("\n")
		}
		if item.AdminComment != nil && *item.AdminComment != "" {
			builder.WriteString("   🧑‍⚖️ ")
			builder.WriteString(*item.AdminComment)
			builder.WriteString("\n")
		}
	}
	return builder.String()
}
