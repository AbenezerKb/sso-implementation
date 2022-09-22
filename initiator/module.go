package initiator

import (
	"sso/internal/module"
	"sso/internal/module/client"
	"sso/internal/module/oauth"
	"sso/internal/module/oauth2"
	"sso/internal/module/scope"
	"sso/internal/module/user"
	"sso/platform/logger"

	"github.com/spf13/viper"

	"github.com/casbin/casbin/v2"
)

type Module struct {
	OAuthModule  module.OAuthModule
	OAuth2Module module.OAuth2Module
	userModule   module.UserModule
	clientModule module.ClientModule
	scopeModule  module.ScopeModule
}

func InitModule(persistence Persistence, cache CacheLayer, privateKeyPath string, platformLayer PlatformLayer, log logger.Logger, enforcer *casbin.Enforcer, state State) Module {

	return Module{
		userModule: user.Init(log.Named("user-module"), persistence.OAuthPersistence, persistence.UserPersistence, platformLayer.Sms, enforcer),
		OAuthModule: oauth.InitOAuth(
			log.Named("oauth-module"),
			persistence.OAuthPersistence,
			cache.OTPCacheLayer,
			cache.SessionCacheLayer,
			platformLayer.Token,
			platformLayer.Sms,
			oauth.SetOptions(oauth.Options{
				AccessTokenExpireTime:  viper.GetDuration("server.login.access_token.expire_time"),
				RefreshTokenExpireTime: viper.GetDuration("server.login.refresh_token.expire_time"),
				IDTokenExpireTime:      viper.GetDuration("server.login.id_token.expire_time"),
			}),
		),
		clientModule: client.InitClient(log.Named("client-module"), persistence.ClientPersistence),
		OAuth2Module: oauth2.InitOAuth2(
			log.Named("oauth2-module"),
			persistence.OAuth2Persistence,
			persistence.OAuthPersistence,
			persistence.ClientPersistence,
			cache.ConsentCacheLayer,
			cache.AuthCodeCacheLayer,
			platformLayer.Token,
			oauth2.SetOptions(
				oauth2.Options{
					AccessTokenExpireTime:  viper.GetDuration("server.client.access_token.expire_time"),
					RefreshTokenExpireTime: viper.GetDuration("server.client.refresh_token.expire_time"),
				},
			),
			persistence.ScopePersistence,
			state.URLs),
		scopeModule: scope.InitScope(log.Named("scope-module"), persistence.ScopePersistence),
	}
}
