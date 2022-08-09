package initiator

import (
	"sso/internal/handler/rest"
	"sso/internal/handler/rest/oauth"
	"sso/platform/logger"
)

type Handler struct {
	// TODO implement
	oauth rest.OAuth
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		// TODO implement
		oauth: oauth.InitOAuth(log, module.OAuthModule),
	}
}
