package initiator

import (
	"sso/internal/handler/rest"
	"sso/internal/handler/rest/asset"
	"sso/internal/handler/rest/client"
	"sso/internal/handler/rest/identity-provider"
	"sso/internal/handler/rest/mini_ride"
	"sso/internal/handler/rest/oauth"
	"sso/internal/handler/rest/oauth2"
	"sso/internal/handler/rest/profile"
	resource_server "sso/internal/handler/rest/resource-server"
	"sso/internal/handler/rest/role"
	rs_api "sso/internal/handler/rest/rs-api"
	"sso/internal/handler/rest/scope"
	"sso/internal/handler/rest/user"
	"sso/platform/logger"
	"sso/platform/utils"

	"github.com/spf13/viper"
)

type Handler struct {
	oauth            rest.OAuth
	oauth2           rest.OAuth2
	user             rest.User
	client           rest.Client
	scope            rest.Scope
	profile          rest.Profile
	miniRide         rest.MiniRide
	resourceServer   rest.ResourceServer
	role             rest.Role
	identityProvider rest.IdentityProvider
	rsAPI            rest.RSAPI
	asset            rest.Asset
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		oauth: oauth.InitOAuth(
			log.Named("oauth-handler"),
			module.OAuthModule,
			oauth.SetOptions(oauth.Options{
				RefreshTokenCookie: utils.CookieOptions{
					Path:     viper.GetString("server.cookies.refresh_token.path"),
					Domain:   viper.GetString("server.cookies.refresh_token.domain"),
					MaxAge:   viper.GetInt("server.cookies.refresh_token.max_age"),
					Secure:   viper.GetBool("server.cookies.refresh_token.secure"),
					HttpOnly: viper.GetBool("server.cookies.refresh_token.http_only"),
					SameSite: viper.GetInt("server.cookies.refresh_token.same_site"),
				},
				OPBSCookie: utils.CookieOptions{
					Path:     viper.GetString("server.cookies.opbs.path"),
					Domain:   viper.GetString("server.cookies.opbs.domain"),
					MaxAge:   viper.GetInt("server.cookies.opbs.max_age"),
					Secure:   viper.GetBool("server.cookies.opbs.secure"),
					HttpOnly: viper.GetBool("server.cookies.opbs.http_only"),
					SameSite: viper.GetInt("server.cookies.opbs.same_site"),
				},
			})),
		user:   user.Init(log.Named("user-handler"), module.userModule),
		client: client.Init(log.Named("client-handler"), module.clientModule),
		oauth2: oauth2.InitOAuth2(
			log.Named("oauth2-handler"),
			module.OAuth2Module,
			oauth2.SetOptions(oauth2.Options{
				OPBSCookie: utils.CookieOptions{
					Path:     viper.GetString("server.cookies.opbs.path"),
					Domain:   viper.GetString("server.cookies.opbs.domain"),
					MaxAge:   viper.GetInt("server.cookies.opbs.max_age"),
					Secure:   viper.GetBool("server.cookies.opbs.secure"),
					HttpOnly: viper.GetBool("server.cookies.opbs.http_only"),
					SameSite: viper.GetInt("server.cookies.opbs.same_site"),
				},
			})),
		scope:            scope.InitScope(log.Named("scope-handler"), module.scopeModule),
		profile:          profile.Init(log.Named("profile-handler"), module.profile),
		miniRide:         mini_ride.Init(log.Named("minRide-handler"), module.MiniRideModule),
		resourceServer:   resource_server.Init(log.Named("resource-server-handler"), module.resourceServer),
		role:             role.InitRole(log.Named("role-handler"), module.RoleModule),
		identityProvider: identity_provider.InitIdentityProvider(log.Named("identity-provider-handler"), module.identityProvider),
		rsAPI:            rs_api.Init(log.Named("rs_api"), module.rsAPI),
		asset:            asset.Init(log.Named("asset-handler"), module.asset),
	}
}
