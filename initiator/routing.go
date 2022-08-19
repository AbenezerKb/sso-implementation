package initiator

import (
	"context"
	"io/ioutil"
	"sso/internal/glue/routing/client"

	"sso/docs"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"sso/internal/glue/routing/oauth"
	"sso/internal/glue/routing/oauth2"
	"sso/internal/glue/routing/user"
	"sso/internal/handler/middleware"
	"sso/platform/logger"
)

func InitRouter(router *gin.Engine, group *gin.RouterGroup, handler Handler, module Module, log logger.Logger, enforcer *casbin.Enforcer, publicKeyPath string) {
	certificate, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal(context.Background(), "Error reading own certificate : \n", zap.Error(err))
	}
	ssoPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(certificate)
	if err != nil {
		log.Fatal(context.Background(), "Error parsing own certificate : \n", zap.Error(err))
	}
	authMiddleware := middleware.InitAuthMiddleware(enforcer, module.OAuthModule, ssoPublicKey, module.clientModule, log.Named("auth-middleware"))

	docs.SwaggerInfo.BasePath = "/v1"
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	oauth.InitRoute(group, handler.oauth, authMiddleware, enforcer)

	// new group for oauth2

	oauth2.InitRoute(group, handler.oauth2, authMiddleware, enforcer)
	user.InitRoute(group, handler.user, authMiddleware, enforcer)
	client.InitRoute(group, handler.client, authMiddleware, enforcer)
}
