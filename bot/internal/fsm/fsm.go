package fsm

import (
	"context"

	lfsm "github.com/looplab/fsm"
)

type FSM struct {
	*lfsm.FSM
}

// NewUserFSM создаёт FSM для управления состояниями пользователя
func NewUserFSM(
	initial State,
	stateUpdate func(newState State) error,
) *FSM {
	callbacks := map[string]lfsm.Callback{
		"leave_state": func(ctx context.Context, e *lfsm.Event) {
			state, err := ParseState(e.Dst)
			if err != nil {
				e.Cancel()
			}
			err = stateUpdate(state)
			if err != nil {
				e.Cancel()
			}
		},
	}

	userFSM := &FSM{
		FSM: lfsm.NewFSM(
			initial.String(),
			lfsm.Events{
				{Name: EmptyToNewUser.String(), Src: []string{Empty.String()}, Dst: NewUser.String()},
				{Name: NewUserToSelectRole.String(), Src: []string{NewUser.String()}, Dst: SelectRole.String()},
				{Name: SelectRoleToMainMenu.String(), Src: []string{SelectRole.String()}, Dst: MainMenu.String()},
				{Name: MainMenuToSelectRole.String(), Src: []string{MainMenu.String()}, Dst: SelectRole.String()},
				{Name: MainMenuToAbout.String(), Src: []string{MainMenu.String()}, Dst: About.String()},
				{Name: MainMenuToApplications.String(), Src: []string{MainMenu.String()}, Dst: Applications.String()},
				{Name: MainMenuToEvents.String(), Src: []string{MainMenu.String()}, Dst: Events.String()},
				{Name: MainMenuToPersonalEvents.String(), Src: []string{MainMenu.String()}, Dst: PersonalEvents.String()},
				{Name: MainMenuToVerifications.String(), Src: []string{MainMenu.String()}, Dst: Verifications.String()},
				{Name: VerificationsToVerification.String(), Src: []string{Verifications.String()}, Dst: Verification.String()},
				{Name: VerificationToVerifications.String(), Src: []string{Verification.String()}, Dst: Verifications.String()},
				{Name: VerificationToReplyVerification.String(), Src: []string{Verification.String()}, Dst: ReplyVerification.String()},
				{Name: ReplyVerificationToVerification.String(), Src: []string{ReplyVerification.String()}, Dst: Verification.String()},
				{Name: VerificationToEditVerification.String(), Src: []string{Verification.String()}, Dst: EditVerification.String()},
				{Name: EditVerificationToVerification.String(), Src: []string{EditVerification.String()}, Dst: Verification.String()},
				{Name: EventsToCategoriesFilter.String(), Src: []string{Events.String()}, Dst: CategoriesFilter.String()},
				{Name: CategoriesFilterToEvents.String(), Src: []string{CategoriesFilter.String()}, Dst: Events.String()},
				{Name: EventsToGeoFilter.String(), Src: []string{Events.String()}, Dst: GeoFilter.String()},
				{Name: GeoFilterToEvents.String(), Src: []string{GeoFilter.String()}, Dst: Events.String()},
				{Name: EventsToMainMenu.String(), Src: []string{Events.String()}, Dst: MainMenu.String()},
				{Name: GeoFilterToEditGeoFilter.String(), Src: []string{GeoFilter.String()}, Dst: EditGeoFilter.String()},
				{Name: EditGeoFilterToGeoFilter.String(), Src: []string{EditGeoFilter.String()}, Dst: GeoFilter.String()},
				{Name: EventToEvents.String(), Src: []string{Event.String()}, Dst: Events.String()},
				{Name: EventsToEvent.String(), Src: []string{Events.String()}, Dst: Event.String()},
				{Name: Reset.String(), Src: []string{"*"}, Dst: Empty.String()},
				{Name: Error.String(), Src: []string{"*"}, Dst: Empty.String()},
			},
			callbacks,
		),
	}

	return userFSM
}
