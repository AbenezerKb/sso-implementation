package initiator

import (
	"context"
	"io/ioutil"
	"sso/internal/module"
	"sso/internal/module/oauth"
	"sso/platform/logger"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Module struct {
	// TODO implement
	OAuthModule module.OAuthModule
}

func InitModule(persistence Persistence, cache *redis.Client, log logger.Logger) Module {
	key := viper.GetString("private_key")

	keyFile, err := ioutil.ReadFile(key)
	if err != nil {
		log.Fatal(context.Background(), "failed to read private key", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		log.Fatal(context.Background(), "failed to parse private key", zap.Error(err))
	}

	return Module{
		OAuthModule: oauth.InitOAuth(log, persistence.OAuthPersistence, cache, privateKey),
	}
}
