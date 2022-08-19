package initiator

import (
	"sso/internal/glue/routing/client"

	"sso/docs"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"sso/internal/glue/routing/oauth"
	"sso/internal/glue/routing/oauth2"
	"sso/internal/glue/routing/user"
	"sso/internal/handler/middleware"
	"sso/platform/logger"
)

func InitRouter(router *gin.Engine, group *gin.RouterGroup, handler Handler, module Module, log logger.Logger, enforcer *casbin.Enforcer, platformLayer PlatformLayer) {

	authMiddleware := middleware.InitAuthMiddleware(enforcer, module.OAuthModule, platformLayer.token, module.clientModule, log.Named("auth-middleware"))

	docs.SwaggerInfo.BasePath = "/v1"
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	oauth.InitRoute(group, handler.oauth, authMiddleware, enforcer)
	oauth2.InitRoute(group, handler.oauth2, authMiddleware, enforcer)
	user.InitRoute(group, handler.user, authMiddleware, enforcer)
	client.InitRoute(group, handler.client, authMiddleware, enforcer)
}
