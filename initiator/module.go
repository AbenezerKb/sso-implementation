package initiator

import (
	"context"
	"io/ioutil"
	"sso/internal/module"
	"sso/internal/module/client"
	"sso/internal/module/oauth"
	"sso/internal/module/user"
	"sso/platform/logger"

	"github.com/spf13/viper"

	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type Module struct {
	// TODO implement
	OAuthModule  module.OAuthModule
	userModule   module.UserModule
	clientModule module.ClientModule
}

func InitModule(persistence Persistence, cache CacheLayer, privateKeyPath string, platformLayer PlatformLayer, log logger.Logger, enforcer *casbin.Enforcer) Module {
	keyFile, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(context.Background(), "failed to read private key", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		log.Fatal(context.Background(), "failed to parse private key", zap.Error(err))
	}

	return Module{
		// OAuthModule: oauth.InitOAuth(log, persistence.OAuthPersistence, cache.OTPCacheLayer, cache.SessionCacheLayer, privateKey, platformLayer.sms),
		userModule: user.Init(log, persistence.OAuthPersistence, platformLayer.sms, enforcer),
		OAuthModule: oauth.InitOAuth(
			log,
			persistence.OAuthPersistence,
			cache.OTPCacheLayer,
			cache.SessionCacheLayer,
			privateKey,
			platformLayer.sms,
			oauth.SetOptions(oauth.Options{
				AccessTokenExpireTime:  viper.GetDuration("server.login.access_token.expire_time"),
				RefreshTokenExpireTime: viper.GetDuration("server.login.refresh_token.expire_time"),
			}),
		),
		clientModule: client.InitClient(log, persistence.ClientPersistence),
	}
}
