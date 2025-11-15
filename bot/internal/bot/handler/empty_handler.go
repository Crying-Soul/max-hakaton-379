package handler

import (
	"context"

	"maxBot/internal/di"
	"maxBot/internal/fsm"

	"github.com/rectid/max-bot-api-client-go/schemes"
)

// EmptyHandler обрабатывает состояние NewUser
type EmptyHandler struct {
	services *di.Services
}

// NewNewUserHandler создаёт новый хендлер для состояния NewUser
func NewEmptyHandler(services *di.Services) *EmptyHandler {
	return &EmptyHandler{services: services}
}

// EnterState отправляет приветственное сообщение новому пользователю
func (h *EmptyHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {

	return nil
}

// LeaveState проверяет апдейт и возвращает событие для выхода из состояния NewUser
func (h *EmptyHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	return fsm.EmptyToNewUser, nil, nil
}
