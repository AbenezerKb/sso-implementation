package initiator

import (
	"github.com/go-redis/redis/v8"
	"sso/platform/logger"
)

type Module struct {
	// TODO implement
}

func InitModule(persistence Persistence, cache *redis.Client, log logger.Logger) Module {
	return Module{}
}
