package handler

import (
	"context"
	"fmt"
	"maxBot/internal/di"
	"maxBot/internal/fsm"
	"slices"

	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

type PersonalEventsHandler struct {
	services *di.Services
}

func NewPersonalEventsHandler(services *di.Services) *PersonalEventsHandler {
	return &PersonalEventsHandler{services: services}
}

// EnterState обрабатывает вход в состояние (вызывается после перехода)
func (h *PersonalEventsHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	keyboard := &maxbot.Keyboard{}
	activePayload := EncodePayload(fsm.PersonalEventsToEvents, map[string]string{"filter": "active"})
	completedPayload := EncodePayload(fsm.Loop, map[string]string{"filter": "completed"})
	canceledPayload := EncodePayload(fsm.Loop, map[string]string{"filter": "canceled"})
	rejectedPayload := EncodePayload(fsm.Loop, map[string]string{"filter": "rejected"})

	keyboard.AddRow().AddCallback("Активные события", schemes.DEFAULT, activePayload)
	keyboard.AddRow().AddCallback("Завершенные события", schemes.DEFAULT, completedPayload)
	keyboard.AddRow().AddCallback("Отмененные события", schemes.DEFAULT, canceledPayload)
	keyboard.AddRow().AddCallback("Участие отклонено", schemes.DEFAULT, rejectedPayload)
	keyboard.AddRow().AddCallback("Назад", schemes.DEFAULT, fsm.PersonalEventsToMainMenu.String())

	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText("Мои события:").
		AddKeyboard(keyboard)

	err := h.services.API.Messages.EditMessage(ctx, update.(*schemes.MessageCallbackUpdate).Message.Body.Mid, msg)
	return err
}

// LeaveState проверяет апдейт и возвращает событие для выхода из состояния, параметры и опциональную ошибку
func (h *PersonalEventsHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
	case *schemes.MessageCallbackUpdate:
		event, params, err := DecodePayload(upd.Callback.Payload)
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("неверный callback")
		}
		if !slices.Contains(availableTransitions, event.String()) {
			return fsm.Error, nil, fmt.Errorf("неверный ответ, воспользуйтесь кнопками")
		}
		return event, params, nil
	}
	return fsm.Error, nil, fmt.Errorf("неверный ответ")
}
