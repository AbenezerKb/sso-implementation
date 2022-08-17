package initiator

import (
	"sso/internal/handler/rest"
	"sso/internal/handler/rest/client"
	"sso/internal/handler/rest/oauth"
	"sso/internal/handler/rest/user"
	"sso/platform/logger"
)

type Handler struct {
	oauth  rest.OAuth
	user   rest.User
	client rest.Client
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		oauth:  oauth.InitOAuth(log.Named("oauth-handler"), module.OAuthModule),
		user:   user.Init(log.Named("user-handler"), module.userModule),
		client: client.Init(log.Named("client-handler"), module.clientModule),
	}
}
