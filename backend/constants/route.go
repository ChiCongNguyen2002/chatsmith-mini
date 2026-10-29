package constants

const (
	ParamConversationID = "id"
	ParamProvider       = "provider"

	PathConversations        = "/conversations"
	PathConversation         = PathConversations + "/:" + ParamConversationID
	PathConversationMessages = PathConversation + "/messages"
	PathConversationContinue = PathConversation + "/continue"

	PathAuthLogin    = "/auth/:" + ParamProvider + "/login"
	PathAuthCallback = "/auth/:" + ParamProvider + "/callback"
	PathAuthLogout   = "/auth/logout"
)
