package internal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/rectid/max-bot-api-client-go/schemes"

	"maxBot/internal/di"
	"maxBot/internal/fsm"
	"maxBot/internal/model"
)

type Bot struct {
	router   *Router
	services *di.Services
}

// NewBot создаёт новый экземпляр бота используя заранее инициализированные сервисы
func NewBot(ctx context.Context, services *di.Services) (*Bot, error) {
	if services == nil {
		return nil, fmt.Errorf("services are required")
	}
	if services.API == nil {
		return nil, fmt.Errorf("api client is required")
	}

	router := NewRouter(services)

	return &Bot{
		router:   router,
		services: services,
	}, nil
}

// Start запускает обработку обновлений от API
func (b *Bot) Start(ctx context.Context) {
	log.Println("Bot started...")

	wg := sync.WaitGroup{}

	api := b.services.API
	for upd := range api.GetUpdates(ctx) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.handleUpdate(ctx, upd)
		}()
	}

	wg.Wait()
}

// handleUpdate обрабатывает одно обновление от пользователя
func (b *Bot) handleUpdate(ctx context.Context, update schemes.UpdateInterface) {
	var username, name string

	userID := update.GetUserID()
	switch upd := update.(type) {
	case *schemes.MessageCallbackUpdate:
		username = upd.Message.Sender.Username
		name = upd.Message.Sender.Name
	case *schemes.MessageCreatedUpdate:
		username = upd.Message.Sender.Username
		name = upd.Message.Sender.Name
	default:
		log.Printf("Unknown update type: %T", update)
		return
	}

	user, err := b.getOrCreateUser(ctx, userID, username, name)
	if err != nil {
		log.Printf("Failed to get/create user: %v", err)
		return
	}

	b.router.RouteUpdate(ctx, user, update)

	log.Printf("User %d processed update in state %s", user.ID, user.State)
}

// getOrCreateUser получает существующего пользователя или создаёт нового
func (b *Bot) getOrCreateUser(ctx context.Context, userID int64, username, name string) (*model.User, error) {
	user, err := b.services.UserService.GetUserByID(ctx, userID)
	if err == nil {
		return &user, nil
	}

	var usernamePtr *string
	if strings.TrimSpace(username) != "" {
		usernamePtr = &username
	}

	userModel := model.User{
		ID:          userID,
		Username:    usernamePtr,
		Name:        name,
		Role:        "volunteer",
		State:       fsm.Empty.String(),
		IsBlocked:   false,
		LocationLat: nil,
		LocationLon: nil,
	}

	created, err := b.services.UserService.UpsertUser(ctx, userModel)
	if err != nil {
		return nil, err
	}

	log.Printf("Created new user: %d (%s)", userID, username)
	return &created, nil
}

// Close закрывает соединения с БД
func (b *Bot) Close() {}
