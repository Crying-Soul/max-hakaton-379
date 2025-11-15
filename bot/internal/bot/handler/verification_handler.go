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

// VerificationHandler показывает детали конкретной заявки и доступные действия.
type VerificationHandler struct {
	services *di.Services
}

func NewVerificationHandler(services *di.Services) *VerificationHandler {
	return &VerificationHandler{services: services}
}

func (h *VerificationHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	userID := update.GetUserID()

	or, err := h.services.OrganizerService.GetOrganizer(ctx, userID)
	if err != nil {
		return h.sendMessage(ctx, update, "Раздел доступен только организаторам.", nil)
	}

	history, err := h.services.OrganizerService.ListOrganizerVerificationHistory(ctx, or.ID, 10, 0)
	if err != nil {
		return h.sendMessage(ctx, update, "Не удалось загрузить данные заявки.", nil)
	}

	notice := params["notice"]
	action := params["action"]

	var builder strings.Builder
	if notice != "" {
		builder.WriteString("✅ ")
		builder.WriteString(notice)
		builder.WriteString("\n\n")
	}

	//status := translateVerificationStatusPtr(or.VerificationStatus)
	builder.WriteString("Текущий статус: ")
	//builder.WriteString(status)
	builder.WriteString("\n\n")

	var latest *model.OrganizerVerificationRequest
	if len(history) > 0 {
		latest = &history[0]
	}
	if latest == nil {
		builder.WriteString("Нет заявок. Нажмите “Заполнить данные”, чтобы отправить первую заявку на проверку.")
	} else {
		hasPending := strings.EqualFold(latest.Status, "pending")
		if action == "new" && hasPending {
			builder.WriteString("У вас уже есть заявка на проверке. Вы можете обновить данные через кнопку “Заполнить данные”.\n\n")
		}
		builder.WriteString("Последняя заявка от ")
		builder.WriteString(latest.SubmittedAt.Format("02.01.2006 15:04"))
		builder.WriteString("\nСтатус: ")
		builder.WriteString(translateVerificationStatus(latest.Status))
		if latest.AdminComment != nil && *latest.AdminComment != "" {
			builder.WriteString("\nКомментарий администратора: ")
			builder.WriteString(*latest.AdminComment)
		}
		if latest.OrganizerComment != nil && *latest.OrganizerComment != "" {
			builder.WriteString("\nВаш комментарий: ")
			builder.WriteString(*latest.OrganizerComment)
		}
	}

	if latest == nil {
		builder.WriteString("\n\nОтправьте данные об организации: название, контакты, ссылки на документы. Нажмите “Заполнить данные”, чтобы оставить сообщение администратору.")
	}

	keyboard := &maxbot.Keyboard{}
	keyboard.AddRow().AddCallback("Заполнить данные", schemes.POSITIVE, EncodePayload(fsm.VerificationToEditVerification, nil))
	keyboard.AddRow().AddCallback("Написать админу", schemes.DEFAULT, EncodePayload(fsm.VerificationToReplyVerification, nil))
	keyboard.AddRow().AddCallback("← История", schemes.NEGATIVE, EncodePayload(fsm.VerificationToVerifications, nil))

	return h.sendMessage(ctx, update, builder.String(), keyboard)
}

func (h *VerificationHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
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

func (h *VerificationHandler) sendMessage(ctx context.Context, update schemes.UpdateInterface, text string, keyboard *maxbot.Keyboard) error {
	msg := maxbot.NewMessage().SetUser(update.GetUserID()).SetText(text)
	if keyboard != nil {
		msg.AddKeyboard(keyboard)
	}
	_, err := h.services.API.Messages.Send(ctx, msg)
	return err
}
