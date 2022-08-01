package initiator

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"sso/platform/logger"
)

func InitRouter(router *gin.Engine, handler Handler, module Module, enforcer *casbin.Enforcer, log logger.Logger) {
	// TODO implement
}
