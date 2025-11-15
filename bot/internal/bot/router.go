package internal

import (
	"context"
	"log"

	"maxBot/internal/bot/handler"
	"maxBot/internal/di"
	"maxBot/internal/fsm"
	"maxBot/internal/model"

	maxbot "github.com/rectid/max-bot-api-client-go"
	"github.com/rectid/max-bot-api-client-go/schemes"
)

// Router управляет маршрутизацией между FSM и хендлерами
type Router struct {
	services              *di.Services
	handlers              map[fsm.State]handler.Handler
	emptyHandler          *handler.EmptyHandler
	newUserHandler        *handler.NewUserHandler
	selectRoleHandler     *handler.SelectRoleHandler
	mainMenuHandler       *handler.MainMenuHandler
	eventsHandler         *handler.EventsHandler
	eventHandler          *handler.EventHandler
	personalEventsHandler *handler.PersonalEventsHandler
	categoryFilterHandler *handler.CategoryFilterHandler
	geoFilterHandler      *handler.GeoFilterHandler
	editGeoFilterHandler  *handler.EditGeoFilterHandler
}

// NewRouter создаёт новый роутер с инициализированными хендлерами
func NewRouter(services *di.Services) *Router {
	r := &Router{
		services: services,
		handlers: make(map[fsm.State]handler.Handler),
	}

	r.newUserHandler = handler.NewNewUserHandler(services)
	r.emptyHandler = handler.NewEmptyHandler(services)
	r.selectRoleHandler = handler.NewSelectRoleHandler(services)
	r.mainMenuHandler = handler.NewMainMenuHandler(services)
	r.eventsHandler = handler.NewEventsHandler(services)
	r.eventHandler = handler.NewEventHandler(services)
	r.personalEventsHandler = handler.NewPersonalEventsHandler(services)
	r.categoryFilterHandler = handler.NewCategoryFilterHandler(services)
	r.geoFilterHandler = handler.NewGeoFilterHandler(services)
	r.editGeoFilterHandler = handler.NewEditGeoFilterHandler(services)

	r.handlers[fsm.Empty] = r.emptyHandler
	r.handlers[fsm.NewUser] = r.newUserHandler
	r.handlers[fsm.SelectRole] = r.selectRoleHandler
	r.handlers[fsm.MainMenu] = r.mainMenuHandler
	r.handlers[fsm.Events] = r.eventsHandler
	r.handlers[fsm.Event] = r.eventHandler
	r.handlers[fsm.PersonalEvents] = r.personalEventsHandler
	r.handlers[fsm.CategoriesFilter] = r.categoryFilterHandler
	r.handlers[fsm.GeoFilter] = r.geoFilterHandler
	r.handlers[fsm.EditGeoFilter] = r.editGeoFilterHandler
	return r
}

// RouteUpdate обрабатывает любой апдейт
func (r *Router) RouteUpdate(ctx context.Context, user *model.User, update schemes.UpdateInterface) {
	currentState, err := fsm.ParseState(user.State)
	if err != nil {
		log.Printf("Failed to parse user state: %v", err)
		return
	}

	onStateUpdate := func(newState fsm.State) error {
		updated, err := r.services.UserService.UpdateUserState(ctx, user.ID, newState.String())
		if err != nil {
			log.Printf("failed to update user state: %v", err)
			return err
		}
		user.State = updated.State
		currentState = newState
		return nil
	}
	machine := fsm.NewUserFSM(currentState, onStateUpdate)
	h := r.handlers[currentState]
	if h == nil {
		log.Printf("No handler found for state: %s", user.State)
		return
	}
	transition, params, err := h.LeaveState(ctx, update, machine.AvailableTransitions())
	if transition == fsm.Loop {
		h.EnterState(ctx, update, transition, params)
		return
	}
	if transition == fsm.Error {
		if err != nil {
			r.services.API.Messages.Send(ctx, maxbot.NewMessage().SetUser(update.GetUserID()).SetText(err.Error()))
		}
		h.EnterState(ctx, update, transition, params)
		return
	}
	err = machine.Event(ctx, transition.String())
	if err != nil {
		log.Printf("State change failed for user: %d", user.ID)
		return
	}
	h = r.handlers[currentState]
	h.EnterState(ctx, update, transition, params)
}
