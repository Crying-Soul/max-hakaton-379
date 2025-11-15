package handler

import (
	"fmt"
	"net/url"
	"strings"

	"maxBot/internal/fsm"
)

// EncodePayload кодирует transition и параметры в строку формата "<transition>?<params>"
// Например: "5?user_id=123&role=organizer"
func EncodePayload(transition fsm.Transition, params map[string]string) string {
	transitionStr := transition.String()

	if len(params) == 0 {
		return transitionStr
	}

	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	return transitionStr + "?" + values.Encode()
}

// DecodePayload декодирует строку формата "<transition>?<params>" в transition и параметры
// Возвращает Transition, map параметров и ошибку, если декодирование не удалось
func DecodePayload(payload string) (fsm.Transition, map[string]string, error) {
	// Разделяем на transition и параметры
	parts := strings.SplitN(payload, "?", 2)

	// Парсим transition
	transition, err := fsm.ParseTransition(parts[0])
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse transition: %w", err)
	}

	// Если параметров нет, возвращаем пустую мапу
	if len(parts) == 1 {
		return transition, make(map[string]string), nil
	}

	// Парсим параметры
	values, err := url.ParseQuery(parts[1])
	if err != nil {
		return transition, nil, fmt.Errorf("failed to parse params: %w", err)
	}

	// Преобразуем url.Values в map[string]string
	params := make(map[string]string)
	for key, vals := range values {
		if len(vals) > 0 {
			params[key] = vals[0] // Берём первое значение
		}
	}

	return transition, params, nil
}
