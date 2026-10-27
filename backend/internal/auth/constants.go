package auth

const (
	CookieSession = "chatsmith_session"
	CookieState   = "chatsmith_oauth_state"

	PathGoogleLogin    = "/auth/google/login"
	PathGoogleCallback = "/auth/google/callback"
	PathLogout         = "/auth/logout"

	googleAuthURL  = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL = "https://oauth2.googleapis.com/token"

	scopeOpenID  = "openid"
	scopeEmail   = "email"
	scopeProfile = "profile"

	claimEmail         = "email"
	claimEmailVerified = "email_verified"
	claimName          = "name"
	claimPicture       = "picture"

	extraIDToken = "id_token"

	logGoogleExchangeFailed  = "google_exchange_failed"
	logGoogleIDTokenMissing  = "google_id_token_missing"
	logGoogleIDTokenInvalid  = "google_id_token_invalid"
	logGoogleEmailUnverified = "google_email_unverified"

	tokenSep = "."
	fieldSep = "|"
)
