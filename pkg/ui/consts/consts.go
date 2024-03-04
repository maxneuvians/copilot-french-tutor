package consts

const (
	// File names
	TokenFile = ".github_copilot_token"

	// Panes
	ChatPane       = "Chat"
	ExercisePane   = "Exercise"
	LoginPane      = "Login"
	QuizPane       = "Quiz"
	VocabularyPane = "Vocabulary"

	// Session states
	LoggedIn  = "Logged in"
	LoggingIn = "Logging in"
	LoggedOut = "Logged out"
)

// Global message types
type InitializingMsg bool
type SessionUpdateMsg struct {
	AccessToken string
	State       string
}
type SessionTokenUpdateMsg string
