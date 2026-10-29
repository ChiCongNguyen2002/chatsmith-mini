package auth

const (
	tableSessions = "sessions"
	whereSession  = "token_hash = ?"

	CookieSession = "chatsmith_session"
	CookieState   = "chatsmith_oauth_state"

	stateSep   = "."
	stateParts = 3
)
