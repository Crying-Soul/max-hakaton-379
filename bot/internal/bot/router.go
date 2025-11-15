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
	services                 *di.Services
	handlers                 map[fsm.State]handler.Handler
	emptyHandler             *handler.EmptyHandler
	newUserHandler           *handler.NewUserHandler
	selectRoleHandler        *handler.SelectRoleHandler
	mainMenuHandler          *handler.MainMenuHandler
	eventsHandler            *handler.EventsHandler
	verificationsHandler     *handler.VerificationsHandler
	verificationHandler      *handler.VerificationHandler
	replyVerificationHandler *handler.ReplyVerificationHandler
	editVerificationHandler  *handler.EditVerificationHandler
	// aboutHandler          *handler.AboutHandler
	// applicationsHandler   *handler.ApplicationsHandler
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
	r.verificationsHandler = handler.NewVerificationsHandler(services)
	r.verificationHandler = handler.NewVerificationHandler(services)
	r.replyVerificationHandler = handler.NewReplyVerificationHandler(services)
	r.editVerificationHandler = handler.NewEditVerificationHandler(services)
	// r.aboutHandler = handler.NewAboutHandler(services)
	// r.applicationsHandler = handler.NewApplicationsHandler(services)

	r.handlers[fsm.Empty] = r.emptyHandler
	r.handlers[fsm.NewUser] = r.newUserHandler
	r.handlers[fsm.SelectRole] = r.selectRoleHandler
	r.handlers[fsm.MainMenu] = r.mainMenuHandler
	// r.handlers[fsm.About] = r.aboutHandler
	// r.handlers[fsm.Applications] = r.applicationsHandler
	r.handlers[fsm.Events] = r.eventsHandler
	r.handlers[fsm.Verifications] = r.verificationsHandler
	r.handlers[fsm.Verification] = r.verificationHandler
	r.handlers[fsm.ReplyVerification] = r.replyVerificationHandler
	r.handlers[fsm.EditVerification] = r.editVerificationHandler
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
