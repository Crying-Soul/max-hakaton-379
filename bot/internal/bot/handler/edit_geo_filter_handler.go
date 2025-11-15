package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"maxBot/internal/di"
	"maxBot/internal/fsm"

	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

type EditGeoFilterHandler struct {
	services *di.Services
}

func NewEditGeoFilterHandler(services *di.Services) *EditGeoFilterHandler {
	return &EditGeoFilterHandler{services: services}
}

func (h *EditGeoFilterHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	vol, err := h.services.VolunteerService.GetVolunteer(ctx, update.GetUserID())
	if err != nil {
		msg := maxbot.NewMessage().SetUser(update.GetUserID()).SetText("Раздел доступен только волонтёрам. Попросите администратора назначить вам роль волонтёра.")
		_, sendErr := h.services.API.Messages.Send(ctx, msg)
		return sendErr
	}

	radiusText := "не задан"
	if vol.SearchRadius != nil {
		radiusText = fmt.Sprintf("%d км", *vol.SearchRadius)
	}

	keyboard := &maxbot.Keyboard{}
	keyboard.AddRow().AddCallback("← Назад", schemes.NEGATIVE, EncodePayload(fsm.EditGeoFilterToGeoFilter, nil))

	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText(fmt.Sprintf("Текущий радиус поиска: %s\n\nВведите новый радиус поиска в км (число от 1 до 100):", radiusText)).
		AddKeyboard(keyboard)

	switch upd := update.(type) {
	case *schemes.MessageCallbackUpdate:
		return h.services.API.Messages.EditMessage(ctx, upd.Message.Body.Mid, msg)
	}
	return nil
}

func (h *EditGeoFilterHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
	case *schemes.MessageCreatedUpdate:
		text := strings.TrimSpace(upd.Message.Body.Text)
		if text == "" {
			return fsm.Error, nil, fmt.Errorf("радиус не может быть пустым")
		}

		radius, err := strconv.Atoi(text)
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("радиус должен быть числом")
		}

		if radius < 1 || radius > 100 {
			return fsm.Error, nil, fmt.Errorf("радиус должен быть от 1 до 100 км")
		}

		vol, err := h.services.VolunteerService.GetVolunteer(ctx, update.GetUserID())
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("раздел доступен только волонтёрам")
		}

		radiusInt32 := int32(radius)
		if _, err := h.services.VolunteerService.UpdateVolunteerSearchRadius(ctx, vol.ID, &radiusInt32); err != nil {
			return fsm.Error, nil, fmt.Errorf("не удалось обновить радиус поиска: %w", err)
		}

		return fsm.EditGeoFilterToGeoFilter, map[string]string{
			"notice": fmt.Sprintf("Радиус поиска обновлён до %d км", radius),
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
		return fsm.Error, nil, fmt.Errorf("введите число для радиуса поиска")
	}
}
