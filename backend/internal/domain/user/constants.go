package user

const (
	tableUsers      = "users"
	tableIdentities = "user_identities"

	whereUserID        = "id = ?"
	whereIdentity      = "user_identities.provider = ? AND user_identities.subject = ?"
	whereVerifiedEmail = "user_identities.email = ? AND user_identities.email_verified = true"
	joinIdentityUser   = "JOIN user_identities ON user_identities.user_id = users.id"
	selectUserColumns  = "users.*"
	updateProfileSQL   = `UPDATE users SET email = ?, name = ?, picture = ? WHERE id = ?`
)
