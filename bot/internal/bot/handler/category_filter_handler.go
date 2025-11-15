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

type CategoryFilterHandler struct {
	services *di.Services
}

func NewCategoryFilterHandler(services *di.Services) *CategoryFilterHandler {
	return &CategoryFilterHandler{services: services}
}

// EnterState обрабатывает вход в состояние (вызывается после перехода)
func (h *CategoryFilterHandler) EnterState(ctx context.Context, update schemes.UpdateInterface, transition fsm.Transition, params map[string]string) error {
	keyboard := &maxbot.Keyboard{}

	// Получаем текущую страницу
	pageStr := params["page"]
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit := int32(8)
	offset := int32(page-1) * limit

	// Получаем волонтера, чтобы узнать выбранные категории
	volunteer, err := h.services.VolunteerService.GetVolunteer(ctx, update.GetUserID())
	if err != nil {
		return err
	}

	// Получаем список активных категорий
	categories, err := h.services.CategoryService.ListActiveCategories(ctx, limit, offset)
	if err != nil {
		return err
	}

	// Получаем общее количество активных категорий
	count, err := h.services.CategoryService.CountActiveCategories(ctx)
	if err != nil {
		return err
	}

	// Создаем кнопки для каждой категории
	for _, category := range categories {
		// Проверяем, выбрана ли категория
		isSelected := slices.Contains(volunteer.CategoryIDs, category.ID)

		// Выбираем эмодзи
		emoji := "❌"
		if isSelected {
			emoji = "✅"
		}

		// Формируем текст кнопки
		buttonText := fmt.Sprintf("%s %s", emoji, category.Name)

		// Создаем payload для переключения категории
		categoryIDStr := strconv.Itoa(int(category.ID))
		categoryPayload := EncodePayload(fsm.Loop, map[string]string{
			"page":        pageStr,
			"category_id": categoryIDStr,
		})

		keyboard.AddRow().AddCallback(buttonText, schemes.DEFAULT, categoryPayload)
	}

	// Добавляем пагинацию
	totalPages := int((count + int64(limit) - 1) / int64(limit))

	if totalPages > 1 {
		row := keyboard.AddRow()

		if page > 1 {
			previousPageStr := strconv.Itoa(page - 1)
			previousPayload := EncodePayload(fsm.Loop, map[string]string{"page": previousPageStr})
			row.AddCallback("<<", schemes.DEFAULT, previousPayload)
		}

		row.AddCallback(strconv.Itoa(page), schemes.DEFAULT, fsm.Loop.String())

		if page < totalPages {
			nextPageStr := strconv.Itoa(page + 1)
			nextPayload := EncodePayload(fsm.Loop, map[string]string{"page": nextPageStr})
			row.AddCallback(">>", schemes.DEFAULT, nextPayload)
		}
	}

	// Кнопка "Назад"
	eventsPayload := EncodePayload(fsm.CategoriesFilterToEvents, map[string]string{"page": "1"})
	keyboard.AddRow().AddCallback("Назад", schemes.DEFAULT, eventsPayload)

	msg := maxbot.NewMessage().
		SetUser(update.GetUserID()).
		SetText("Выберите категории событий, которые вас интересуют:").
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
func (h *CategoryFilterHandler) LeaveState(ctx context.Context, update schemes.UpdateInterface, availableTransitions []string) (fsm.Transition, map[string]string, error) {
	switch upd := update.(type) {
	case *schemes.MessageCallbackUpdate:
		// Обработка callback кнопок
		event, params, err := DecodePayload(upd.Callback.Payload)
		if err != nil {
			return fsm.Error, nil, fmt.Errorf("неверный callback")
		}

		// Если это fsm.Loop, обрабатываем переключение категории
		if event == fsm.Loop {
			// Проверяем, есть ли category_id в параметрах
			if categoryIDStr, exists := params["category_id"]; exists {
				categoryID, err := strconv.Atoi(categoryIDStr)
				if err != nil {
					return fsm.Error, nil, fmt.Errorf("неверный ID категории")
				}

				// Получаем текущие категории волонтера
				volunteer, err := h.services.VolunteerService.GetVolunteer(ctx, update.GetUserID())
				if err != nil {
					return fsm.Error, nil, err
				}

				// Переключаем категорию (добавляем, если её нет, удаляем, если есть)
				categoryID32 := int32(categoryID)
				newCategoryIDs := make([]int32, 0, len(volunteer.CategoryIDs))
				found := false

				for _, id := range volunteer.CategoryIDs {
					if id == categoryID32 {
						found = true
						// Пропускаем эту категорию (удаляем)
					} else {
						newCategoryIDs = append(newCategoryIDs, id)
					}
				}

				// Если категория не была найдена, добавляем её
				if !found {
					newCategoryIDs = append(newCategoryIDs, categoryID32)
				}

				// Обновляем категории волонтера
				_, err = h.services.VolunteerService.UpdateVolunteerCategories(ctx, update.GetUserID(), newCategoryIDs)
				if err != nil {
					return fsm.Error, nil, err
				}

				// Возвращаем fsm.Loop, чтобы остаться на текущей странице
				return fsm.Loop, params, nil
			}

			// Если category_id нет, это просто пагинация
			return fsm.Loop, params, nil
		}

		// Проверяем, доступен ли переход
		if !slices.Contains(availableTransitions, event.String()) {
			return fsm.Error, nil, fmt.Errorf("действие недоступно")
		}

		return event, params, nil
	default:
		return fsm.Error, nil, fmt.Errorf("воспользуйтесь кнопками меню")
	}
}
