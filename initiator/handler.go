package initiator

import (
	"sso/internal/handler/rest"
	"sso/internal/handler/rest/client"
	"sso/internal/handler/rest/oauth"
	"sso/internal/handler/rest/oauth2"
	"sso/internal/handler/rest/user"
	"sso/platform/logger"
)

type Handler struct {
	// TODO implement
	oauth  rest.OAuth
	oauth2 rest.OAuth2
	user   rest.User
	client rest.Client
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		// TODO implement
		oauth:  oauth.InitOAuth(log.Named("oauth-handler"), module.OAuthModule),
		user:   user.Init(log.Named("user-handler"), module.userModule),
		client: client.Init(log.Named("client-handler"), module.clientModule),
		oauth2: oauth2.InitOAuth2(log.Named("oauth2-handler"), module.OAuth2Module),
	}
}
