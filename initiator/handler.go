package initiator

import (
	"github.com/spf13/viper"
	"sso/internal/handler/rest"
	"sso/internal/handler/rest/client"
	"sso/internal/handler/rest/oauth"
	"sso/internal/handler/rest/oauth2"
	"sso/internal/handler/rest/user"
	"sso/platform/logger"
)

type Handler struct {
	oauth  rest.OAuth
	oauth2 rest.OAuth2
	user   rest.User
	client rest.Client
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		oauth:  oauth.InitOAuth(log.Named("oauth-handler"), module.OAuthModule),
		user:   user.Init(log.Named("user-handler"), module.userModule),
		client: client.Init(log.Named("client-handler"), module.clientModule),
		oauth2: oauth2.InitOAuth2(log.Named("oauth2-handler"), module.OAuth2Module, oauth2.SetOptions(oauth2.Options{
			ConsentURL: viper.GetString("server.oauth2.consent_url"),
			ErrorURL:   viper.GetString("server.oauth2.error_url"),
		})),
	}
}
