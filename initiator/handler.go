package initiator

import (
	"github.com/casbin/casbin/v2"
	"sso/platform/logger"
)

type Handler struct {
	// TODO implement
}

func InitHandler(module Module, enforcer *casbin.Enforcer, log logger.Logger) Handler {
	return Handler{}
}
