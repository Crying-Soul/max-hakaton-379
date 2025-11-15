package handler

import (
	"context"
	"fmt"
	"maxBot/internal/di"
	"maxBot/internal/fsm"
	"maxBot/internal/model"
	"slices"
	"strconv"

	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

type EventHandler struct {
	services *di.Services
}

func NewEventHandler(services *di.Services) *EventHandler {
	return &EventHandler{services: services}
}

func (h *EventHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	idStr, ok := params["id"]
	if !ok {
		return fmt.Errorf("event id not provided")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid event id")
	}

	event, err := h.services.EventService.GetEventByID(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	user, err := h.services.UserService.GetUserByID(ctx, update.GetUserID())
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user has applied
	var application *model.VolunteerApplication
	if user.Role == "volunteer" {
		app, err := h.services.ApplicationService.GetVolunteerApplication(ctx, &event.ID, &user.ID)
		if err == nil {
			application = &app
		}
	}

	// Handle actions
	if action := params["action"]; action != "" {
		switch action {
		case "apply":
			if user.Role != "volunteer" {
				return fmt.Errorf("only volunteers can apply")
			}
			if application != nil {
				return fmt.Errorf("you have already applied")
			}
			_, err := h.services.ApplicationService.CreateVolunteerApplication(ctx, event.ID, user.ID)
			if err != nil {
				return fmt.Errorf("failed to apply: %w", err)
			}
			// Send success message
			h.services.API.Messages.Send(ctx, maxbot.NewMessage().
				SetUser(update.GetUserID()).
				SetText("Вы успешно подали заявку на участие в событии!"))
			// Re-fetch application
			app, _ := h.services.ApplicationService.GetVolunteerApplication(ctx, &event.ID, &user.ID)
			application = &app
		case "cancel":
			if application == nil {
				return fmt.Errorf("you have not applied")
			}
			err := h.services.ApplicationService.DeleteVolunteerApplication(ctx, application.ID)
			if err != nil {
				return fmt.Errorf("failed to cancel application: %w", err)
			}
			// Send success message
			h.services.API.Messages.Send(ctx, maxbot.NewMessage().
				SetUser(update.GetUserID()).
				SetText("Заявка отменена."))
			application = nil
		}
	}

	// Build message text
	text := fmt.Sprintf("**%s**\n\n", event.Title)
	if event.Description != nil {
		text += fmt.Sprintf("Описание: %s\n\n", *event.Description)
	}
	text += fmt.Sprintf("Дата: %s\n", event.Date.Format("02.01.2006 15:04"))
	if event.DurationHours != nil {
		text += fmt.Sprintf("Длительность: %d часов\n", *event.DurationHours)
	}
	text += fmt.Sprintf("Место: %s\n", event.Location)
	text += fmt.Sprintf("Максимум волонтёров: %d\n", event.MaxVolunteers)
	if event.CurrentVolunteers != nil {
		text += fmt.Sprintf("Текущих волонтёров: %d\n", *event.CurrentVolunteers)
	}
	if event.Contacts != nil {
		text += fmt.Sprintf("Контакты: %s\n", *event.Contacts)
	}
	if event.Status != nil {
		text += fmt.Sprintf("Статус: %s\n", *event.Status)
	}

	keyboard := &maxbot.Keyboard{}

	if user.Role == "volunteer" {
		if application == nil {
			applyPayload := EncodePayload(fsm.Loop, map[string]string{"id": idStr, "action": "apply"})
			keyboard.AddRow().AddCallback("Подать заявку", schemes.DEFAULT, applyPayload)
		} else {
			cancelPayload := EncodePayload(fsm.Loop, map[string]string{"id": idStr, "action": "cancel"})
			keyboard.AddRow().AddCallback("Отменить заявку", schemes.DEFAULT, cancelPayload)
		}
	}

	backPayload := EncodePayload(fsm.EventToEvents, map[string]string{"page": "1"})
	keyboard.AddRow().AddCallback("Назад", schemes.DEFAULT, backPayload)

	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText(text).
		AddKeyboard(keyboard)

	err = h.services.API.Messages.EditMessage(ctx, update.(*schemes.MessageCallbackUpdate).Message.Body.Mid, msg)
	return err
}

func (h *EventHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
	case *schemes.MessageCallbackUpdate:
		event, params, err := DecodePayload(upd.Callback.Payload)
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("неверный callback")
		}
		if event == fsm.Loop {
			return fsm.Loop, params, nil
		}
		if !slices.Contains(availableTransitions, event.String()) {
			return fsm.Error, nil, fmt.Errorf("неверный ответ, воспользуйтесь кнопками")
		}
		return event, params, nil
	}
	return fsm.Error, nil, fmt.Errorf("неверный ответ")
}
