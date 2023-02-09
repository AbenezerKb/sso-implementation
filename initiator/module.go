package initiator

import (
	"sso/internal/module"
	"sso/internal/module/asset"
	"sso/internal/module/client"
	identity_provider "sso/internal/module/identity-provider"
	"sso/internal/module/mini_ride"
	"sso/internal/module/oauth"
	"sso/internal/module/oauth2"
	"sso/internal/module/profile"
	resource_server "sso/internal/module/resource-server"
	"sso/internal/module/role"
	rs_api "sso/internal/module/rs-api"
	"sso/internal/module/scope"
	"sso/internal/module/user"
	"sso/platform/logger"

	"github.com/spf13/viper"

	"github.com/casbin/casbin/v2"
)

type Module struct {
	OAuthModule      module.OAuthModule
	OAuth2Module     module.OAuth2Module
	userModule       module.UserModule
	clientModule     module.ClientModule
	scopeModule      module.ScopeModule
	profile          module.ProfileModule
	resourceServer   module.ResourceServerModule
	MiniRideModule   module.MiniRideModule
	RoleModule       module.RoleModule
	identityProvider module.IdentityProviderModule
	rsAPI            module.RSAPI
	asset            module.Asset
}

func InitModule(persistence Persistence, cache CacheLayer, privateKeyPath string, platformLayer PlatformLayer, log logger.Logger, enforcer *casbin.Enforcer, state State) Module {
	miniRideModule := mini_ride.InitMinRide(log, persistence.MiniRidePersistence, platformLayer.Kafka)

	return Module{
		userModule: user.Init(
			log.Named("user-module"),
			persistence.OAuthPersistence,
			persistence.UserPersistence,
			persistence.RolePersistence,
			platformLayer.Sms, enforcer),
		OAuthModule: oauth.InitOAuth(
			log.Named("oauth-module"),
			persistence.OAuthPersistence,
			persistence.IdentityProviderPersistence,
			cache.OTPCacheLayer,
			cache.SessionCacheLayer,
			platformLayer.Token,
			platformLayer.Sms,
			platformLayer.SelfIP,
			cache.ResetCodeCacheLayer,
			oauth.SetOptions(oauth.Options{
				AccessTokenExpireTime:  viper.GetDuration("server.login.access_token.expire_time"),
				RefreshTokenExpireTime: viper.GetDuration("server.login.refresh_token.expire_time"),
				IDTokenExpireTime:      viper.GetDuration("server.login.id_token.expire_time"),
				ExcludedPhones:         state.ExcludedPhones,
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
		profile: profile.InitProfile(
			log.Named("profile-module"),
			persistence.OAuthPersistence,
			persistence.ProfilePersistence,
			cache.OTPCacheLayer,
			profile.SetOptions(profile.Options{
				ProfilePictureDist:    viper.GetString("assets.profile_picture_dst"),
				ProfilePictureMaxSize: viper.GetInt("assets.profile_picture_max_size"),
			})),
		resourceServer:   resource_server.InitResourceServer(log.Named("resource-server-module"), persistence.ResourceServerPersistence, persistence.ScopePersistence),
		RoleModule:       role.InitRole(log.Named("role-module"), persistence.RolePersistence),
		identityProvider: identity_provider.InitIdentityProvider(log.Named("identity-provider-module"), persistence.IdentityProviderPersistence),
		rsAPI:            rs_api.Init(log.Named("rs_api_module"), persistence.UserPersistence),
		asset:            asset.Init(log.Named("asset-module"), platformLayer.Asset, state.UploadParams),
		MiniRideModule:   miniRideModule,
	}
}

func InitMockModule(persistence Persistence, cache CacheLayer, privateKeyPath string, platformLayer PlatformLayer, log logger.Logger, enforcer *casbin.Enforcer, state State, path string) Module {
	return Module{
		userModule: user.Init(
			log.Named("user-module"),
			persistence.OAuthPersistence,
			persistence.UserPersistence,
			persistence.RolePersistence,
			platformLayer.Sms, enforcer),
		OAuthModule: oauth.InitOAuth(
			log.Named("oauth-module"),
			persistence.OAuthPersistence,
			persistence.IdentityProviderPersistence,
			cache.OTPCacheLayer,
			cache.SessionCacheLayer,
			platformLayer.Token,
			platformLayer.Sms,
			platformLayer.SelfIP,
			cache.ResetCodeCacheLayer,
			oauth.SetOptions(oauth.Options{
				AccessTokenExpireTime:  viper.GetDuration("server.login.access_token.expire_time"),
				RefreshTokenExpireTime: viper.GetDuration("server.login.refresh_token.expire_time"),
				IDTokenExpireTime:      viper.GetDuration("server.login.id_token.expire_time"),
				ExcludedPhones:         state.ExcludedPhones,
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
		profile: profile.InitProfile(
			log.Named("profile-module"),
			persistence.OAuthPersistence,
			persistence.ProfilePersistence,
			cache.OTPCacheLayer,
			profile.SetOptions(profile.Options{
				ProfilePictureDist:    path + viper.GetString("assets.profile_picture_dst"),
				ProfilePictureMaxSize: viper.GetInt("assets.profile_picture_max_size"),
			})),
		resourceServer:   resource_server.InitResourceServer(log.Named("resource-server-module"), persistence.ResourceServerPersistence, persistence.ScopePersistence),
		MiniRideModule:   mini_ride.InitMinRide(log.Named("mini-ride-module"), persistence.MiniRidePersistence, platformLayer.Kafka),
		RoleModule:       role.InitRole(log.Named("role-module"), persistence.RolePersistence),
		identityProvider: identity_provider.InitIdentityProvider(log.Named("identity-provider-module"), persistence.IdentityProviderPersistence),
		rsAPI:            rs_api.Init(log.Named("rs_api_module"), persistence.UserPersistence),
		asset:            asset.Init(log.Named("asset-module"), platformLayer.Asset, state.UploadParams),
	}
}
