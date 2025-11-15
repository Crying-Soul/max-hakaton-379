package handler

import (
	"context"
	"fmt"
	"maxBot/internal/di"
	"maxBot/internal/fsm"
	"slices"
	"strconv"

	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

type EventsHandler struct {
	services *di.Services
}

func NewEventsHandler(services *di.Services) *EventsHandler {
	return &EventsHandler{services: services}
}

// EnterState обрабатывает вход в состояние (вызывается после перехода)
func (h *EventsHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	keyboard := &maxbot.Keyboard{}

	keyboard.AddRow().
		AddCallback("Фильтр по категориям", schemes.DEFAULT, fsm.EventsToCategoriesFilter.String()).
		AddCallback("Фильтр по геолокации", schemes.DEFAULT, fsm.EventsToGeoFilter.String())

	page, err := strconv.Atoi(params["page"])
	if err != nil {
		return err
	}

	limit := int32(8)
	offset := int32(page-1) * limit // страницы начинаются с 1, поэтому offset = (page-1) * limit

	events, _ := h.services.EventService.ListAvailableEventsForVolunteer(ctx, update.GetUserID(), limit, offset)

	count, err := h.services.EventService.CountAvailableEventsForVolunteer(ctx, update.GetUserID())
	if err != nil {
		return err
	}

	for _, event := range events {
		id := strconv.Itoa(int(event.ID))
		eventPayload := EncodePayload(fsm.EventsToEvent, map[string]string{"id": id})
		keyboard.AddRow().AddCallback(event.Title, schemes.DEFAULT, eventPayload)
	}

	totalPages := int((count + int64(limit) - 1) / int64(limit)) // ceil division

	row := keyboard.AddRow()

	if page > 1 {
		previousPageStr := strconv.Itoa(page - 1)
		previousPayload := EncodePayload(fsm.Loop, map[string]string{"page": previousPageStr})
		row.AddCallback("<<", schemes.DEFAULT, previousPayload)
	}

	row.AddCallback(params["page"], schemes.DEFAULT, fsm.Loop.String())

	if page < totalPages {
		nextPageStr := strconv.Itoa(page + 1)
		nextPayload := EncodePayload(fsm.Loop, map[string]string{"page": nextPageStr})
		row.AddCallback(">>", schemes.DEFAULT, nextPayload)
	}
	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText("События:").
		AddKeyboard(keyboard)

	err = h.services.API.Messages.EditMessage(ctx, update.(*schemes.MessageCallbackUpdate).Message.Body.Mid, msg)
	return err
}

// LeaveState проверяет апдейт и возвращает событие для выхода из состояния и опциональную ошибку
func (h *EventsHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
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
