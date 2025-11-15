package handler

import (
	"context"

	"maxBot/internal/fsm"

	"github.com/rectid/max-bot-api-client-go/schemes"
)

// Handler интерфейс для всех state-хендлеров
type Handler interface {
	// EnterState обрабатывает вход в состояние (вызывается после перехода)
	EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error

	// LeaveState проверяет апдейт и возвращает событие для выхода из состояния, параметры и опциональную ошибку
	LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error)
}
