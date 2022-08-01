package initiator

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"sso/platform/logger"
)

type Module struct {
	// TODO implement
}

func InitModule(persistence Persistence, cache *redis.Client, enforcer *casbin.Enforcer, log logger.Logger) Module {
	return Module{}
}
