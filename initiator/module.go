package initiator

import (
	"sso/internal/module"
	"sso/internal/module/oauth"
	"sso/platform/logger"

	"github.com/go-redis/redis/v8"
)

type Module struct {
	// TODO implement
	OAuthModule module.OAuthModule
}

func InitModule(persistence Persistence, cache *redis.Client, log logger.Logger) Module {
	return Module{
		OAuthModule: oauth.InitOAuth(log, persistence.OAuthPersistence, cache),
	}
}
