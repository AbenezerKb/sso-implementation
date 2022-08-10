package initiator

import (
	"github.com/casbin/casbin/v2"
	"sso/internal/glue/routing/oauth"
	"sso/platform/logger"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine, group *gin.RouterGroup, handler Handler, module Module, log logger.Logger, enforcer *casbin.Enforcer) {
	oauth.InitRoute(group, handler.oauth)
}
