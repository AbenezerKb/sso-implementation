package initiator

import (
	"sso/internal/handler/rest"
	"sso/internal/handler/rest/oauth"
	"sso/internal/handler/rest/oauth2"
	"sso/platform/logger"
)

type Handler struct {
	// TODO implement
	oauth  rest.OAuth
	oauth2 rest.OAuth2
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		// TODO implement
		oauth:  oauth.InitOAuth(log, module.OAuthModule),
		oauth2: oauth2.InitOAuth2(log, module.OAuth2Module),
	}
}
