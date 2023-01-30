package initiator

import (
	"context"

	"sso/internal/glue/routing/asset"
	"sso/internal/glue/routing/client"
	identity_provider "sso/internal/glue/routing/identity-provider"
	"sso/internal/glue/routing/mini_ride"
	resource_server "sso/internal/glue/routing/resource-server"
	"sso/internal/glue/routing/role"
	rs_api "sso/internal/glue/routing/rs-api"

	"sso/docs"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"sso/internal/glue/routing/oauth"
	"sso/internal/glue/routing/oauth2"
	"sso/internal/glue/routing/profile"
	"sso/internal/glue/routing/scope"
	"sso/internal/glue/routing/user"
	"sso/internal/handler/middleware"
	"sso/platform/logger"
)

func InitRouter(router *gin.Engine, group *gin.RouterGroup, handler Handler, module Module, log logger.Logger, enforcer *casbin.Enforcer, platformLayer PlatformLayer) {

	log.Info(context.Background(), "serving static files")
	group.Static("/static", viper.GetString("assets.profile_picture_dst"))

	authMiddleware := middleware.InitAuthMiddleware(
		enforcer,
		module.OAuthModule,
		platformLayer.Token,
		module.clientModule,
		middleware.MiniRideCredential{
			UserName: viper.GetString("mini_ride.username"),
			Password: viper.GetString("mini_ride.password")},
		module.RoleModule,
		module.resourceServer,
		log.Named("auth-middleware"))

	docs.SwaggerInfo.BasePath = "/v1"
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	oauth.InitRoute(group, handler.oauth, authMiddleware, enforcer)
	oauth2.InitRoute(group, handler.oauth2, authMiddleware, enforcer)
	user.InitRoute(group, handler.user, authMiddleware, enforcer)
	client.InitRoute(group, handler.client, authMiddleware, enforcer)
	scope.InitRoute(group, handler.scope, authMiddleware, enforcer)
	profile.InitRoute(group, handler.profile, authMiddleware, enforcer)
	mini_ride.InitRoute(group, handler.miniRide, authMiddleware, enforcer)
	resource_server.InitRoute(group, handler.resourceServer, authMiddleware, enforcer)
	role.InitRoute(group, handler.role, authMiddleware, enforcer)
	identity_provider.InitRoute(group, handler.identityProvider, authMiddleware, enforcer)
	rs_api.InitRoute(group, handler.rsAPI, authMiddleware, enforcer)
	asset.InitRoute(group, handler.asset, authMiddleware, enforcer)
}
