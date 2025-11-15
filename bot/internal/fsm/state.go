package fsm

import "strconv"

type State int
type Transition int

const (
	Empty State = iota
	NewUser
	SelectRole
	MainMenu
	Verifications
	About
	Applications
	PersonalEvents
	Events
	Event
	CategoriesFilter
	GeoFilter
	EditGeoFilter
	Verification
	ReplyVerification
	EditVerification
)

const (
	EmptyToNewUser Transition = iota
	NewUserToSelectRole
	SelectRoleToMainMenu

	MainMenuToSelectRole
	MainMenuToAbout
	MainMenuToApplications
	MainMenuToEvents
	MainMenuToPersonalEvents
	MainMenuToVerifications

	PersonalEventsToEvents

	EventsToCategoriesFilter
	EventsToGeoFilter
	EventsToEvent

	GeoFilterToEditGeoFilter

	VerificationsToVerification
	VerificationToReplyVerification
	VerificationToEditVerification
	VerificationToVerifications
	ReplyVerificationToVerification
	EditVerificationToVerification

	Reset
	Error
	Loop
)

func (s State) String() string {
	return strconv.Itoa(int(s))
}

func (e Transition) String() string {
	return strconv.Itoa(int(e))
}

// ParseState преобразует строку в State
func ParseState(s string) (State, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return Empty, err
	}
	return State(i), nil
}

// ParseTransition преобразует строку в Transition
func ParseTransition(s string) (Transition, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return EmptyToNewUser, err
	}
	return Transition(i), nil
}
