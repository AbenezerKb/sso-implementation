package initiator

import (
	"sso/internal/handler/rest"
	"sso/internal/handler/rest/oauth"
	"sso/internal/handler/rest/user"
	"sso/platform/logger"
)

type Handler struct {
	// TODO implement
	oauth rest.OAuth
	user  rest.User
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		// TODO implement
		oauth: oauth.InitOAuth(log, module.OAuthModule),
		user:  user.Init(log, module.userModule),
	}
}
