package initiator

import (
	"sso/internal/handler/rest"
	"sso/internal/handler/rest/client"
	"sso/internal/handler/rest/mini_ride"
	"sso/internal/handler/rest/oauth"
	"sso/internal/handler/rest/oauth2"
	"sso/internal/handler/rest/profile"
	"sso/internal/handler/rest/scope"
	"sso/internal/handler/rest/user"
	"sso/platform/logger"
)

type Handler struct {
	oauth    rest.OAuth
	oauth2   rest.OAuth2
	user     rest.User
	client   rest.Client
	scope    rest.Scope
	profile  rest.Profile
	miniRide rest.MiniRide
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		oauth:    oauth.InitOAuth(log.Named("oauth-handler"), module.OAuthModule),
		user:     user.Init(log.Named("user-handler"), module.userModule),
		client:   client.Init(log.Named("client-handler"), module.clientModule),
		oauth2:   oauth2.InitOAuth2(log.Named("oauth2-handler"), module.OAuth2Module),
		scope:    scope.InitScope(log.Named("scope-handler"), module.scopeModule),
		profile:  profile.Init(log.Named("profile-handler"), module.profile),
		miniRide: mini_ride.Init(log.Named("minRide-handler"), module.Mini_rideModule),
	}
}
