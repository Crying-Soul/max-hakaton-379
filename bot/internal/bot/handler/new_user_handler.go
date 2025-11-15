package handler

import (
	"context"
	"fmt"

	"maxBot/internal/di"
	"maxBot/internal/fsm"

	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

// NewUserHandler обрабатывает состояние NewUser
type NewUserHandler struct {
	services *di.Services
}

// NewNewUserHandler создаёт новый хендлер для состояния NewUser
func NewNewUserHandler(services *di.Services) *NewUserHandler {
	return &NewUserHandler{services: services}
}

// EnterState отправляет приветственное сообщение новому пользователю
func (h *NewUserHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	keyboard := &maxbot.Keyboard{}
	payload := EncodePayload(fsm.NewUserToSelectRole, nil)
	keyboard.AddRow().AddCallback("Начать", schemes.POSITIVE, payload)

	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText("Добро пожаловать! Нажмите кнопку для начала работы.").
		AddKeyboard(keyboard)

	_, err := h.services.API.Messages.Send(ctx, msg)
	return err
}

// LeaveState проверяет апдейт и возвращает событие для выхода из состояния NewUser
func (h *NewUserHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
	case *schemes.MessageCallbackUpdate:
		event, params, err := DecodePayload(upd.Callback.Payload)
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("неправильный callback")
		}
		return event, params, nil
	}

	return fsm.Error, nil, fmt.Errorf("пожалуйста, нажмите 'Начать'")
}
