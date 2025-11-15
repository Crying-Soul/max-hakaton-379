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

type MainMenuHandler struct {
	services *di.Services
}

func NewMainMenuHandler(services *di.Services) *MainMenuHandler {
	return &MainMenuHandler{services: services}
}

func (h *MainMenuHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	keyboard := &maxbot.Keyboard{}

	keyboard.AddRow().AddCallback("Мои события", schemes.DEFAULT, fsm.MainMenuToEvents.String())

	events := EncodePayload(fsm.MainMenuToEvents, map[string]string{"page": "1"})
	keyboard.AddRow().AddCallback("События", schemes.DEFAULT, events)

	keyboard.AddRow().AddCallback("Заявки", schemes.DEFAULT, fsm.MainMenuToApplications.String())
	keyboard.AddRow().AddCallback("О себе", schemes.DEFAULT, fsm.MainMenuToAbout.String())
	keyboard.AddRow().AddCallback("Верификация", schemes.DEFAULT, fsm.MainMenuToVerifications.String())
	keyboard.AddRow().AddCallback("Назад", schemes.DEFAULT, fsm.MainMenuToSelectRole.String())

	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText("Главное меню:").
		AddKeyboard(keyboard)

	err := h.services.API.Messages.EditMessage(ctx, update.(*schemes.MessageCallbackUpdate).Message.Body.Mid, msg)
	return err
}

func (h *MainMenuHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
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
