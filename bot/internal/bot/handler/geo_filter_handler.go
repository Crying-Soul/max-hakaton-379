package handler

import (
	"context"
	"fmt"
	"maxBot/internal/di"
	"maxBot/internal/fsm"

	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

type GeoFilterHandler struct {
	services *di.Services
}

func NewGeoFilterHandler(services *di.Services) *GeoFilterHandler {
	return &GeoFilterHandler{services: services}
}

// EnterState обрабатывает вход в состояние (вызывается после перехода)
func (h *GeoFilterHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	keyboard := &maxbot.Keyboard{}
	eventsPayload := EncodePayload(fsm.GeoFilterToEvents, map[string]string{"page": "1"})
	keyboard.AddRow().AddGeolocation("Отправить геолокацию", true)
	keyboard.AddRow().AddCallback("Изменить радиус поиска", schemes.DEFAULT, fsm.GeoFilterToEditGeoFilter.String())
	keyboard.AddRow().AddCallback("Назад", schemes.DEFAULT, eventsPayload)

	vol, err := h.services.VolunteerService.GetVolunteer(ctx, update.GetUserID())
	if err != nil {
		return err
	}
	radiusText := "не задан"
	if vol.SearchRadius != nil {
		radiusText = fmt.Sprintf("%d км", *vol.SearchRadius)
	}

	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText(fmt.Sprintf("Меню геолокации:\nТекущий радиус поиска: %s", radiusText)).
		AddKeyboard(keyboard)

	switch upd := update.(type) {
	case *schemes.MessageCallbackUpdate:
		return h.services.API.Messages.EditMessage(ctx, upd.Message.Body.Mid, msg)
	case *schemes.MessageCreatedUpdate:
		_, err := h.services.API.Messages.Send(ctx, msg)
		return err
	}
	return nil
}

// LeaveState проверяет апдейт и возвращает событие для выхода из состояния, параметры и опциональную ошибку
func (h *GeoFilterHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
	case *schemes.MessageCallbackUpdate:
		// Обработка callback кнопок
		event, params, err := DecodePayload(upd.Callback.Payload)
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("неверный callback")
		}
		if !containsTransition(availableTransitions, event.String()) {
			return fsm.Error, nil, fmt.Errorf("действие недоступно")
		}
		return event, params, nil
	default:
		return fsm.Error, nil, fmt.Errorf("воспользуйтесь кнопками меню")
	}
}
