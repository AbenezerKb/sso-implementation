package initiator

import (
	"context"
	"io/ioutil"
	"sso/internal/glue/routing/oauth"
	"sso/internal/handler/middleware"
	"sso/platform/logger"

	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
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
	authMiddleware := middleware.InitAuthMiddleware(enforcer, module.OAuthModule, ssoPublicKey)

	oauth.InitRoute(group, handler.oauth, authMiddleware)
}
