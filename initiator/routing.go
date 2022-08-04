package initiator

import (
	"sso/internal/glue/routing/oauth"
	"sso/platform/logger"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine, group *gin.RouterGroup, handler Handler, module Module, log logger.Logger) {
	oauth.InitRoute(group, handler.oauth)
}
